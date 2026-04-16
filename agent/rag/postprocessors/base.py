"""
Post-processor interface — the chain of responsibility contract.

Equivalent to CAgent's SearchResultPostProcessor interface.

Post-processors execute sequentially in ascending order:
  dedup (1) → RRF (5) → rerank (10) → quality_filter (20)

Each processor receives:
  - chunks: current list of RetrievedChunks (output of previous processor)
  - all_results: original per-channel results (for RRF rank calculation)
  - context: shared SearchContext
"""

from __future__ import annotations

from abc import ABC, abstractmethod

from rag.models import ChannelResult, RetrievedChunk, SearchContext


class PostProcessor(ABC):
    """Abstract base for all post-processors."""

    @abstractmethod
    def order(self) -> int:
        """
        Execution order — lower numbers run first.

        Convention:
          1-4:   deduplication
          5-9:   rank fusion
          10-14: reranking
          15-19: (reserved)
          20+:   quality filtering / truncation
        """
        ...

    def is_enabled(self, context: SearchContext) -> bool:
        """
        Runtime gate — return False to skip this processor.

        Default: always enabled.
        """
        return True

    @abstractmethod
    def process(
        self,
        chunks: list[RetrievedChunk],
        all_results: list[ChannelResult],
        context: SearchContext,
    ) -> list[RetrievedChunk]:
        """
        Process the chunk list and return a (possibly filtered/reordered) list.

        Implementations MUST NOT mutate the input list — create a new list.
        They MAY mutate individual RetrievedChunk.score values.
        """
        ...
