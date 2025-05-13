# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from the matching-service directory
COPY services/matching-service/go.mod services/matching-service/go.sum ./

# Copy the proto directory from the root context into the correct location inside the container
COPY proto /app/proto

# Modify the replace path in go.mod to reflect the correct location of the proto directory
RUN sed -i 's|../../proto|/app/proto|' go.mod

RUN go mod download

# Copy all the service code for the matching service
COPY services/matching-service/. .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o matching-service ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Add non-root user
RUN adduser -D appuser

# Copy the binary from the builder
COPY --from=builder /app/matching-service .

# Switch to non-root user
USER appuser

EXPOSE 8080
CMD ["./matching-service"]

