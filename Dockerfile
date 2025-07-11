FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY cmd/ ./cmd/
COPY client/ ./client/
COPY models/ ./models/
COPY pkg/ ./pkg/

ARG TARGETOS=linux
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-w -s" -trimpath -o bin/yaca ./cmd/yaca/main.go

# --- Runtime stage ---
FROM alpine:3.22 AS runtime

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser -D -u 10001 appuser

WORKDIR /app

# Copy the binary and entrypoint from builder
COPY --from=builder /app/bin/yaca .
COPY entrypoint.sh .

# Ensure entrypoint is executable
RUN chmod +x ./entrypoint.sh

# Switch to non-root user
USER appuser

ENV ENVIRONMENT=production

ENTRYPOINT ["/app/entrypoint.sh"]
