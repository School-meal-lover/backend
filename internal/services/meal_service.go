package services

import (
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

// 특정 레스토랑의 주간 식단을 조회
func (s *MealService) GetRestaurantWeekMeals(restaurantNameParam string, date string) (*models.RestaurantMealsResponse, error) {
	// 날짜 형식 검증
	var restaurantType models.RestaurantType
	switch restaurantNameParam {
	case string(models.Restaurant1):
		restaurantType = models.Restaurant1
	case string(models.Restaurant2):
		restaurantType = models.Restaurant2
	default:
		return &models.RestaurantMealsResponse{Success: false, Error: "Invalid restaurant name", Code: "INVALID_RESTAURANR_NAME"}, nil
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return &models.RestaurantMealsResponse{
			Success: false,
			Error:   "Invalid date format. Use YYYY-MM-DD",
			Code:    "INVALID_DATE_FORMAT",
		}, nil
	}
	// 주차 정보 조회 및 에러 처리
	week, err := s.mealRepo.GetWeekInfo(restaurantType, date)
	if err != nil {
		return &models.RestaurantMealsResponse{
			Success: false,
			Error:   "WEEK_DATA_NOT_FOUND",
			Code:    "No meal data found for the specified week",
		}, nil
	}

	// 식단 데이터 조회 및 에러 처리
	orderedmealsByDay, summary, err := s.mealRepo.GetMealsData(week.ID)
	if err != nil {
		// 식단 데이터 조회 실패는 특정 코드로 처리
		return &models.RestaurantMealsResponse{
			Success: false,
			Error:   "Failed to retrieve meal data",
			Code:    "MEAL_DATA_RETRIEVAL_FAILED",
		}, nil
	}

	// 성공 응답 구성
	response := &models.RestaurantMealsData{
		Restaurant: string(restaurantType),
		Week:       week,
		MealsByDay: orderedmealsByDay,
		Summary:    summary,
	}

	return &models.RestaurantMealsResponse{
		Success: true,
		Data:    response,
	}, nil
}
