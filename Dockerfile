# Build stage
FROM golang:1.24-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM base AS builder

WORKDIR /app

COPY . .

RUN go build -o /todo-app ./cmd/app/main.go

# Run stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /todo-app /app/todo-app
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8000

CMD ["/app/todo-app"]