from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, Field, ConfigDict, EmailStr


class UserBase(BaseModel):
    username: str = Field(..., min_length=3, max_length=50)
    email: EmailStr


class UserCreate(UserBase):
    """Schema for user registration input"""
    password: str = Field(..., min_length=8, max_length=128)


class UserResponse(UserBase):
    """Schema for user response output (omit password hash)"""
    id: UUID
    created_at: datetime

    model_config = ConfigDict(from_attributes=True)


class UserInDB(UserResponse):
    password_hash: str
