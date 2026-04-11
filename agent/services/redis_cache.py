import hashlib
import json
import os
import asyncio
import uuid
from typing import Any, Dict, Iterable, Optional

import redis
from redis.asyncio import ConnectionPool, Redis


def _redis_auth_kwargs() -> Dict[str, Any]:
    redis_username = os.environ.get("REDIS_USERNAME", None)
    redis_password = os.environ.get("REDIS_PASSWORD", None)
    auth_kwargs: Dict[str, Any] = {}
    if redis_username:
        auth_kwargs["username"] = redis_username
    if redis_password:
        auth_kwargs["password"] = redis_password
    return auth_kwargs


def _redis_base_config() -> Dict[str, Any]:
    redis_host = os.environ.get("REDIS_HOST", "localhost")
    redis_port = int(os.environ.get("REDIS_PORT", 6379))
    redis_db = int(os.environ.get("REDIS_DB", 0)) + 1
    return {
        "host": redis_host,
        "port": redis_port,
        "db": redis_db,
        **_redis_auth_kwargs(),
    }


def build_cache_key(module: str, endpoint: str, params: Dict[str, Any]) -> str:
    serialized = json.dumps(params, sort_keys=True, separators=(",", ":"), default=str)
    params_hash = hashlib.md5(serialized.encode("utf-8")).hexdigest()
    return f"{module}:{endpoint}:{params_hash}"


class SyncRedisCache:
    def __init__(self):
        self._client = redis.Redis(**_redis_base_config())

    def get_json(self, key: str) -> Optional[Any]:
        try:
            value = self._client.get(key)
            if value is None:
                return None
            if isinstance(value, bytes):
                value = value.decode("utf-8")
            return json.loads(value)
        except Exception:
            return None

    def set_json(self, key: str, value: Any, ttl_seconds: int) -> None:
        try:
            self._client.set(key, json.dumps(value, default=str), ex=ttl_seconds)
        except Exception:
            pass

    def delete(self, key: str) -> None:
        try:
            self._client.delete(key)
        except Exception:
            pass

    def delete_by_prefix(self, prefix: str) -> int:
        deleted = 0
        try:
            cursor = 0
            pattern = f"{prefix}*"
            while True:
                cursor, keys = self._client.scan(cursor=cursor, match=pattern, count=100)
                if keys:
                    deleted += self._client.delete(*keys)
                if cursor == 0:
                    break
        except Exception:
            return deleted
        return deleted


class AsyncRedisCache:
    def __init__(self):
        config = _redis_base_config()
        auth_kwargs = _redis_auth_kwargs()
        redis_url = f"redis://{config['host']}:{config['port']}/{config['db']}"
        self._pool = ConnectionPool.from_url(redis_url, max_connections=20, **auth_kwargs)
        self._client = Redis(connection_pool=self._pool)

    async def get_json(self, key: str) -> Optional[Any]:
        try:
            value = await self._client.get(key)
            if value is None:
                return None
            if isinstance(value, bytes):
                value = value.decode("utf-8")
            return json.loads(value)
        except Exception:
            return None

    async def set_json(self, key: str, value: Any, ttl_seconds: int) -> None:
        try:
            await self._client.set(key, json.dumps(value, default=str), ex=ttl_seconds)
        except Exception:
            pass

    async def delete(self, key: str) -> None:
        try:
            await self._client.delete(key)
        except Exception:
            pass

    async def delete_by_prefix(self, prefix: str) -> int:
        deleted = 0
        try:
            cursor = 0
            pattern = f"{prefix}*"
            while True:
                cursor, keys = await self._client.scan(cursor=cursor, match=pattern, count=100)
                if keys:
                    deleted += await self._client.delete(*keys)
                if cursor == 0:
                    break
        except Exception:
            return deleted
        return deleted

    async def delete_many(self, keys: Iterable[str]) -> None:
        keys = list(keys)
        if not keys:
            return
        try:
            await self._client.delete(*keys)
        except Exception:
            pass

    async def _release_lock(self, lock_key: str, owner_token: str) -> None:
        release_lua = """
        local lock_value = redis.call('GET', KEYS[1])
        if lock_value == ARGV[1] then
            redis.call('DEL', KEYS[1])
            return 1
        end
        return 0
        """
        try:
            await self._client.eval(release_lua, 1, lock_key, owner_token)
        except Exception:
            pass

    async def get_or_set_json_singleflight(
        self,
        key: str,
        ttl_seconds: int,
        loader,
        lock_ttl_seconds: int = 8,
        wait_timeout_ms: int = 1200,
        retry_interval_ms: int = 80,
    ) -> Any:
        cached = await self.get_json(key)
        if cached is not None:
            return cached

        lock_key = f"lock:singleflight:{key}"
        owner_token = str(uuid.uuid4())
        got_lock = False
        try:
            got_lock = bool(
                await self._client.set(
                    lock_key,
                    owner_token,
                    nx=True,
                    ex=max(1, lock_ttl_seconds),
                )
            )
        except Exception:
            got_lock = False

        if got_lock:
            try:
                # Double-check after acquiring lock.
                cached_after_lock = await self.get_json(key)
                if cached_after_lock is not None:
                    return cached_after_lock
                value = await loader()
                await self.set_json(key, value, ttl_seconds)
                return value
            finally:
                await self._release_lock(lock_key, owner_token)

        waited_ms = 0
        while waited_ms < wait_timeout_ms:
            await asyncio.sleep(retry_interval_ms / 1000.0)
            waited_ms += retry_interval_ms
            cached_retry = await self.get_json(key)
            if cached_retry is not None:
                return cached_retry

        # Fallback: avoid request starvation if lock owner failed.
        value = await loader()
        await self.set_json(key, value, ttl_seconds)
        return value


sync_redis_cache = SyncRedisCache()
async_redis_cache = AsyncRedisCache()
