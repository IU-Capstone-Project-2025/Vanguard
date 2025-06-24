class TokenError(Exception):
    """Base exception for token related errors"""

    def __init__(self, message: str = "Token authentication failed"):
        self.message = message
        super().__init__(message)

    def __str__(self):
        return self.message

class TokenExpiredError(TokenError):
    """Raised when token has expired"""

    def __init__(self):
        super().__init__("Token has expired")

class InvalidTokenError(TokenError):
    """Raised when token is invalid"""

    def __init__(self):
        super().__init__("Invalid token")

class MissingTokenError(TokenError):
    """Raised when token is missing"""

    def __init__(self):
        super().__init__("Authorization token is missing")

class InvalidTokenPayloadError(TokenError):
    """Raised when token payload is invalid"""
    def __init__(self, field: str):
        self.field = field
        message = f"Missing required field in token: '{field}'"
        super().__init__(message)
