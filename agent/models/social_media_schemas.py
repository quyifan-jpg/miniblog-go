from datetime import datetime

from pydantic import BaseModel


class PostBase(BaseModel):
    post_id: str
    platform: str
    user_display_name: str | None = None
    user_handle: str | None = None
    user_profile_pic_url: str | None = None
    post_timestamp: str | None = None
    post_display_time: str | None = None
    post_url: str | None = None
    post_text: str | None = None
    post_mentions: str | None = None


class PostEngagement(BaseModel):
    replies: int | None = None
    retweets: int | None = None
    likes: int | None = None
    bookmarks: int | None = None
    views: int | None = None


class MediaItem(BaseModel):
    type: str
    url: str


class Post(PostBase):
    engagement: PostEngagement | None = None
    media: list[MediaItem] | None = None
    media_count: int | None = 0
    is_ad: bool | None = False
    sentiment: str | None = None
    categories: list[str] | None = None
    tags: list[str] | None = None
    analysis_reasoning: str | None = None
    created_at: datetime | None = None
    updated_at: datetime | None = None


class PaginatedPosts(BaseModel):
    items: list[Post]
    total: int
    page: int
    per_page: int
    total_pages: int
    has_next: bool
    has_prev: bool


class PostFilterParams(BaseModel):
    platform: str | None = None
    user_handle: str | None = None
    sentiment: str | None = None
    category: str | None = None
    date_from: str | None = None
    date_to: str | None = None
    search: str | None = None
