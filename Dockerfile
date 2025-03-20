# Step 1: Use the official Golang image as the base
FROM golang:1.24-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY .env.development /app/.env.development

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o /app/letsgo ./cmd/api/main.go

# Step 2: Create a smaller image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary from the previous stage
COPY --from=builder /app/letsgo .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the application
CMD ["./letsgo"]
