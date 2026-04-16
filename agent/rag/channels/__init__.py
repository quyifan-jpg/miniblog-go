"""Search channel implementations — strategy pattern."""

from rag.channels.base import SearchChannel
from rag.channels.chunk_vector_channel import ChunkVectorChannel
from rag.channels.article_vector_channel import ArticleVectorChannel
from rag.channels.keyword_channel import KeywordChannel
from rag.channels.social_media_channel import SocialMediaChannel
from rag.channels.external_search_channel import ExternalSearchChannel

__all__ = [
    "SearchChannel",
    "ChunkVectorChannel",
    "ArticleVectorChannel",
    "KeywordChannel",
    "SocialMediaChannel",
    "ExternalSearchChannel",
]
