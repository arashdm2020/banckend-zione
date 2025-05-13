# Build stage
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Install required packages
RUN apk add --no-cache git gcc musl-dev

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zione-backend ./cmd/api

# Final stage
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder stage
COPY --from=builder /app/zione-backend .

# Copy config files
COPY --from=builder /app/configs ./configs

# Create necessary directories
RUN mkdir -p /app/uploads /app/logs

# Set environment variables
ENV GIN_MODE=release

# Expose port
EXPOSE 8080

# Run the application
CMD ["./zione-backend"] 