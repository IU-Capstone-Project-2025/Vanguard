from datetime import datetime
from typing import Optional
from uuid import UUID

from pydantic import BaseModel, ConfigDict, EmailStr, Field, SecretStr


class UserCreate(BaseModel):
    email: EmailStr
    username: str = Field(..., min_length=3, max_length=50)
    password: SecretStr

class UserLogin(BaseModel):
    email: EmailStr
    password: SecretStr

class TokenResponse(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str = "bearer"

class UserResponse(BaseModel):
    id: UUID
    email: EmailStr
    username: str = Field(..., min_length=3, max_length=50)
    created_at: datetime

    model_config = ConfigDict(from_attributes=True)

class UpdateProfile(BaseModel):
    email: Optional[EmailStr] = None
    username: Optional[str] = Field(default=None, min_length=3, max_length=50)
    password: Optional[SecretStr] = None

class RefreshTokenRequest(BaseModel):
    refresh_token: UUID
