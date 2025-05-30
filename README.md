# Backend

grrrrr의 벡엔드 어플리케이션입니다.

## Environment

이 어플리케이션을 작동시키기 위해서 필수적으로 필요한 환경변수는 다음과 같습니다.

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydatabase
DB_USER=myuser
DB_PASSWORD=mypassword
```

## how to upload excel file

- 로컬 파일 처리

```go
curl -X POST http://localhost:8080/api/v1/process/excel/local \
  -H "Content-Type: application/json" \
  -d '{"file_path": "assets/2025_5_5_ko.xlsx"}'
```

- 파일 업로드

```go
curl -X POST http://localhost:8080/api/v1/upload/excel \
  -F "excel=@asssets/2025_5_5_ko.xlsx"
```

## how to build swagger file

```bash
swag init -d ./cmd
```