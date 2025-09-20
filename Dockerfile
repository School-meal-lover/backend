# Builder stage
FROM golang:1.24.3-alpine AS builder

ENV GOOS=linux CGO_ENABLED=0 GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server -ldflags="-s -w" ./cmd/server/main.go

# Runner stage
FROM alpine:3.18

RUN apk update && apk add --no-cache curl tar ca-certificates tzdata \
    && adduser -D appuser \
    && rm -rf /var/cache/apk/*

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

WORKDIR /app

COPY --from=builder /build/server ./
COPY migrations ./migrations/

RUN mkdir -p uploads && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "migrate -path ./migrations -database \"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable\" up && ./server"]

