FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

COPY auth_service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY shared /shared
RUN pip install /shared

COPY auth_service/auth_app ./auth_app

EXPOSE 8000

CMD ["uvicorn", "auth_app.main:app", "--host", "0.0.0.0", "--port", "8000", "--no-access-log"]