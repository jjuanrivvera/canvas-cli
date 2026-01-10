# Canvas CLI - Dockerfile
# Multi-stage build for minimal final image

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w -X main.version=docker" \
    -o canvas ./cmd/canvas

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 canvas && \
    adduser -D -u 1000 -G canvas canvas

# Set working directory
WORKDIR /home/canvas

# Copy binary from builder
COPY --from=builder /build/canvas /usr/local/bin/canvas

# Set ownership
RUN chown -R canvas:canvas /home/canvas

# Switch to non-root user
USER canvas

# Set entrypoint
ENTRYPOINT ["canvas"]

# Default command
CMD ["--help"]
