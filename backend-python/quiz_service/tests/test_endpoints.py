from uuid import UUID

import pytest
from fastapi import status


QUIZ_EXAMPLE = {
    "title": "Advanced Python",
    "description": "A quiz about advanced Python topics.",
    "is_public": True,
    "questions": [
        {
            "type": "single_choice",
            "text": "What is the output of print(type(lambda x: x))?",
            "options": [
                {"text": "<class 'function'>", "is_correct": True},
                {"text": "<class 'lambda'>", "is_correct": False},
                {"text": "<class 'method'>", "is_correct": False},
                {"text": "<lambda>", "is_correct": False}
            ]
        }
    ]
}


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


@pytest.mark.asyncio
async def test_create_quiz(test_client_authed, test_user):
    response = await test_client_authed.post(
        "/",
        json=QUIZ_EXAMPLE
    )
    assert response.status_code == status.HTTP_201_CREATED
    data = response.json()
    assert data["title"] == QUIZ_EXAMPLE["title"]
    assert len(data["questions"]) == 1


@pytest.mark.asyncio
async def test_get_quiz_by_id(test_client_authed, test_user, test_quiz):
    response = await test_client_authed.get(
        f"/{test_quiz.id}"
    )
    assert response.status_code == status.HTTP_200_OK
    data = response.json()
    assert data["id"] == str(test_quiz.id)
    assert data["title"] == test_quiz.title


@pytest.mark.asyncio
async def test_update_quiz(test_client_authed, test_user, test_quiz):
    updated_data = {
        "title": "Updated Title",
        "description": "Updated Description",
        "is_public": False
    }
    response = await test_client_authed.put(
        f"/{test_quiz.id}",
        json=updated_data
    )
    assert response.status_code == status.HTTP_200_OK
    data = response.json()
    assert data["title"] == "Updated Title"
    assert data["is_public"] is False


@pytest.mark.asyncio
async def test_delete_quiz(test_client_authed, test_user, test_quiz):
    response = await test_client_authed.delete(
        f"/{test_quiz.id}"
    )
    assert response.status_code == status.HTTP_204_NO_CONTENT

    # Ensure it no longer exists
    get_response = await test_client_authed.get(f"/{test_quiz.id}")
    assert get_response.status_code == status.HTTP_404_NOT_FOUND


@pytest.mark.asyncio
async def test_list_quizzes_all_public(test_client, test_user, test_quiz):
    response = await test_client.get("/")
    assert response.status_code == status.HTTP_200_OK
    quizzes = response.json()
    assert isinstance(quizzes, list)
    assert any(q["id"] == str(test_quiz.id) for q in quizzes)


@pytest.mark.asyncio
async def test_list_quizzes_mine(test_client_authed, test_user, test_quiz):
    response = await test_client_authed.get("/", params={"mine": "true"})
    assert response.status_code == status.HTTP_200_OK
    quizzes = response.json()
    assert all(q["id"] == str(test_quiz.id) for q in quizzes)


@pytest.mark.asyncio
async def test_list_quizzes_filter_by_user_unauthorized(test_client, test_user, test_quiz):
    response = await test_client.get("/", params={"user_id": str(test_user.id), "public": "true"})
    assert response.status_code == status.HTTP_400_BAD_REQUEST


@pytest.mark.asyncio
async def test_list_quizzes_filter_by_user_authorized(test_client_authed, test_user, test_quiz):
    response = await test_client_authed.get("/", params={"user_id": str(test_user.id), "public": "true"})
    assert response.status_code == status.HTTP_200_OK
    quizzes = response.json()
    assert any(q["id"] == str(test_quiz.id) for q in quizzes)


@pytest.mark.asyncio
async def test_list_quizzes_unauthenticated_forbidden_private(test_client):
    response = await test_client.get("/", params={"public": "false"})
    assert response.status_code == status.HTTP_400_BAD_REQUEST


@pytest.mark.asyncio
async def test_list_quizzes_search(test_client, test_user, test_quiz):
    response = await test_client.get("/", params={"search": "Python"})
    assert response.status_code == status.HTTP_200_OK
    quizzes = response.json()
    assert any("Python" in q["title"] for q in quizzes)

@pytest.mark.asyncio
async def test_upload_and_delete_image_success(test_client, test_image):
    files = {"file": ("test.png", test_image, "image/png")}
    response = await test_client.post("/images/upload", files=files)

    assert response.status_code == status.HTTP_200_OK
    assert "url" in response.json()
    assert response.json()["url"].startswith("https://")

    img_url = response.json()["url"]

    response = await test_client.delete(f"/images/?img_url={img_url}")
    assert response.status_code == status.HTTP_204_NO_CONTENT
