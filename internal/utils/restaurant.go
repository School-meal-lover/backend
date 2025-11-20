package utils

import (
	"fmt"
	"strings"

	"github.com/School-meal-lover/backend/internal/models"
)

// ParseRestaurantType 문자열을 RestaurantType으로 변환
// 데이터베이스 ENUM 값과 일치하는지 검증
func ParseRestaurantType(s string) (models.RestaurantType, error) {
	normalized := strings.ToUpper(strings.TrimSpace(s))
	
	// ENUM 값과 직접 비교
	if normalized == string(models.Restaurant1) || normalized == "RESTAURANT_1" {
		return models.Restaurant1, nil
	}
	if normalized == string(models.Restaurant2) || normalized == "RESTAURANT_2" {
		return models.Restaurant2, nil
	}
	
	// 숫자로도 파싱 시도 (식당 1, 식당 2 등)
	if strings.Contains(normalized, "1") {
		return models.Restaurant1, nil
	}
	if strings.Contains(normalized, "2") {
		return models.Restaurant2, nil
	}
	
	return "", fmt.Errorf("invalid restaurant type: %s (must be RESTAURANT_1 or RESTAURANT_2)", s)
}

// ValidateRestaurantType RestaurantType이 유효한지 검증
func ValidateRestaurantType(rt models.RestaurantType) error {
	if rt != models.Restaurant1 && rt != models.Restaurant2 {
		return fmt.Errorf("invalid restaurant type: %s", rt)
	}
	return nil
}

// GetDaysForRestaurant 식당 타입에 따른 요일 수 반환
func GetDaysForRestaurant(rt models.RestaurantType) (int, error) {
	if err := ValidateRestaurantType(rt); err != nil {
		return 0, err
	}
	
	if rt == models.Restaurant1 {
		return 5, nil // 월화수목금
	}
	return 7, nil // 월화수목금토일
}

