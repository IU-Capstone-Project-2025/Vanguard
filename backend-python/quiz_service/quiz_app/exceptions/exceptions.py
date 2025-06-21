from uuid import UUID


class NotFoundError(Exception):
    """Base exception for when a resource is not found"""
    def __init__(self, resource: str, identifier: str | int | UUID):
        self.resource = resource
        self.identifier = identifier
        super().__init__(f"{resource} with id '{identifier}' not found")


class ForbiddenError(Exception):
    """Base exception for when access is forbidden"""
    def __init__(self, message: str = "Action is forbidden"):
        super().__init__(message)


class QuizNotFoundError(NotFoundError):
    """Raised when a quiz is not found"""
    def __init__(self, identifier: str | int | UUID):
        super().__init__("Quiz", identifier)
