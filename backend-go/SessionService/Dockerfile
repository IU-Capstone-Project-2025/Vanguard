# Сборка бинарника
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /quiz-engine ./cmd/main.go

# Минимальный рантайм
FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /quiz-engine /usr/local/bin/quiz-engine

CMD ["quiz-engine"]
