# Build stage - compile Go application
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /build

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary with optimizations
# -a: force rebuild of packages
# -ldflags="-w -s": strip debug info for smaller binary
# -installsuffix cgo: for static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o ip-verifier-api \
    ./cmd/ip-verifier-api/main.go

# Runtime stage - minimal distroless image
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/ip-verifier-api .

# Copy initial GeoIP database (will be replaced by init container in K8s)
COPY --from=builder /build/data/GeoLite2-Country.mmdb ./data/

# Expose HTTP port
EXPOSE 8080

# Run as non-root user (distroless provides nonroot user)
USER nonroot:nonroot

# Start application
CMD ["/app/ip-verifier-api"]
