
# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . . 

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o chatting-service ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Add non-root user
RUN adduser -D appuser

# Copy the binary from builder
COPY --from=builder /app/chatting-service .

# Switch to non-root user
USER appuser

EXPOSE 8080
CMD ["./chatting-service"]

