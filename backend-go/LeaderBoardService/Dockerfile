FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /Leader-Board ./cmd/LeaderBoardService/main.go

FROM alpine:3.19

COPY --from=builder /Leader-Board /usr/local/bin/Leader-Board

ENTRYPOINT ["/usr/local/bin/Leader-Board"]