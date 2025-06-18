from datetime import datetime
from typing import List, Optional
from uuid import UUID, uuid4

from sqlalchemy import ForeignKey, String, Boolean, Text, Integer, DateTime
from sqlalchemy.dialects.postgresql import UUID as PG_UUID, JSONB
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship
from sqlalchemy.sql import func


class Base(DeclarativeBase):
    pass


class User(Base):
    __tablename__ = "users"

    id: Mapped[UUID] = mapped_column(
        PG_UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
        server_default=func.gen_random_uuid()
    )
    username: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    email: Mapped[str] = mapped_column(String(100), unique=True, nullable=False)
    password_hash: Mapped[str] = mapped_column(Text, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())

    # Relationship to quizzes
    quizzes: Mapped[List["Quiz"]] = relationship(
        back_populates="owner",
        cascade="all, delete-orphan"
    )


class Quiz(Base):
    __tablename__ = "quizzes"

    id: Mapped[UUID] = mapped_column(
        PG_UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
        server_default=func.gen_random_uuid()
    )
    title: Mapped[str] = mapped_column(String(200), nullable=False)
    description: Mapped[Optional[str]] = mapped_column(Text)
    is_public: Mapped[bool] = mapped_column(Boolean, default=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

    # Questions stored as JSONB
    questions: Mapped[list] = mapped_column(JSONB, nullable=False)

    # Owner relationship
    owner_id: Mapped[UUID] = mapped_column(
        PG_UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE")
    )
    owner: Mapped["User"] = relationship(back_populates="quizzes")

    # Tags relationship (many-to-many)
    # Deleting a Quiz will remove its entries in quiz_tags (via FK), but not delete the tags themselves
    tags: Mapped[List["Tag"]] = relationship(
        secondary="quiz_tags",
        back_populates="quizzes",
        passive_deletes=True
    )


class Tag(Base):
    __tablename__ = "tags"

    id: Mapped[int] = mapped_column(
        Integer,
        primary_key=True,
        autoincrement=True
    )
    name: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)

    # Relationship to quizzes
    # Deleting a Tag will remove its entries in quiz_tags (via FK), but not delete any quizzes
    quizzes: Mapped[List["Quiz"]] = relationship(
        secondary="quiz_tags",
        back_populates="tags",
        passive_deletes=True
    )


class QuizTag(Base):
    __tablename__ = "quiz_tags"

    quiz_id: Mapped[UUID] = mapped_column(
        PG_UUID(as_uuid=True),
        ForeignKey("quizzes.id", ondelete="CASCADE"),  # When a Quiz is deleted, all its quiz_tags rows are auto-deleted
        primary_key=True
    )
    tag_id: Mapped[int] = mapped_column(
        Integer,
        ForeignKey("tags.id", ondelete="CASCADE"),  # When a Tag is deleted, all its quiz_tags rows are auto-deleted
        primary_key=True
    )
