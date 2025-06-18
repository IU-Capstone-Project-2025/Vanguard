from enum import Enum
from typing import Annotated, List, Literal, Optional, Union

from pydantic import BaseModel, ConfigDict, Field, field_validator


class QuestionType(str, Enum):
    SINGLE_CHOICE = "single_choice"
    TRUE_FALSE = "true_false"
    OPEN_ENDED = "open_ended"


class Option(BaseModel):
    text: str
    image_url: Optional[str] = None
    is_correct: bool


class QuestionBase(BaseModel):
    type: QuestionType
    text: str
    image_url: Optional[str] = None
    time_limit: Annotated[Optional[int], Field(ge=5, le=120)] = None


class SingleChoiceQuestion(QuestionBase):
    type: Literal[QuestionType.SINGLE_CHOICE]
    options: Annotated[List[Option], Field(min_length=3, max_length=4)]

    @field_validator("options")
    @classmethod
    def validate_exactly_one_correct(cls, v: List[Option]) -> List[Option]:
        if sum(opt.is_correct for opt in v) != 1:
            raise ValueError("Exactly one option must be marked as correct.")
        return v


class TrueFalseQuestion(QuestionBase):
    type: Literal[QuestionType.TRUE_FALSE]
    is_true: bool


class OpenEndedQuestion(QuestionBase):
    type: Literal[QuestionType.OPEN_ENDED]
    accepted_answers: Optional[List[str]] = None


Question = Annotated[
    Union[SingleChoiceQuestion, TrueFalseQuestion, OpenEndedQuestion],
    Field(discriminator="type")
]


class QuizCreate(BaseModel):
    model_config = ConfigDict(strict=True)

    title: Annotated[str, Field(max_length=200)]
    description: Optional[str] = None
    is_public: bool = False
    tags: Annotated[List[str], Field(max_length=10)] = []
    questions: Annotated[List[Question], Field(min_length=1, max_length=100)]
