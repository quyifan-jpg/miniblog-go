"""Post-processor implementations — chain of responsibility pattern."""

from rag.postprocessors.base import PostProcessor
from rag.postprocessors.dedup import DeduplicationProcessor
from rag.postprocessors.rrf import RRFProcessor
from rag.postprocessors.rerank import RerankProcessor
from rag.postprocessors.quality_filter import QualityFilterProcessor

__all__ = [
    "PostProcessor",
    "DeduplicationProcessor",
    "RRFProcessor",
    "RerankProcessor",
    "QualityFilterProcessor",
]
