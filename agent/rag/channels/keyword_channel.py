"""
Keyword search channel — BM25 over crawled_articles.

Priority 3 — complements vector search by catching exact / lexical matches
that embedding models miss (product names, version numbers, acronyms like
"GPT-4o" or "K8s") and by adding a BM25-style relevance signal that the
vector channels do not provide.

Implementation:
  Two-stage retrieval.
  1. The in-memory BM25 index (rag/keyword_index.py, Okapi BM25 from
     rank-bm25) returns the top-k article ids ranked by lexical
     relevance. The index is loaded once from MySQL and refreshed on a
     TTL.
  2. We hydrate those ids back to full rows (title, url, summary) in a
     single SQL fetch so RetrievedChunk has the metadata the rest of
     the pipeline expects.

Why not MySQL FULLTEXT? MySQL InnoDB FULLTEXT NATURAL LANGUAGE MODE
uses a TF*IDF^2 variant — not BM25. We want honest BM25 for fairness
when fusing with vector channels via RRF.
"""

from __future__ import annotations

import asyncio
import time

from loguru import logger

from rag.channels.base import SearchChannel
from rag.models import ChannelResult, ChannelType, MetadataKey, RetrievedChunk, SearchContext


class KeywordChannel(SearchChannel):
    """BM25 keyword search on crawled_articles."""

    def __init__(self, *, top_k_multiplier: int = 2):
        self._top_k_multiplier = top_k_multiplier

    def channel_type(self) -> ChannelType:
        return ChannelType.KEYWORD

    def priority(self) -> int:
        return 3

    async def search(self, context: SearchContext) -> ChannelResult:
        return await asyncio.to_thread(self._search_sync, context)

    def _search_sync(self, context: SearchContext) -> ChannelResult:
        start = time.perf_counter()

        # Use original query — rewritten queries add synonyms that help
        # vector recall but inflate BM25 noise.
        query = context.original_query
        top_k = context.top_k * self._top_k_multiplier

        # ── Stage 1: BM25 ranking ───────────────────────────────────
        try:
            hits = get_bm25_index().search(query, top_k=top_k)
        except Exception as e:
            logger.error("KeywordChannel: BM25 search failed: {e}", e=e)
            return ChannelResult.empty(ChannelType.KEYWORD)

        if not hits:
            elapsed = (time.perf_counter() - start) * 1000
            return ChannelResult(
                channel_type=ChannelType.KEYWORD,
                chunks=[],
                latency_ms=elapsed,
            )

        # ── Stage 2: Hydrate ids → full rows ────────────────────────
        article_ids = [aid for aid, _ in hits]
        scores_by_id = dict(hits)

        try:
            rows = self._fetch_rows(article_ids)
        except Exception as e:
            logger.error("KeywordChannel: hydration failed: {e}", e=e)
            return ChannelResult.empty(ChannelType.KEYWORD)

        # Preserve BM25 ranking — DB IN(...) returns unordered, so we
        # reorder by the score map.
        rows_by_id = {int(r["id"]): r for r in rows}
        ordered_rows = [
            rows_by_id[aid] for aid in article_ids if aid in rows_by_id
        ]

        # ── Build chunks ────────────────────────────────────────────
        max_score = max(scores_by_id.values(), default=0.0)
        query_lower = query.lower()
        chunks: list[RetrievedChunk] = []
        for row in ordered_rows:
            aid = int(row["id"])
            raw_score = scores_by_id.get(aid, 0.0)
            normalized = (raw_score / max_score) if max_score > 0 else 0.0

            title = row.get("title", "") or ""
            title_match = query_lower in title.lower()

            chunks.append(RetrievedChunk(
                id=str(aid),
                content=row.get("content", "") or "",
                title=title,
                url=row.get("url", "") or "",
                score=normalized,
                source_channel=ChannelType.KEYWORD,
                metadata={
                    MetadataKey.PUBLISHED_DATE: row.get("published_date", ""),
                    MetadataKey.SOURCE_ID:      str(row.get("source_id", "")),
                    MetadataKey.MATCH_TYPE:     "title" if title_match else "content",
                },
            ))

        elapsed = (time.perf_counter() - start) * 1000
        logger.info(
            "KeywordChannel: {count} results in {ms:.0f}ms (BM25)",
            count=len(chunks),
            ms=elapsed,
        )
        return ChannelResult(
            channel_type=ChannelType.KEYWORD,
            chunks=chunks,
            latency_ms=elapsed,
        )

    @staticmethod
    def _fetch_rows(article_ids: list[int]) -> list[dict]:
        """Fetch full row data for the BM25-ranked ids in one SQL call."""
        from db.config import get_tracking_db_path
        from db.connection import execute_query

        db_path = get_tracking_db_path()
        placeholders = ",".join(["%s"] * len(article_ids))
        sql = f"""
            SELECT id,
                   title,
                   url,
                   published_date,
                   COALESCE(summary, SUBSTRING(content, 1, 500)) AS content,
                   source_id
            FROM crawled_articles
            WHERE id IN ({placeholders})
        """
        return execute_query(db_path, sql, tuple(article_ids), fetch=True)
