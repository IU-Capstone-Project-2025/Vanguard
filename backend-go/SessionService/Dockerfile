# Stage 1: build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем все файлы проекта
COPY .. .

# Собираем бинарь с отключенным CGO и для linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /quiz-engine ./cmd/SessionService/main.go

# Stage 2: финальный образ
FROM alpine:3.19

# Для логов и т.п. добавим bash, ca-certificates
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бинарь из билд-стадии
COPY --from=builder /quiz-engine /usr/local/bin/quiz-engine

# Запускаем наш сервис (передадим флаги через docker-compose)
CMD ["quiz-engine"]
