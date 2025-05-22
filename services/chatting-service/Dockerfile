# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /chatting_server

# Copy go.mod and go.sum from the chatting-service directory
COPY services/chatting-service/go.mod services/chatting-service/go.sum ./

RUN go mod download

# Copy all the service code for the chatting service
COPY services/chatting-service/. .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o chatting-service ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Add non-root user
RUN adduser -D appuser

# Copy the binary from the builder
COPY --from=builder /chatting_server/chatting-service .

# Switch to non-root user
USER appuser

EXPOSE 8080
CMD ["./chatting-service"]

