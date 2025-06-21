FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

COPY quiz_service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY shared /shared
RUN pip install /shared

COPY quiz_service/quiz_app ./quiz_app

EXPOSE 8000

CMD ["uvicorn", "quiz_app.main:app", "--host", "0.0.0.0", "--port", "8000"]