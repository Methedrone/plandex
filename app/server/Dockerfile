# Multi-stage build for security and efficiency
# Stage 1: Builder with security hardening
FROM golang:1.23.10-alpine AS builder

# Install security updates and build dependencies using Alpine package manager
RUN apk update && \
    apk add --no-cache \
        git \
        gcc \
        g++ \
        make \
        musl-dev \
        ca-certificates \
        tzdata && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

# Create non-root user for build process
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy go.mod and go.sum files with proper ownership
COPY --chown=appuser:appgroup ./shared/go.mod ./shared/go.sum ./shared/
COPY --chown=appuser:appgroup ./server/go.mod ./server/go.sum ./server/

# Download dependencies with cache mount for efficiency
RUN --mount=type=cache,target=/go/pkg/mod \
    cd shared && go mod download && \
    cd ../server && go mod download

# Copy source code with proper ownership
COPY --chown=appuser:appgroup ./server ./server
COPY --chown=appuser:appgroup ./shared ./shared
COPY --chown=appuser:appgroup ./scripts /scripts

# Switch to non-root user for build
USER appuser

# Set working directory and build with optimizations
WORKDIR /app/server

# Build the application with security flags
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s -extldflags=-static" \
    -a -installsuffix cgo \
    -o plandex-server .

# Stage 2: Production image with minimal attack surface
FROM alpine:3.20 AS production

# Install only runtime dependencies and security updates
RUN apk update && \
    apk add --no-cache \
        ca-certificates \
        tzdata \
        dumb-init && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

# Create non-root user for runtime
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create necessary directories with proper permissions
RUN mkdir -p /app /scripts /plandex-server && \
    chown -R appuser:appgroup /app /scripts /plandex-server

# Copy binary and scripts from builder with proper ownership
COPY --from=builder --chown=appuser:appgroup /app/server/plandex-server /app/
COPY --from=builder --chown=appuser:appgroup /scripts /scripts/

# Make scripts executable
RUN chmod +x /scripts/*.sh /app/plandex-server

# Set working directory
WORKDIR /app

# Switch to non-root user
USER appuser

# Security: Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD /app/plandex-server --health-check || exit 1

# Set environment variables
ENV PORT=8099 \
    GOGC=50 \
    GOMEMLIMIT=512MiB \
    GOMAXPROCS=4

# Expose port
EXPOSE 8099

# Command to run the application
CMD ["/app/plandex-server"]
