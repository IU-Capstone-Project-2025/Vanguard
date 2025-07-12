from uuid import UUID

from fastapi import APIRouter, Body, Depends, Header, Request, status

from shared.schemas.auth import (
    UserCreate, UserLogin, TokenResponse, UserResponse, UpdateProfile, RefreshTokenRequest
)

from shared.core.dependencies import get_current_user_id

from auth_app.core.dependencies import get_auth_service
from auth_app.services.auth_service import AuthService

router = APIRouter(tags=["auth-service"])


@router.post(
    "/register",
    summary="Register New User",
    response_model=UserResponse,
    status_code=status.HTTP_201_CREATED,
    response_description="Newly created user record (without password)"
)
async def register(
        data: UserCreate = Body(..., description="User registration data"),
        svc: AuthService = Depends(get_auth_service)
):
    return await svc.register(data)


@router.post(
    "/login",
    summary="User Login",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
    response_description="Access and refresh tokens"
)
async def login(
        request: Request,
        data: UserLogin = Body(..., description="User login credentials"),
        svc: AuthService = Depends(get_auth_service)
):
    return await svc.login(data, request)



@router.post(
    "/refresh",
    summary="Refresh Access Token",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
    response_description="New access and refresh tokens"
)
async def refresh(
        request: Request,
        req: RefreshTokenRequest = Body(..., description="Refresh token to be exchanged"),
        svc: AuthService = Depends(get_auth_service)
):
    return await svc.refresh(req, request)


@router.get(
    "/me",
    summary="Get Current User",
    response_model=UserResponse,
    status_code=status.HTTP_200_OK,
    response_description="Profile data for the authenticated user",
    responses={401: {"description": "Missing or invalid JWT"}}
)
async def me(
        _authorization: str = Header(..., alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        svc: AuthService = Depends(get_auth_service)
):
    return await svc.me(user_id)


@router.put(
    "/me",
    summary="Update Current User",
    response_model=UserResponse,
    status_code=status.HTTP_200_OK,
    response_description="Updated profile data",
    responses={
        400: {"description": "Invalid input"},
        401: {"description": "Missing or invalid JWT"},
        403: {"description": "New email/username already in use"}
    }
)
async def update_me(
        data: UpdateProfile = Body(..., description="Fields to update (email, username, or password)"),
        _authorization: str = Header(..., alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        svc: AuthService = Depends(get_auth_service)
):
    return await svc.update_me(user_id, data)



@router.post(
    "/logout",
    summary="Logout (Revoke One Refresh Token)",
    status_code=status.HTTP_204_NO_CONTENT,
    responses={
        204: {"description": "Refresh token revoked"},
        400: {"description": "Invalid token format"},
    }
)
async def logout(
        body: RefreshTokenRequest = Body(..., description="Refresh token to revoke"),
        svc: AuthService = Depends(get_auth_service)
):
    await svc.logout(body.refresh_token)


@router.post(
    "/logout/all",
    summary="Logout All Sessions",
    status_code=status.HTTP_204_NO_CONTENT,
    responses={
        204: {"description": "All refresh tokens revoked for user"},
        401: {"description": "Missing or invalid JWT"},
    }
)
async def logout_all(
        _authorization: str = Header(..., alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        svc: AuthService = Depends(get_auth_service)
):
    await svc.logout_all(user_id)
