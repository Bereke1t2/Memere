# syntax=docker/dockerfile:1.7

# ---------------------------------------------------------------------------
# Build stage — compiles the API server and the migration runner.
# A single build stage keeps the layer cache efficient: go mod download is
# cached separately from the source copy so dependency changes don't bust the
# source layer and vice-versa.
# ---------------------------------------------------------------------------
FROM golang:1.26-bookworm AS build

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the API server. -trimpath removes local path prefixes from stack
# traces; -s -w strips debug info and DWARF for a smaller binary.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w \
      -X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.Version=${VERSION}' \
      -X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.Commit=${COMMIT}' \
      -X 'github.com/Bereke1t2/Memere/memere-backend/internal/version.BuildTime=${BUILD_TIME}'" \
    -o /out/api ./cmd/api

# Build the migration runner (no ldflags needed — it logs nothing version-specific).
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath -ldflags="-s -w" \
    -o /out/migrate ./cmd/migrate

# ---------------------------------------------------------------------------
# Runtime stage — distroless/static-debian12:nonroot
# No shell, no package manager, no setuid binaries. The nonroot variant
# drops root before the process starts, satisfying Non-Negotiable #6 and CIS
# Benchmark for containers. Digest is pinned so CI detects upstream changes.
#
# Pin update procedure: `docker pull gcr.io/distroless/static-debian12:nonroot`
# then replace the digest below with `docker inspect --format '{{index
# .RepoDigests 0}}' gcr.io/distroless/static-debian12:nonroot`.
# ---------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot@sha256:3d0f463de06b7ddff27684ec3bfd0b54a425149d0f8685308b1d9d7b6579cf4c

# Copy the compiled binaries and the embedded migrations directory.
COPY --from=build /out/api     /api
COPY --from=build /out/migrate /migrate
COPY --from=build /src/migrations /migrations

# The nonroot base image already sets USER 65532:65532 (nonroot).
# Declaring it explicitly makes intent clear to static scanners.
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/api"]
