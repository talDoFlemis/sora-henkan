FROM golang:1.25.3-bookworm AS builder

# Install libvips development headers and build tools
# libvips-dev pulls in all dependencies (libjpeg, libpng, libwebp, etc.)
RUN apt-get update && apt-get install -y \
  build-essential \
  pkg-config \
  libvips-dev \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum /app/

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

COPY . .

ARG package_name
ENV PACKAGE_NAME=${package_name}

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target="/root/.cache/go-build" \
  CGO_ENABLED=1 go build -o /main -tags vips ./cmd/${PACKAGE_NAME}

RUN mkdir /deps

# Use ldd to find all shared library dependencies of our compiled binary
# Then, copy all those libraries (and their parent directories) into /deps
#
# - ldd /main: Lists all .so files the /main binary needs.
# - grep '=> /': Filters for dynamically linked libraries on the filesystem.
# - awk '{print $3}': Extracts just the file path (e.g., /usr/lib/x86_64-linux-gnu/libvips.so.42).
# - xargs -I '{}' cp -L --parents '{}' /deps:
#   - cp -L: Follows symlinks to copy the actual library file.
#   - --parents: Copies the file *with its full directory structure*
#     (e.g., creates /deps/usr/lib/x86_64-linux-gnu/libvips.so.42)
RUN ldd /main | grep '=> /' | awk '{print $3}' | xargs -I '{}' cp -L --parents '{}' /deps


# We use -debian12 to match our builder stage.
FROM gcr.io/distroless/cc-debian12:debug-nonroot

WORKDIR /app

COPY --from=builder /main .

COPY --from=builder /deps/ /

ENTRYPOINT ["/app/main"]
