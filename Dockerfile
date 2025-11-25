# Dockerfile for srn-coreutils

# ---- Build Stage ----
# This stage compiles the Go binary.
FROM golang:1.22-alpine AS build

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Get version from git
ARG VERSION=$(git describe --tags --always --dirty --first-parent 2>/dev/null || echo "dev")

# Build the binary with version information
RUN go build -ldflags="-s -w -X main.version=${VERSION}" -o /srn .

# ---- Final Stage ----
# This stage creates the final, minimal image.
FROM alpine:latest

# Install runtime dependencies for srn commands
# yt-dlp requires python and pip
RUN apk add --no-cache \
    bash \
    git \
    podman \
    ffmpeg \
    p7zip \
    python3 \
    py3-pip

# Install yt-dlp using pip
RUN pip3 install --no-cache-dir yt-dlp

# Copy the compiled binary from the build stage
COPY --from=build /srn /usr/local/bin/srn

# Set the entrypoint to the srn binary
ENTRYPOINT ["srn"]

# Default command to show help
CMD ["--help"]
