from pydantic import BaseModel


class SourceFeed(BaseModel):
    id: int
    feed_url: str
    feed_type: str
    description: str | None = None
    is_active: bool
    created_at: str
    last_crawled: str | None = None


class SourceFeedCreate(BaseModel):
    feed_url: str
    feed_type: str = "main"
    description: str | None = None
    is_active: bool = True


class SourceBase(BaseModel):
    name: str
    url: str | None = None
    categories: list[str] | None = []
    description: str | None = None
    is_active: bool = True


class Source(SourceBase):
    id: int
    created_at: str | None = None
    last_crawled: str | None = None

    class Config:
        from_attributes = True


class SourceCreate(SourceBase):
    feeds: list[SourceFeedCreate] | None = []


class SourceUpdate(BaseModel):
    name: str | None = None
    url: str | None = None
    categories: list[str] | None = None
    description: str | None = None
    is_active: bool | None = None


class SourceWithFeeds(Source):
    feeds: list[SourceFeed] = []


class PaginatedSources(BaseModel):
    items: list[Source]
    total: int
    page: int
    per_page: int
    total_pages: int
    has_next: bool
    has_prev: bool


class Category(BaseModel):
    id: int
    name: str
    description: str | None = None
