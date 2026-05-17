from pydantic import BaseModel, Field


class PodcastConfigBase(BaseModel):
    name: str
    prompt: str
    description: str | None = None
    time_range_hours: int = Field(24, ge=1, le=168)
    limit_articles: int = Field(20, ge=5, le=50)
    is_active: bool = True
    tts_engine: str = "kokoro"
    language_code: str = "en"
    podcast_script_prompt: str | None = None
    image_prompt: str | None = None


class PodcastConfig(PodcastConfigBase):
    id: int
    created_at: str | None = None
    updated_at: str | None = None

    class Config:
        from_attributes = True


class PodcastConfigCreate(PodcastConfigBase):
    pass


class PodcastConfigUpdate(BaseModel):
    name: str | None = None
    prompt: str | None = None
    description: str | None = None
    time_range_hours: int | None = Field(None, ge=1, le=168)
    limit_articles: int | None = Field(None, ge=5, le=50)
    is_active: bool | None = None
    tts_engine: str | None = None
    language_code: str | None = None
    podcast_script_prompt: str | None = None
    image_prompt: str | None = None
