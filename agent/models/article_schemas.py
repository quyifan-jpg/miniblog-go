from datetime import datetime
from typing import Any

from pydantic import BaseModel, ConfigDict, field_validator


class ArticleBase(BaseModel):
    title: str
    url: str | None = None
    published_date: str

    @field_validator("published_date", mode="before")
    @classmethod
    def coerce_datetime(cls, v):
        if isinstance(v, datetime):
            return v.isoformat()
        return v

    summary: str | None = None
    content: str | None = None
    categories: list[str] | None = []
    source_name: str | None = None


class Article(ArticleBase):
    id: int
    metadata: dict[str, Any] | None = {}

    model_config = ConfigDict(from_attributes=True)


class PaginatedArticles(BaseModel):
    items: list[Article]
    total: int
    page: int
    per_page: int
    total_pages: int
    has_next: bool
    has_prev: bool
