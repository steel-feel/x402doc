# Stage 1: Build the Go application
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy dependency definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binaries (pure Go SQLite doesn't require CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrate ./cmd/migrate

# Stage 2: Create a lightweight runner image
FROM alpine:3.19

WORKDIR /app

# Copy built binaries
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/migrate /app/migrate

# Copy migrations folder (required for migrations to run)
COPY --from=builder /app/migrations /app/migrations

# Expose REST (9000) and gRPC (9001) ports
EXPOSE 9000
EXPOSE 9001

# Default command starts the server
CMD ["/app/server"]
