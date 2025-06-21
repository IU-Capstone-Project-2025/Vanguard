FROM python:3.12-slim

WORKDIR /app

RUN pip install alembic==1.14.0 asyncpg==0.30.0 pydantic==2.10.5 pydantic-settings==2.7.1 SQLAlchemy==2.0.37

COPY . .

CMD ["alembic", "upgrade", "head"]