
# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o matching-service ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/matching-service .

# Add non-root user (optional but secure)
RUN adduser -D appuser
USER appuser

# Expose port (if needed for metrics or HTTP)
EXPOSE 8080

# Run the service
CMD ["./matching-service"]
