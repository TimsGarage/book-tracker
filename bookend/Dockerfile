# Stage 1: Build binary
FROM golang:alpine AS builder

# Install C compiler and musl development headers required for SQLite CGO driver
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Cache Go dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source files
COPY . .

# Compile optimized binary with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o bookend main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

# Install root CA certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create dedicated non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bookend /app/bookend

# Prepare persistent data directory with non-root ownership
RUN mkdir -p /data && chown -R appuser:appgroup /data /app

# Set production defaults for environment variables
ENV PORT=8080 \
    GIN_MODE=release \
    DB_PATH=/data/bookend.db

# Expose application port and SQLite storage volume
EXPOSE 8080
VOLUME ["/data"]

# Run container as unprivileged user
USER appuser

ENTRYPOINT ["/app/bookend"]

