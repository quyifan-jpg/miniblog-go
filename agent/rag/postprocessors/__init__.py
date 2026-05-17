"""Post-processor implementations — chain of responsibility pattern."""

from rag.postprocessors.base import PostProcessor
from rag.postprocessors.dedup import DeduplicationProcessor
from rag.postprocessors.quality_filter import QualityFilterProcessor
from rag.postprocessors.rerank import RerankProcessor
from rag.postprocessors.rrf import RRFProcessor

__all__ = [
    "PostProcessor",
    "DeduplicationProcessor",
    "RRFProcessor",
    "RerankProcessor",
    "QualityFilterProcessor",
]
