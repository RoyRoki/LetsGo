# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /matching_server

# Copy go.mod and go.sum from the matching-service directory
COPY services/matching-service/go.mod services/matching-service/go.sum ./

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
COPY --from=builder /matching_server/matching-service .


# Switch to non-root user
USER appuser

EXPOSE 8080
CMD ["./matching-service"]

