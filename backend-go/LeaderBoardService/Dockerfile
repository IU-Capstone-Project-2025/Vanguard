FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /leaderboard ./cmd/LeaderBoard/main.go

FROM alpine:3.19

COPY --from=builder /leaderboard /usr/local/bin/leaderboard

ENTRYPOINT ["/usr/local/bin/leaderboard"]