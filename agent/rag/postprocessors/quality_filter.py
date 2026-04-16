"""
Quality filter post-processor.

Order 20 — runs last. Filters out low-quality results based on:
- Minimum relevance score
- Minimum content length

Equivalent to CAgent's QualityFilterPostProcessor.

Safety: if all chunks are filtered out, keeps the single best chunk
to avoid returning empty results (configurable).
"""

from __future__ import annotations

from loguru import logger

from rag.models import ChannelResult, RetrievedChunk, SearchContext
from rag.postprocessors.base import PostProcessor


class QualityFilterProcessor(PostProcessor):
    """Score and length-based quality gate."""

    def __init__(
        self,
        *,
        min_score: float = 0.3,
        min_content_length: int = 10,
        keep_best_when_all_filtered: bool = True,
    ):
        """
        Args:
            min_score: Chunks with score below this are dropped.
            min_content_length: Chunks with content shorter than this
                                (in characters) are dropped.
            keep_best_when_all_filtered: If True and ALL chunks are
                                         filtered out, keep the one
                                         with the highest score.
        """
        self._min_score = min_score
        self._min_content_length = min_content_length
        self._keep_best = keep_best_when_all_filtered

    def order(self) -> int:
        return 20

    def process(
        self,
        chunks: list[RetrievedChunk],
        all_results: list[ChannelResult],
        context: SearchContext,
    ) -> list[RetrievedChunk]:
        if not chunks:
            return []

        filtered = [
            c for c in chunks
            if c.score >= self._min_score
            and len(c.content) >= self._min_content_length
        ]

        dropped = len(chunks) - len(filtered)
        if dropped > 0:
            logger.info(
                "QualityFilter: dropped {n} chunks "
                "(score < {s} or length < {l})",
                n=dropped,
                s=self._min_score,
                l=self._min_content_length,
            )

        if not filtered and self._keep_best:
            # All filtered out — keep the best one as a safety net
            best = max(chunks, key=lambda c: c.score)
            logger.warning(
                "QualityFilter: all chunks filtered, keeping best "
                "(score={s:.3f}, id={id})",
                s=best.score,
                id=best.id,
            )
            return [best]

        # Final truncation to top_k
        return filtered[: context.top_k]
