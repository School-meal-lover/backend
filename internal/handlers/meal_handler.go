package handlers

import (
	"net/http"

	"github.com/School-meal-lover/backend/internal/models"
	"github.com/School-meal-lover/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type MealHandler struct {
	mealService *services.MealService
}

func NewMealHandler(mealService *services.MealService) *MealHandler {
	return &MealHandler{mealService: mealService}
}

// @Summary Get meals for a specific restaurant on a specific date
// @Description Retrieve meals for a restaurant on a given date
// @Tags Meals
// @Accept json
// @Produce json
// @Param id path string true "레스토랑 ID (UUID 형식)" format:"uuid" example:"db6017d7-094a-4252-bdb6-28cc3f832447"
// @Param date query string true "조회할 날짜 (YYYY-MM-DD 형식)" example:"2025-06-28"
// @Success 200 {object} models.RestaurantMealsResponse "성공적으로 식단 정보 조회"
// @Failure 400 "잘못된 요청 파라미터 (ID 또는 날짜 형식 오류)"
// @Failure 404 "레스토랑 또는 해당 날짜의 식단 정보를 찾을 수 없음"
// @Failure 500 "서버 내부 오류 발생"
// @Router /restaurants/{id} [get]
func (h *MealHandler) GetRestaurantMeals(c *gin.Context) {
	restaurantID := c.Param("id")
	date := c.Query("date")

	// 파라미터 검증
	if restaurantID == "" {
		c.JSON(http.StatusBadRequest, models.RestaurantMealsResponse{
			Success: false,
			Error:   "restaurant_id is required",
			Code:    "MISSING_RESTAURANT_ID",
		})
		return
	}

	if date == "" {
		c.JSON(http.StatusBadRequest, models.RestaurantMealsResponse{
			Success: false,
			Error:   "date parameter is required (YYYY-MM-DD)",
			Code:    "MISSING_DATE",
		})
		return
	}

	// 서비스 호출
	response, err := h.mealService.GetRestaurantWeekMeals(restaurantID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RestaurantMealsResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    "INTERNAL_ERROR",
		})
		return
	}

	// HTTP 상태 코드 결정
	statusCode := http.StatusOK
	if !response.Success {
		switch response.Code {
		case "RESTAURANT_NOT_FOUND", "WEEK_DATA_NOT_FOUND":
			statusCode = http.StatusNotFound
		case "INVALID_DATE_FORMAT", "MISSING_RESTAURANT_ID", "MISSING_DATE":
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
	}

	c.JSON(statusCode, response)
}
