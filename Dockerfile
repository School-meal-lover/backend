# Builder stage
FROM golang:1.24.3-alpine AS builder

ENV GOOS=linux \
    CGO_ENABLED=0 \
    GOARCH=amd64 

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY app/ ./app/
COPY docs/ ./docs/
COPY migrations/ ./migrations/

RUN go build -o server -ldflags="-s -w" ./app/cmd/main.go

# Runner stage
FROM alpine:3.18

RUN apk update && apk add --no-cache curl tar && curl --version && tar --version

# migrate 도구 설치
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz -o migrate.tar.gz \
    && tar xvzf migrate.tar.gz \
    && ls -l /  \
    && mkdir -p /app \
    && mv migrate /app/migrate 2>/dev/null || mv migrate.linux-amd64 /app/migrate 2>/dev/null || { echo "Failed to move migrate binary"; ls -l /; exit 1; } \
    && chmod +x /app/migrate \
    && ls -l /app \
    && /app/migrate --version || { echo "migrate version check failed"; exit 1; }

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /build/server ./
COPY --from=builder /build/migrations ./migrations/

# S3 사용시 변경 필요
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "/app/migrate -path /app/migrations -database \"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable\" up && /app/server"]