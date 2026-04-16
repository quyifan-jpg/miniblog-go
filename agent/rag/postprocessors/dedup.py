"""
Deduplication post-processor.

Order 1 — runs first. Removes duplicate documents that appear across
multiple channels.

Equivalent to CAgent's DeduplicationPostProcessor:
- Results from higher-priority channels take precedence
- When duplicates exist, the version with the higher score is kept
- Uses URL as the dedup key (falls back to chunk.id)

Example:
  ChunkVector returns article_42 with score=0.85
  Keyword returns article_42 with score=1.0
  → ChunkVector has priority=1, Keyword has priority=3
  → Keep ChunkVector's version (higher priority channel)
  → But use max(0.85, 1.0) = 1.0 as the score
"""

from __future__ import annotations

from rag.models import CHANNEL_PRIORITY, ChannelResult, RetrievedChunk, SearchContext
from rag.postprocessors.base import PostProcessor


class DeduplicationProcessor(PostProcessor):
    """Priority-aware deduplication across channels."""

    def order(self) -> int:
        return 1

    def process(
        self,
        chunks: list[RetrievedChunk],
        all_results: list[ChannelResult],
        context: SearchContext,
    ) -> list[RetrievedChunk]:
        if not chunks:
            return []

        # Sort channel results by priority (lower = higher precedence)
        sorted_results = sorted(
            all_results,
            key=lambda r: CHANNEL_PRIORITY.get(r.channel_type, 99),
        )

        # Walk through results in priority order
        # For each dedup_key: keep the first (highest priority) channel's version,
        # but take the max score across all appearances
        seen: dict[str, RetrievedChunk] = {}
        score_max: dict[str, float] = {}

        for result in sorted_results:
            for chunk in result.chunks:
                key = chunk.dedup_key
                if key not in seen:
                    # First time seeing this document — keep it
                    seen[key] = chunk
                    score_max[key] = chunk.score
                else:
                    # Duplicate — keep higher score but preserve the
                    # higher-priority channel's metadata/content
                    score_max[key] = max(score_max[key], chunk.score)

        # Apply max scores and return in insertion order (priority order)
        deduped = []
        for key, chunk in seen.items():
            chunk.score = score_max[key]
            deduped.append(chunk)

        return deduped
