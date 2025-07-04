FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /real-time-service ./cmd/real_time/main.go

FROM alpine:3.19

COPY --from=builder /real-time-service /usr/local/bin/real-time-service

ENTRYPOINT ["/usr/local/bin/real-time-service"]