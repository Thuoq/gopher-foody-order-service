# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Development stage with Air
FROM golang:1.26-alpine AS dev

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Default command for development
CMD ["air", "-c", ".air.toml"]

# Final production stage
FROM alpine:latest AS production

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

EXPOSE 8082 9092

CMD ["./main"]
