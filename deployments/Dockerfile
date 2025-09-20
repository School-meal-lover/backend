# Builder stage
FROM golang:1.24.3-alpine AS builder

ENV GOOS=linux \
    CGO_ENABLED=0 \
    GOARCH=amd64 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY migrations/ ./migrations/

RUN go build -o server -ldflags="-s -w" ./server/main.go

# Runner stage
FROM alpine:3.18

RUN apk update && apk add --no-cache curl tar ca-certificates tzdata \
    && adduser -D appuser \
    && rm -rf /var/cache/apk/*

# migrate 도구 설치# Builder stage: Go 소스 코드를 컴파일하여 실행 파일을 만듭니다.
FROM golang:1.24.3-alpine AS builder

# 빌드에 필요한 환경 변수 설정
ENV GOOS=linux \
    CGO_ENABLED=0 \
    GOARCH=amd64

WORKDIR /build

# 1. 의존성 파일을 먼저 복사하여 Docker 레이어 캐시를 활용합니다.
COPY go.mod go.sum ./
RUN go mod download

# 2. 전체 소스 코드를 복사합니다.
COPY . .

# 3. [수정] main.go 파일의 정확한 경로를 지정하여 빌드합니다.
RUN go build -o server -ldflags="-s -w" ./cmd/server/main.go

# ---

# Runner stage: 빌드된 실행 파일과 필요 파일만으로 가벼운 최종 이미지를 만듭니다.
FROM alpine:3.18

# 1. 필요한 패키지 설치 및 비-루트 사용자 생성 (보안 강화)
RUN apk update && apk add --no-cache curl tar ca-certificates tzdata \
    && adduser -D appuser \
    && rm -rf /var/cache/apk/*

# 2. migrate 도구 설치
# v4.17.1 -> 최신 버전으로 변경할 수도 있습니다.
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate.linux-amd64 /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

WORKDIR /app

# 3. Builder 스테이지에서 컴파일된 실행 파일과 마이그레이션 파일만 복사합니다.
COPY --from=builder /build/server ./
COPY migrations ./migrations/

# 4. uploads 디렉토리 생성 및 권한 설정
RUN mkdir -p uploads && chown -R appuser:appuser /app

# 5. 비-루트 사용자로 전환
USER appuser

# 6. 서버 포트 노출
EXPOSE 8080

# 7. [수정] migrate 경로 수정 및 서버 실행
# 컨테이너 시작 시 마이그레이션을 먼저 실행하고, 성공하면 웹 서버를 시작합니다.
ENTRYPOINT ["/bin/sh", "-c", "migrate -path ./migrations -database \"${DATABASE_URL}\" up && ./server"]

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xz \
    && mkdir -p /app \
    && mv migrate /app/ 2>/dev/null || mv migrate.linux-amd64 /app/migrate \
    && chmod +x /app/migrate

WORKDIR /app

COPY --from=builder /build/server ./
COPY migrations ./migrations/

# S3 사용시 변경 필요
RUN mkdir -p uploads && chown -R appuser:appuser /app

USER appuser
  
EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "/app/migrate -path ./migrations -database \"${DATABASE_URL}\" up && ./server"]