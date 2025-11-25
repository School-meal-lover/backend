package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/School-meal-lover/backend/internal/models"
	"github.com/School-meal-lover/backend/internal/repository"
)

type TextService struct {
	mealRepo *repository.MealRepository
}

func NewTextService(mealRepo *repository.MealRepository) *TextService {
	return &TextService{
		mealRepo: mealRepo,
	}
}

// ProcessText는 텍스트 형식의 식단 데이터를 처리합니다
// 텍스트 형식:
// RESTAURANT_1 또는 RESTAURANT_2
// 2025-05-26 (주차 시작 날짜)
// Monday 2025-05-26
// Breakfast
// 밥
// 국
// 반찬
// Lunch_1
// 메인메뉴
// ...
func (s *TextService) ProcessText(text string) (*models.ExcelProcessResult, error) {
	lines := strings.Split(text, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid text format: too few lines")
	}

	// 1. Restaurant 타입 파싱
	restaurantLine := strings.TrimSpace(lines[0])
	var restaurantType models.RestaurantType
	switch strings.ToUpper(restaurantLine) {
	case "RESTAURANT_1", "RESTAURANT_1\n", "RESTAURANT_1\r":
		restaurantType = models.Restaurant1
	case "RESTAURANT_2", "RESTAURANT_2\n", "RESTAURANT_2\r":
		restaurantType = models.Restaurant2
	default:
		return nil, fmt.Errorf("invalid restaurant type: %s (expected RESTAURANT_1 or RESTAURANT_2)", restaurantLine)
	}

	// 2. 주차 시작 날짜 파싱
	dateLine := strings.TrimSpace(lines[1])
	weekStartDate, err := time.Parse("2006-01-02", dateLine)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", dateLine)
	}

	// 3. Week 생성
	weekID, err := s.mealRepo.FindOrCreateWeek(weekStartDate, restaurantType)
	if err != nil {
		return nil, fmt.Errorf("failed to create week: %w", err)
	}

	// 4. 날짜별 식사 데이터 파싱
	totalMeals := 0
	totalMenuItems := 0
	currentDate := ""
	currentDayOfWeek := ""
	currentMealType := ""
	var currentMenuItems []string

	processedDates := make(map[string]bool)

	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 날짜 라인 확인 (예: "Monday 2025-05-26")
		if strings.Contains(line, "2025-") || strings.Contains(line, "2024-") || strings.Contains(line, "2026-") {
			// 이전 식사 데이터 저장
			if currentDate != "" && currentMealType != "" {
				_, err := s.saveMeal(weekID, currentDate, currentDayOfWeek, currentMealType, currentMenuItems)
				if err != nil {
					log.Printf("Failed to save meal: %v", err)
				} else {
					totalMeals++
					totalMenuItems += len(currentMenuItems)
				}
			}

			// 새 날짜 파싱
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentDayOfWeek = parts[0]
				dateStr := parts[1]
				_, err := time.Parse("2006-01-02", dateStr)
				if err == nil {
					currentDate = dateStr
					if !processedDates[currentDate] {
						processedDates[currentDate] = true
					}
					currentMealType = ""
					currentMenuItems = []string{}
				}
			}
			continue
		}

		// MealType 확인
		upperLine := strings.ToUpper(line)
		if upperLine == "BREAKFAST" || upperLine == "LUNCH_1" || upperLine == "LUNCH_2" || upperLine == "DINNER" {
			// 이전 식사 데이터 저장
			if currentDate != "" && currentMealType != "" {
				_, err := s.saveMeal(weekID, currentDate, currentDayOfWeek, currentMealType, currentMenuItems)
				if err != nil {
					log.Printf("Failed to save meal: %v", err)
				} else {
					totalMeals++
					totalMenuItems += len(currentMenuItems)
				}
			}

			// 새 MealType 설정
			if upperLine == "LUNCH_1" {
				currentMealType = "Lunch_1"
			} else if upperLine == "LUNCH_2" {
				currentMealType = "Lunch_2"
			} else if upperLine == "BREAKFAST" {
				currentMealType = "Breakfast"
			} else if upperLine == "DINNER" {
				currentMealType = "Dinner"
			} else {
				currentMealType = upperLine
			}
			currentMenuItems = []string{}
			continue
		}

		// 메뉴 아이템 추가
		if currentDate != "" && currentMealType != "" {
			currentMenuItems = append(currentMenuItems, line)
		}
	}

	// 마지막 식사 데이터 저장
	if currentDate != "" && currentMealType != "" {
		_, err := s.saveMeal(weekID, currentDate, currentDayOfWeek, currentMealType, currentMenuItems)
		if err != nil {
			log.Printf("Failed to save meal: %v", err)
		} else {
			totalMeals++
			totalMenuItems += len(currentMenuItems)
		}
	}

	return &models.ExcelProcessResult{
		Success:        true,
		RestaurantType: string(restaurantType),
		WeekStartDate:  weekStartDate.Format("2006-01-02"),
		WeekID:         weekID,
		TotalMeals:     totalMeals,
		TotalMenuItems: totalMenuItems,
		Message:        "Text processed successfully",
	}, nil
}

func (s *TextService) saveMeal(weekID, date, dayOfWeek, mealType string, menuItems []string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid date: %w", err)
	}

	meal := &models.Meal{
		WeekID:    weekID,
		Date:      parsedDate,
		DayOfWeek: dayOfWeek,
		MealType:  mealType,
	}

	mealID, err := s.mealRepo.FindOrCreateMeal(meal)
	if err != nil {
		return "", err
	}

	// 메뉴 아이템 저장
	if len(menuItems) > 0 {
		menuItemModels := s.buildMenuItems(mealID, menuItems, mealType)
		if err := s.mealRepo.InsertMenuItems(menuItemModels); err != nil {
			log.Printf("Failed to insert menu items: %v", err)
		}
	}

	return mealID, nil
}

func (s *TextService) buildMenuItems(mealID string, itemNames []string, mealType string) []models.MenuItem {
	categories := s.getCategoriesForMealType(mealType)
	var menuItems []models.MenuItem

	for idx, name := range itemNames {
		var category string
		if idx < len(categories) {
			category = categories[idx]
		} else {
			category = "기타"
		}
		menuItems = append(menuItems, models.MenuItem{
			MealID:   mealID,
			Category: category,
			Name:     strings.TrimSpace(name),
			NameEn:   "",
			Price:    0.00,
		})
	}
	return menuItems
}

func (s *TextService) getCategoriesForMealType(mealType string) []string {
	if mealType == "Lunch_1" {
		return []string{"메인메뉴"}
	}
	if mealType == "Breakfast" {
		return []string{"밥", "국", "반찬", "메인메뉴", "반찬", "반찬", "반찬", "반찬", "반찬"}
	}
	return []string{"밥", "국", "메인메뉴", "메인메뉴", "반찬", "반찬"}
}

