package services

import (
	"strings"
	"testing"
	"time"

	"github.com/School-meal-lover/backend/internal/models"
)

func TestParseRestaurantType(t *testing.T) {
	service := &TextService{}

	tests := []struct {
		name     string
		text     string
		expected models.RestaurantType
		wantErr  bool
	}{
		{
			name:     "RESTAURANT_1",
			text:     "식당: RESTAURANT_1\n주 시작일: 2025-06-28",
			expected: models.Restaurant1,
			wantErr:  false,
		},
		{
			name:     "RESTAURANT_2",
			text:     "식당: RESTAURANT_2\n주 시작일: 2025-06-28",
			expected: models.Restaurant2,
			wantErr:  false,
		},
		{
			name:     "식당 1",
			text:     "식당: 식당 1\n주 시작일: 2025-06-28",
			expected: models.Restaurant1,
			wantErr:  false,
		},
		{
			name:     "식당 2",
			text:     "식당: 식당 2\n주 시작일: 2025-06-28",
			expected: models.Restaurant2,
			wantErr:  false,
		},
		{
			name:     "식당 정보 없음",
			text:     "주 시작일: 2025-06-28",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseRestaurantType(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRestaurantType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("parseRestaurantType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseWeekStartDate(t *testing.T) {
	service := &TextService{}

	tests := []struct {
		name     string
		text     string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "정상 날짜",
			text:     "식당: RESTAURANT_1\n주 시작일: 2025-06-28",
			expected: time.Date(2025, 6, 28, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "날짜 정보 없음",
			text:     "식당: RESTAURANT_1",
			expected: time.Time{},
			wantErr:  true,
		},
		{
			name:     "잘못된 날짜 형식",
			text:     "식당: RESTAURANT_1\n주 시작일: 2025/06/28",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseWeekStartDate(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWeekStartDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !result.Equal(tt.expected) {
				t.Errorf("parseWeekStartDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBuildDatesFromText(t *testing.T) {
	service := &TextService{}

	weekStartDate := time.Date(2025, 6, 28, 0, 0, 0, 0, time.UTC)

	t.Run("식당1_5일", func(t *testing.T) {
		dates, err := service.buildDatesFromText(weekStartDate, models.Restaurant1)

		if err != nil {
			t.Fatalf("buildDatesFromText() error = %v", err)
		}

		if len(dates) != 5 {
			t.Fatalf("buildDatesFromText() returned %d dates, want 5", len(dates))
		}

		expectedDays := []string{"월요일", "화요일", "수요일", "목요일", "금요일"}
		for i, date := range dates {
			if date.DayOfWeek != expectedDays[i] {
				t.Errorf("buildDatesFromText() date[%d].DayOfWeek = %v, want %v", i, date.DayOfWeek, expectedDays[i])
			}
			expectedDate := weekStartDate.AddDate(0, 0, i)
			if date.Date != expectedDate.Format("2006-01-02") {
				t.Errorf("buildDatesFromText() date[%d].Date = %v, want %v", i, date.Date, expectedDate.Format("2006-01-02"))
			}
		}
	})

	t.Run("식당2_7일", func(t *testing.T) {
		dates, err := service.buildDatesFromText(weekStartDate, models.Restaurant2)

		if err != nil {
			t.Fatalf("buildDatesFromText() error = %v", err)
		}

		if len(dates) != 7 {
			t.Fatalf("buildDatesFromText() returned %d dates, want 7", len(dates))
		}

		expectedDays := []string{"월요일", "화요일", "수요일", "목요일", "금요일", "토요일", "일요일"}
		for i, date := range dates {
			if date.DayOfWeek != expectedDays[i] {
				t.Errorf("buildDatesFromText() date[%d].DayOfWeek = %v, want %v", i, date.DayOfWeek, expectedDays[i])
			}
			expectedDate := weekStartDate.AddDate(0, 0, i)
			if date.Date != expectedDate.Format("2006-01-02") {
				t.Errorf("buildDatesFromText() date[%d].Date = %v, want %v", i, date.Date, expectedDate.Format("2006-01-02"))
			}
		}
	})
}

func TestSplitTextByDays(t *testing.T) {
	service := &TextService{}

	text := `월요일 (2025-06-28)
아침:
밥: 쌀밥

화요일 (2025-06-29)
아침:
밥: 현미밥

수요일 (2025-06-30)
아침:
밥: 쌀밥

토요일 (2025-07-03)
아침:
밥: 현미밥

일요일 (2025-07-04)
아침:
밥: 쌀밥`

	daySections := service.splitTextByDays(text)

	if len(daySections) != 5 {
		t.Fatalf("splitTextByDays() returned %d sections, want 5", len(daySections))
	}

	if !strings.Contains(daySections["월요일"], "쌀밥") {
		t.Error("splitTextByDays() 월요일 section does not contain expected content")
	}

	if !strings.Contains(daySections["화요일"], "현미밥") {
		t.Error("splitTextByDays() 화요일 section does not contain expected content")
	}

	if !strings.Contains(daySections["토요일"], "현미밥") {
		t.Error("splitTextByDays() 토요일 section does not contain expected content")
	}

	if !strings.Contains(daySections["일요일"], "쌀밥") {
		t.Error("splitTextByDays() 일요일 section does not contain expected content")
	}
}

func TestParseMenuItemsFromDayText(t *testing.T) {
	service := &TextService{}

	dayText := `아침:
밥: 쌀밥
국: 된장국
반찬: 김치
메인메뉴: 계란후라이

점심1:
메인메뉴: 제육볶음

점심2:
밥: 쌀밥
국: 미역국`

	items, err := service.parseMenuItemsFromDayText(dayText, "Breakfast")
	if err != nil {
		t.Fatalf("parseMenuItemsFromDayText() error = %v", err)
	}

	expectedItems := []struct {
		Category string
		Name     string
	}{
		{"밥", "쌀밥"},
		{"국", "된장국"},
		{"반찬", "김치"},
		{"메인메뉴", "계란후라이"},
	}

	if len(items) != len(expectedItems) {
		t.Fatalf("parseMenuItemsFromDayText() returned %d items, want %d", len(items), len(expectedItems))
	}

	for i, item := range items {
		if item.Category != expectedItems[i].Category {
			t.Errorf("parseMenuItemsFromDayText() items[%d].Category = %v, want %v", i, item.Category, expectedItems[i].Category)
		}
		if item.Name != expectedItems[i].Name {
			t.Errorf("parseMenuItemsFromDayText() items[%d].Name = %v, want %v", i, item.Name, expectedItems[i].Name)
		}
	}
}

