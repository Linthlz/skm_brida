# syntax=docker/dockerfile:1

# ---------- dev: source di-mount dari host, Air me-rebuild saat file berubah ----------
FROM golang:1.26-alpine AS dev
RUN apk add --no-cache git tzdata \
 && go install github.com/air-verse/air@v1.67.4
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 3000
CMD ["air", "-c", ".air.docker.toml"]

# ---------- build: binary statis untuk produksi ----------
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/singa-api . \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/singa-seed ./cmd/seed

# ---------- prod: image kecil, user non-root, folder upload sebagai volume ----------
FROM alpine:3.20 AS prod
RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -S app && adduser -S -G app app \
 && mkdir -p /app/public/storage && chown -R app:app /app
WORKDIR /app
COPY --from=build /out/singa-api /out/singa-seed ./
USER app
VOLUME ["/app/public/storage"]
EXPOSE 3000
CMD ["./singa-api"]
