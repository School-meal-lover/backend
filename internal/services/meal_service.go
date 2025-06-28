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
func (s *MealService) GetRestaurantWeekMeals(restaurantID, date string) (*models.RestaurantMealsResponse, error) {
	// 날짜 형식 검증
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return &models.RestaurantMealsResponse{
			Success: false,
			Error:   "Invalid date format. Use YYYY-MM-DD",
			Code:    "INVALID_DATE_FORMAT",
		}, nil
	}

	// 레스토랑 정보 조회 및 에러 처리
	restaurant, err := s.mealRepo.GetRestaurantInfo(restaurantID)
	if err != nil {
		return s.mealRepo.HandleRepositoryError(err, "RESTAURANT_NOT_FOUND", "Restaurant not found")
	}

	// 주차 정보 조회 및 에러 처리
	week, err := s.mealRepo.GetWeekInfo(restaurantID, date)
	if err != nil {
		return s.mealRepo.HandleRepositoryError(err, "WEEK_DATA_NOT_FOUND", "No meal data found for the specified week")
	}

	// 식단 데이터 조회 및 에러 처리
	mealsByDay, summary, err := s.mealRepo.GetMealsData(week.ID)
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
		Restaurant: restaurant,
		Week:       week,
		MealsByDay: mealsByDay,
		Summary:    summary,
	}

	return &models.RestaurantMealsResponse{
		Success: true,
		Data:    response,
	}, nil
}