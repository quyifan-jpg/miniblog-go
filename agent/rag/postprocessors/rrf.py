"""
Reciprocal Rank Fusion (RRF) post-processor.

Order 5 — runs after dedup, before rerank.

Equivalent to CAgent's RRFPostProcessor.

RRF solves the problem of combining ranked lists from different channels
where the raw scores are NOT comparable:
  - ChunkVector scores are cosine similarities in [0, 1]
  - Keyword scores are binary (1.0 or 0.8)
  - Social media scores are normalized engagement counts
  - External scores are position-based decrements

RRF only uses RANK (position), not raw scores:
  RRF_score(doc) = Σ  1 / (k + rank_in_channel_i)

Where:
  - k = smoothing parameter (default 60, from Cormack et al. 2009)
  - rank starts at 1 (top result)
  - The sum is over all channels that returned this document

Properties:
  - Documents ranked highly by MULTIPLE channels get boosted
  - Documents ranked highly by ONE channel still get a reasonable score
  - Immune to score-scale differences between channels
"""

from __future__ import annotations

from loguru import logger

from rag.models import ChannelResult, RetrievedChunk, SearchContext
from rag.postprocessors.base import PostProcessor


class RRFProcessor(PostProcessor):
    """Reciprocal Rank Fusion — CAgent-equivalent implementation."""

    def __init__(self, k: int = 60):
        """
        Args:
            k: Smoothing parameter. Higher k → more emphasis on being
               returned by multiple channels. Lower k → more emphasis
               on being ranked high in a single channel.
               Default 60 (Cormack et al. recommendation).
        """
        self._k = k

    def order(self) -> int:
        return 5

    def is_enabled(self, context: SearchContext) -> bool:
        # RRF is only meaningful with 2+ channels
        return True  # let process() handle the <2 case

    def process(
        self,
        chunks: list[RetrievedChunk],
        all_results: list[ChannelResult],
        context: SearchContext,
    ) -> list[RetrievedChunk]:
        if not chunks:
            return []

        # Only fuse if we have results from 2+ channels
        effective_results = [r for r in all_results if r.chunks]
        if len(effective_results) < 2:
            logger.debug(
                "RRF skipped: only {n} effective channel(s)",
                n=len(effective_results),
            )
            return chunks

        # Build RRF scores
        rrf_scores: dict[str, float] = {}
        chunk_map: dict[str, RetrievedChunk] = {}

        for result in effective_results:
            for rank, chunk in enumerate(result.chunks, start=1):
                key = chunk.dedup_key
                rrf_scores[key] = rrf_scores.get(key, 0.0) + 1.0 / (self._k + rank)

                # Keep the version with the best original score
                if key not in chunk_map or chunk.score > chunk_map[key].score:
                    chunk_map[key] = chunk

        if not rrf_scores:
            return chunks

        # Normalize RRF scores to [0, 1]
        max_rrf = max(rrf_scores.values())
        if max_rrf <= 0:
            return chunks

        # Update scores and sort
        fused = []
        for key, chunk in chunk_map.items():
            chunk.score = rrf_scores[key] / max_rrf
            fused.append(chunk)

        fused.sort(key=lambda c: c.score, reverse=True)

        logger.debug(
            "RRF fused {n} chunks from {c} channels (k={k})",
            n=len(fused),
            c=len(effective_results),
            k=self._k,
        )
        return fused
