# pyright: reportMissingModuleSource=false
import os
import re
from contextlib import contextmanager
from urllib.parse import unquote, urlparse


def _is_mysql_enabled():
    database_url = os.environ.get("DATABASE_URL", "")
    return database_url.startswith("mysql://") or database_url.startswith("mysql+pymysql://")


def _build_mysql_config():
    import pymysql

    database_url = os.environ.get("DATABASE_URL", "")
    if not database_url:
        raise RuntimeError("DATABASE_URL is required when using MySQL")

    normalized_url = database_url.replace("mysql+pymysql://", "mysql://", 1)
    parsed = urlparse(normalized_url)

    if not parsed.hostname or not parsed.path:
        raise RuntimeError("Invalid MySQL DATABASE_URL format")

    return {
        "host": parsed.hostname,
        "port": parsed.port or 3306,
        "user": unquote(parsed.username or ""),
        "password": unquote(parsed.password or ""),
        "database": parsed.path.lstrip("/"),
        "charset": "utf8mb4",
        "autocommit": False,
        "cursorclass": pymysql.cursors.DictCursor,
    }


def _to_mysql_query(query):
    transformed = re.sub(r"\bINSERT\s+OR\s+IGNORE\b", "INSERT IGNORE", query, flags=re.IGNORECASE)
    transformed = transformed.replace("?", "%s")
    return transformed


class _MySQLCursorWrapper:
    def __init__(self, cursor):
        self._cursor = cursor

    def execute(self, query, params=()):
        return self._cursor.execute(_to_mysql_query(query), params)

    def executemany(self, query, params_seq):
        return self._cursor.executemany(_to_mysql_query(query), params_seq)

    def __getattr__(self, item):
        return getattr(self._cursor, item)


class _MySQLConnectionWrapper:
    def __init__(self, conn):
        self._conn = conn

    def cursor(self):
        return _MySQLCursorWrapper(self._conn.cursor())

    def __getattr__(self, item):
        return getattr(self._conn, item)


@contextmanager
def db_connection(db_path):
    if not _is_mysql_enabled():
        raise RuntimeError("SQLite is disabled. Please set DATABASE_URL to a MySQL URL.")

    import pymysql

    conn = pymysql.connect(**_build_mysql_config())
    wrapped_conn = _MySQLConnectionWrapper(conn)
    try:
        yield wrapped_conn
    finally:
        conn.close()


def execute_query(db_path, query, params=(), fetch=False, fetch_one=False):
    with db_connection(db_path) as conn:
        cursor = conn.cursor()
        cursor.execute(query, params)

        if fetch_one:
            result = cursor.fetchone()
            return dict(result) if result else None
        elif fetch:
            rows = cursor.fetchall()
            return [dict(row) for row in rows]
        else:
            conn.commit()
            return cursor.lastrowid
