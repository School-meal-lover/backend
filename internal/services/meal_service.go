package services

import (
	"strings"
	"time"

	"github.com/School-meal-lover/backend/internal/models"
	"github.com/School-meal-lover/backend/internal/repository"
)

type MealService struct {
	mealRepo *repository.MealRepository
}

func NewMealService(mealRepo *repository.MealRepository) *MealService {
    return &MealService{
        mealRepo: mealRepo,
    }
}
// 간소화된 레스토랑 식단 조회
func (s *MealService) GetRestaurantWeekMeals(restaurantID, date string) (*models.RestaurantMealsResponse, error) {
    // 날짜 형식 검증
    if _, err := time.Parse("2006-01-02", date); err != nil {
        return &models.RestaurantMealsResponse{
            Success: false,
            Error:   "Invalid date format. Use YYYY-MM-DD",
            Code:    "INVALID_DATE_FORMAT",
        }, nil
    }

    // 데이터 조회
    data, err := s.mealRepo.GetRestaurantWeekMeals(restaurantID, date)
    if err != nil {
        var code string
        var message string

        switch {
        case strings.Contains(err.Error(), "restaurant not found"):
            code = "RESTAURANT_NOT_FOUND"
            message = "Restaurant not found"
        case strings.Contains(err.Error(), "week not found"):
            code = "WEEK_DATA_NOT_FOUND"
            message = "No meal data found for the specified week"
        default:
            code = "INTERNAL_ERROR"
            message = "Internal server error"
        }

        return &models.RestaurantMealsResponse{
            Success: false,
            Error:   message,
            Code:    code,
        }, nil
    }

    return &models.RestaurantMealsResponse{
        Success: true,
        Data:    data,
    }, nil
}