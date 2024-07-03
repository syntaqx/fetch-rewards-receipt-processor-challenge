FROM golang:1.22-alpine AS builder

RUN apk add --update --no-cache \
  ca-certificates tzdata openssh git mercurial && update-ca-certificates \
  && rm -rf /var/cache/apk/*

WORKDIR /src

COPY go.mod* go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install ./cmd/...

# RUN --mount=type=cache,target=/go/pkg/mod go generate ./...

FROM alpine

RUN adduser -S -D -H -h /app appuser

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=appuser:appuser /go/bin/* /bin/

COPY --chown=appuser:appuser ./static /app/static

WORKDIR /app

USER appuser

ENV PORT=8080
EXPOSE $PORT

CMD ["app"]
