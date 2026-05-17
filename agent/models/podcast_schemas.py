from typing import Any

from pydantic import BaseModel


class PodcastBase(BaseModel):
    title: str
    date: str
    audio_generated: bool = False
    banner_img: str | None = None
    identifier: str
    language_code: str | None = "en"
    tts_engine: str | None = "kokoro"


class Podcast(PodcastBase):
    id: int
    created_at: str | None = None
    audio_path: str | None = None

    class Config:
        from_attributes = True


class PodcastContent(BaseModel):
    title: str
    sections: list[dict[str, Any]]


class PodcastSource(BaseModel):
    title: str | None = None
    url: str | None = None
    source: str | None = None

    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        if isinstance(v, str):
            return cls(url=v)
        if isinstance(v, dict):
            return cls(**v)
        raise ValueError("Source must be a string or a dict")


class PodcastDetail(BaseModel):
    podcast: Podcast
    content: PodcastContent
    audio_url: str | None = None
    sources: list[PodcastSource | str] | None = None
    banner_images: list[str] | None = None


class PodcastCreate(BaseModel):
    title: str
    date: str | None = None
    content: dict[str, Any]
    sources: list[dict[str, str] | str] | None = None
    language_code: str | None = "en"
    tts_engine: str | None = "kokoro"


class PodcastUpdate(BaseModel):
    title: str | None = None
    date: str | None = None
    content: dict[str, Any] | None = None
    audio_generated: bool | None = None
    sources: list[dict[str, str] | str] | None = None
    language_code: str | None = None
    tts_engine: str | None = None


class PaginatedPodcasts(BaseModel):
    items: list[Podcast]
    total: int
    page: int
    per_page: int
    total_pages: int
    has_next: bool
    has_prev: bool
