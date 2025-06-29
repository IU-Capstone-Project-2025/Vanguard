from datetime import datetime, timedelta, timezone

import pytest

from shared.db.models import User, RefreshToken
from shared.repositories import UserRepository, RefreshTokenRepository

@pytest.mark.asyncio
async def test_register_and_me_and_update(test_client):
    # Register
    resp = await test_client.post("/register", json={
        "email": "bob@example.com",
        "username": "bob",
        "password": "hunter2"
    })
    assert resp.status_code == 201
    bob = resp.json()
    assert "id" in bob and bob["email"] == "bob@example.com"

    # Login
    resp = await test_client.post("/login", json={
        "email": "bob@example.com",
        "password": "hunter2"
    })
    assert resp.status_code == 200
    tokens = resp.json()
    assert "access_token" in tokens and "refresh_token" in tokens

    headers = {"Authorization": f"Bearer {tokens['access_token']}"}

    # Me
    resp = await test_client.get("/me", headers=headers)
    assert resp.status_code == 200
    me = resp.json()
    assert me["email"] == "bob@example.com"

    # Update profile (change email)
    resp = await test_client.put("/me", json={"email": "bob2@example.com"}, headers=headers)
    assert resp.status_code == 200
    updated = resp.json()
    assert updated["email"] == "bob2@example.com"

@pytest.mark.asyncio
async def test_register_conflicts(test_client, registered_user):
    # Same email
    resp = await test_client.post("/register", json={
        "email": "alice@example.com",
        "username": "newalice",
        "password": "pw"
    })
    assert resp.status_code == 400

    # Same username
    resp = await test_client.post("/register", json={
        "email": "other@example.com",
        "username": "alice",
        "password": "pw"
    })
    assert resp.status_code == 400

@pytest.mark.asyncio
async def test_login_failures(test_client):
    # No such user
    resp = await test_client.post("/login", json={
        "email": "nouser@example.com",
        "password": "pw"
    })
    assert resp.status_code == 401

    # Wrong password
    await test_client.post("/register", json={
        "email": "carl@example.com",
        "username": "carl",
        "password": "rightpw"
    })
    resp = await test_client.post("/login", json={
        "email": "carl@example.com",
        "password": "wrongpw"
    })
    assert resp.status_code == 401

@pytest.mark.asyncio
async def test_refresh_and_logout(test_client):
    # Register + login
    await test_client.post("/register", json={
        "email": "dan@example.com",
        "username": "dan",
        "password": "pw"
    })
    login = (await test_client.post("/login", json={
        "email": "dan@example.com",
        "password": "pw"
    })).json()

    # Use refresh
    resp = await test_client.post("/refresh", json={"refresh_token": login["refresh_token"]})
    assert resp.status_code == 200
    new_tokens = resp.json()
    assert new_tokens["access_token"] != login["access_token"]

    # Logout old token
    resp = await test_client.post("/logout", json={"refresh_token": login["refresh_token"]})
    assert resp.status_code == 204

    # Using old refresh should now fail
    resp = await test_client.post("/refresh", json={"refresh_token": login["refresh_token"]})
    assert resp.status_code == 401

@pytest.mark.asyncio
async def test_refresh_expired(test_client, uow_test):
    # Directly insert expired token
    async with uow_test.transaction() as sess:
        user_repo = UserRepository(sess)
        user = await user_repo.create(User(email="eve@example.com", username="eve", password_hash="hash"))
        expire = datetime.now(timezone.utc) - timedelta(days=1)
        rt_repo = RefreshTokenRepository(sess)
        rt = await rt_repo.create(RefreshToken(
            user_id=user.id,
            user_agent="x",
            ip_address="127.0.0.1",
            expires_at=expire
        ))
    # Attempt refresh
    resp = await test_client.post("/refresh", json={"refresh_token": str(rt.token)})
    assert resp.status_code == 401

@pytest.mark.asyncio
async def test_logout_all(test_client):
    # Register + login twice
    await test_client.post("/register", json={
        "email": "frank@example.com",
        "username": "frank",
        "password": "pw"
    })
    first = (await test_client.post("/login", json={
        "email": "frank@example.com",
        "password": "pw"
    })).json()
    second = (await test_client.post("/login", json={
        "email": "frank@example.com",
        "password": "pw"
    })).json()

    # Logout all
    headers = {"Authorization": f"Bearer {first['access_token']}"}
    resp = await test_client.post("/logout/all", headers=headers)
    assert resp.status_code == 204

    # Both refresh tokens are now invalid
    for tok in (first["refresh_token"], second["refresh_token"]):
        r = await test_client.post("/refresh", json={"refresh_token": tok})
        assert r.status_code == 401
