# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and .env
COPY . .
COPY .env .env

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o matching-service ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Add non-root user
RUN adduser -D appuser

# Copy the binary from builder
COPY --from=builder /app/matching-service .

# Copy the .env file from builder and set permissions
COPY --from=builder /app/.env .env
RUN chown appuser:appuser .env && chmod 644 .env

# Switch to non-root user
USER appuser

EXPOSE 8080
CMD ["./matching-service"]

