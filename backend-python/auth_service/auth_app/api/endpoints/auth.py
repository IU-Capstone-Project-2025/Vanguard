from uuid import UUID

from fastapi import APIRouter, Body, Depends, Header, Request, status

from shared.schemas.auth import (
    UserCreate, UserLogin, TokenResponse, UserResponse, UpdateProfile, RefreshTokenRequest
)

from shared.core.dependencies import get_current_user_id

from auth_app.core.dependencies import get_auth_service
from auth_app.services.auth_service import AuthService

router = APIRouter(
    prefix="/api",
    tags=["auth-service"],
    responses={
        401: {"description": "Unauthorized - Invalid or missing credentials"},
        500: {"description": "Internal Server Error"}
    }
)


@router.post(
    "/register",
    summary="Register New User",
    description="Creates a new user account. Email and username must be unique.",
    response_model=UserResponse,
    status_code=status.HTTP_201_CREATED,
    response_description="Newly created user record (password excluded)",
    responses={
        400: {"description": "Email/username already exists"},
        422: {"description": "Validation error in request body"}
    }
)
async def register(
        data: UserCreate = Body(..., description="User registration data"),
        svc: AuthService = Depends(get_auth_service)
):
    return await svc.register(data)


@router.post(
    "/login",
    summary="User Login",
    description="Authenticates a user and returns JWT tokens.",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
    response_description="Access and refresh tokens",
    responses={
        401: {"description": "Invalid credentials"},
        422: {"description": "Validation error in request body"}
    }
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
    description="Exchanges a valid refresh token for a new access and refresh token pair.",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
    response_description="New access and refresh tokens",
    responses={
        401: {"description": "Invalid or expired refresh token"},
        422: {"description": "Validation error in request body"}
    }
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
    description="Returns the authenticated user's profile data based on the access token.",
    response_model=UserResponse,
    status_code=status.HTTP_200_OK,
    response_description="Profile data for the authenticated user",
    responses={
        401: {"description": "Missing or invalid JWT"},
        404: {"description": "User extracted from access token not found"}
    }
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
    description="Updates the authenticated user's profile. Email and username must remain unique.",
    response_model=UserResponse,
    status_code=status.HTTP_200_OK,
    response_description="Updated profile data",
    responses={
        400: {"description": "Invalid input"},
        401: {"description": "Missing or invalid JWT"},
        403: {"description": "New email/username already in use"},
        404: {"description": "User extracted from access token not found"},
        422: {"description": "Validation error in request body"}
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
    description="Revokes a single refresh token to log the user out of one device/session.",
    status_code=status.HTTP_204_NO_CONTENT,
    responses={
        204: {"description": "Refresh token revoked"},
        422: {"description": "Validation error in request body"}
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
    description="Revokes all refresh tokens associated with the authenticated user, logging them out everywhere.",
    status_code=status.HTTP_204_NO_CONTENT,
    responses={
        204: {"description": "All refresh tokens revoked for user"},
        401: {"description": "Missing or invalid JWT"}
    }
)
async def logout_all(
        _authorization: str = Header(..., alias="Authorization", description="Bearer <access_token>"),
        user_id: UUID = Depends(get_current_user_id),
        svc: AuthService = Depends(get_auth_service)
):
    await svc.logout_all(user_id)
