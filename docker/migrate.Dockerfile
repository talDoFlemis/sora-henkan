FROM golang:1.25.3-alpine AS builder
WORKDIR /app

RUN apk add --no-cache \
  gcc \
  musl-dev \
  ca-certificates

COPY go.mod go.sum /app/

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target="/root/.cache/go-build" \
  CGO_ENABLED=0 GOOS=linux go build -o migrate  -ldflags '-s -w -extldflags "-static"' ./cmd/migrate

FROM ubuntu:oracular AS user
RUN useradd -u 10001 scratchuser

FROM scratch
WORKDIR /app


COPY --from=builder /app/migrate ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=user /etc/passwd /etc/passwd

USER scratchuser
STOPSIGNAL SIGINT
EXPOSE 8080

CMD ["/app/migrate", "-direction=up"]