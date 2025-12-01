package services

import (
	"fmt"
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
func (s *TextService) ProcessText(restaurantName string, requestBody models.TextUploadRequest) (*models.ExcelProcessResult, error) {
	restaurantType := models.RestaurantType(restaurantName)
	weekStartDate, err := time.Parse("2006-01-02", requestBody.WeekStartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD)", requestBody.WeekStartDate)
	}

	// 3. Week 생성
	weekID, err := s.mealRepo.FindOrCreateWeek(weekStartDate, restaurantType)
	if err != nil {
		return nil, fmt.Errorf("failed to create week: %w", err)
	}

	var totalMeals, totalMenuItems int
	for _, dailyMeal := range requestBody.Meals {
		date, err := time.Parse("2006-01-02", dailyMeal.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format in meals: %s (expected YYYY-MM-DD)", dailyMeal.Date)
		}

		mealID, err := s.mealRepo.InsertMeal(&models.Meal{WeekID: weekID, Date: date, DayOfWeek: dailyMeal.DayOfWeek, MealType: dailyMeal.MealType})
		if err != nil {
			return nil, fmt.Errorf("failed to insert meal: %w", err)
		}

		menuItems := []models.MenuItem{}
		for _, menuItemName := range dailyMeal.MenuItems {
			menuItems = append(menuItems, models.MenuItem{MealID: mealID, Name: menuItemName.Name, Category: menuItemName.Category})
			totalMenuItems++
		}
		if err := s.mealRepo.InsertMenuItems(menuItems); err != nil {
			return nil, fmt.Errorf("failed to insert menu items: %w", err)
		}
		totalMeals++
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
