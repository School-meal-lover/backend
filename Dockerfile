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
RUN go build -o migrate -ldflags="-s -w" ./app/cmd/migrate/main.go

# Runner stage
FROM alpine:3.18

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /build/server ./
COPY --from=builder /build/migrate ./migrate
COPY migrations ./migrations

# S3 사용시 변경 필요
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["./server"]

EXPOSE 8080
