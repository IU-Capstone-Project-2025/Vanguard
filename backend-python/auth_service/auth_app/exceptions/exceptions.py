from fastapi import status


class AuthServiceError(Exception):
    """Base exception for auth service."""


class EmailAlreadyExists(AuthServiceError):
    status_code = status.HTTP_400_BAD_REQUEST
    detail = "Email already registered."


class UsernameAlreadyExists(AuthServiceError):
    status_code = status.HTTP_400_BAD_REQUEST
    detail = "Username already registered."


class InvalidCredentials(AuthServiceError):
    status_code = status.HTTP_401_UNAUTHORIZED
    detail = "Invalid email or password."


class UserNotFound(AuthServiceError):
    status_code = status.HTTP_404_NOT_FOUND
    detail = "User not found."


class EmailInUse(AuthServiceError):
    status_code = status.HTTP_403_FORBIDDEN
    detail = "Email is already in use."


class UsernameInUse(AuthServiceError):
    status_code = status.HTTP_403_FORBIDDEN
    detail = "Username is already in use."


class InvalidRefreshToken(AuthServiceError):
    status_code = status.HTTP_401_UNAUTHORIZED
    detail = "Invalid refresh token."


class ExpiredRefreshToken(AuthServiceError):
    status_code = status.HTTP_401_UNAUTHORIZED
    detail = "Refresh token has expired."
