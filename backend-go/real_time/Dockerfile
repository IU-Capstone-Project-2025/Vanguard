FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /real-time-service ./cmd/real_time/main.go

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /real-time-service /usr/local/bin/real-time-service

CMD ["real-time-service"]
