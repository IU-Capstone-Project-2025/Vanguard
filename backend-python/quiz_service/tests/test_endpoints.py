from uuid import UUID

import pytest


@pytest.mark.asyncio
async def test_get_public_quizzes(test_client, test_quiz):
    response = await test_client.get("/")

    assert response.status_code == 200

    quizzes = response.json()
    assert isinstance(quizzes, list)
    assert len(quizzes) >= 1  # At least the test_quiz

    # Check that only public quizzes are included
    for quiz in quizzes:
        assert quiz["is_public"] is True
        assert "id" in quiz
        assert "title" in quiz
        assert isinstance(UUID(quiz["id"]), UUID)
        assert isinstance(quiz["questions"], list)
        assert len(quiz["questions"]) > 0


@pytest.mark.asyncio
async def test_get_quiz_by_id_as_unauthenticated(test_client, test_quiz):
    response = await test_client.get(f"/{test_quiz.id}")

    assert response.status_code == 200
    data = response.json()

    assert data["id"] == str(test_quiz.id)
    assert data["title"] == test_quiz.title
    assert data["is_public"] is True
    assert "questions" in data
    assert len(data["questions"]) == 3
