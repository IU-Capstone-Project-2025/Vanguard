from shared.db.models import User as DBUser
from shared.repositories.base_repo import SQLAlchemyRepository


class UserRepository(SQLAlchemyRepository[DBUser]):
    model = DBUser