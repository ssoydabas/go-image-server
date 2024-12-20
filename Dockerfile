# Build stage
FROM golang:1.22-alpine AS builder

# Install required dependencies for webp
RUN apk add --no-cache gcc musl-dev libwebp-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install libwebp for runtime
RUN apk add --no-cache libwebp

# Create app user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder FIRST
COPY --from=builder /app/main .

# Create data directory and set permissions
RUN mkdir -p /app/data /app/dev-data && \
    chown -R appuser:appgroup /app && \
    chmod 755 /app/main && \
    chmod 755 /app/data && \
    chmod 755 /app/dev-data

# Switch to app user
USER appuser

EXPOSE 8080

CMD ["./main"]