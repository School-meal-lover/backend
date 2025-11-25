# 텍스트 업로드 API 사용 가이드

## 텍스트 형식

텍스트 파일은 다음과 같은 형식으로 작성해야 합니다:

```
RESTAURANT_1 또는 RESTAURANT_2
YYYY-MM-DD (주차 시작 날짜)
요일 YYYY-MM-DD
MealType
메뉴아이템1
메뉴아이템2
...
```

## 형식 설명

### 1. 첫 번째 줄: 레스토랑 타입

- `RESTAURANT_1`: 평일만 (월~금, 5일)
- `RESTAURANT_2`: 주말 포함 (월~일, 7일)

### 2. 두 번째 줄: 주차 시작 날짜

- 형식: `YYYY-MM-DD` (예: `2025-05-26`)

### 3. 날짜별 식사 정보

각 날짜는 다음과 같은 형식으로 작성:

```
요일 YYYY-MM-DD
MealType
메뉴아이템1
메뉴아이템2
...
```

#### 요일 형식

- `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`
- 날짜와 함께: `Monday 2025-05-26`

#### MealType

- `Breakfast`: 아침
- `Lunch_1`: 일품 메뉴 (점심)
- `Lunch_2`: 일반 메뉴 (점심)
- `Dinner`: 저녁

#### 메뉴 아이템

- 각 MealType 아래에 메뉴 아이템을 한 줄씩 작성
- 빈 줄은 무시됩니다

## 예제

### Restaurant1 예제 (평일만)

`testdata/example_restaurant1.txt` 파일 참조

### Restaurant2 예제 (주말 포함)

`testdata/example_restaurant2.txt` 파일 참조

## API 사용 예제

### Bearer Token 인증

- **토큰 설정**: 환경변수 `BEARER_TOKEN`에 토큰 값을 설정하세요
- **기본값**: 환경변수가 없으면 `gistsikdang` 사용 (개발 환경용)
- **프로덕션**: 반드시 강력한 토큰으로 변경하세요
- **사용법**: Authorization 헤더에 `Bearer <토큰값>` 형식으로 전송

### curl 명령어

```bash
# 토큰을 환경변수로 설정 (선택사항)
export BEARER_TOKEN="your-secret-token"

# Restaurant1 업로드
curl -X POST https://grrrr.me/api/v1/upload/text \
  -H "Authorization: Bearer ${BEARER_TOKEN:-gistsikdang}" \
  -H "Content-Type: text/plain" \
  --data-binary @testdata/example_restaurant1.txt

# Restaurant2 업로드
curl -X POST https://grrrr.me/api/v1/upload/text \
  -H "Authorization: Bearer ${BEARER_TOKEN:-gistsikdang}" \
  -H "Content-Type: text/plain" \
  --data-binary @testdata/example_restaurant2.txt
```

### Python 예제

```python
import os
import requests

# 환경변수에서 토큰 가져오기 (없으면 기본값 사용)
token = os.getenv("BEARER_TOKEN", "gistsikdang")

url = "https://grrrr.me/api/v1/upload/text"
headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "text/plain"
}

with open("testdata/example_restaurant1.txt", "r", encoding="utf-8") as f:
    text_data = f.read()

response = requests.post(url, headers=headers, data=text_data)
print(response.json())
```

### JavaScript 예제

```javascript
const fs = require("fs");

// 환경변수에서 토큰 가져오기 (없으면 기본값 사용)
const token = process.env.BEARER_TOKEN || "gistsikdang";
const textData = fs.readFileSync("testdata/example_restaurant1.txt", "utf8");

fetch("https://grrrr.me/api/v1/upload/text", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "text/plain",
  },
  body: textData,
})
  .then((response) => response.json())
  .then((data) => console.log(data));
```

## 주의사항

1. **인코딩**: 텍스트 파일은 UTF-8 인코딩을 사용해야 합니다.
2. **날짜 형식**: 날짜는 반드시 `YYYY-MM-DD` 형식을 사용해야 합니다.
3. **MealType 대소문자**: MealType은 대소문자를 구분하지 않지만, 표준 형식은 `Breakfast`, `Lunch_1`, `Lunch_2`, `Dinner`입니다.
4. **빈 줄**: 빈 줄은 무시되므로 가독성을 위해 사용할 수 있습니다.
5. **Restaurant1 vs Restaurant2**:
   - Restaurant1은 평일(월~금)만 처리됩니다.
   - Restaurant2는 주말(토~일)까지 포함하여 처리됩니다.

## 응답 형식

성공 시:

```json
{
  "success": true,
  "restaurant_type": "RESTAURANT_1",
  "week_id": "uuid-string",
  "week_start_date": "2025-05-26",
  "total_meals": 20,
  "total_menu_items": 100,
  "message": "Text processed successfully"
}
```

실패 시:

```json
{
  "success": false,
  "error": "에러 메시지"
}
```
