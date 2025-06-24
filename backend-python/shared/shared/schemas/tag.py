import re
from typing import Annotated

from pydantic import BaseModel, ConfigDict, Field, ValidationError, field_validator

TAG_NAME_REGEX = re.compile(r"^[a-z0-9_-]{1,50}$")


class TagBase(BaseModel):
    model_config = ConfigDict(strict=True)

    name: Annotated[
        str,
        Field(
            max_length=50,
            description="Lowercase letters, numbers, hyphens or underscores only"
        )
    ]

    @field_validator("name")
    @classmethod
    def validate_and_normalize(cls, v: str) -> str:
        normalized = v.strip().lower()
        if not TAG_NAME_REGEX.match(normalized):
            raise ValidationError("Tag must be 1–50 chars of [a–z0–9_-] only")
        return normalized


class TagResponse(TagBase):
    id: int

    model_config = ConfigDict(from_attributes=True)
