package handlers

import (
	"net/http"

	"github.com/School-meal-lover/backend/internal/models"
	"github.com/School-meal-lover/backend/internal/services"
	"github.com/School-meal-lover/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type MealHandler struct {
	mealService *services.MealService
}

func NewMealHandler(mealService *services.MealService) *MealHandler {
	return &MealHandler{mealService: mealService}
}

// @Summary      특정 식당의 주간 식단 조회
// @Description  경로 파라미터로 받은 식당 타입과 쿼리로 받은 날짜를 기준으로 주간 식단을 조회합니다. 식당 1은 월화수목금(5일), 식당 2는 월화수목금토일(7일)입니다.
// @Tags         Meals
// @Accept       json
// @Produce      json
// @Param        restaurant_type path string true "식당 타입 (RESTAURANT_1 또는 RESTAURANT_2)" example:"RESTAURANT_1"
// @Param        date query string true "조회할 날짜 (YYYY-MM-DD 형식)" example:"2025-06-28"
// @Success      200 {object} models.RestaurantMealsResponse "성공적으로 식단 정보 조회"
// @Failure      400 {object} models.ErrorResponse "잘못된 요청 파라미터 (식당 타입 또는 날짜 형식 오류)"
// @Failure      404 {object} models.ErrorResponse "해당 식당 또는 해당 날짜의 식단 정보를 찾을 수 없음"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /restaurants/{restaurant_type}/meals [get]
// @Router       /restaurants/{name} [get]
func (h *MealHandler) GetRestaurantMeals(c *gin.Context) {
	// 경로 파라미터에서 식당 타입 파싱 (utils 사용)
	restaurantTypeParam := c.Param("restaurant_type")
	if restaurantTypeParam == "" {
		// 기존 호환성을 위해 name 파라미터도 지원
		restaurantTypeParam = c.Param("name")
	}
	
	restaurantType, err := utils.ParseRestaurantType(restaurantTypeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.RestaurantMealsResponse{
			Success: false,
			Error:   "Invalid restaurant type: " + err.Error(),
			Code:    "INVALID_RESTAURANT_TYPE",
		})
		return
	}
	
	date := c.Query("date")

	// 파라미터 검증
	if date == "" {
		c.JSON(http.StatusBadRequest, models.RestaurantMealsResponse{
			Success: false,
			Error:   "date parameter is required (YYYY-MM-DD)",
			Code:    "MISSING_DATE",
		})
		return
	}

	// 서비스 호출 (식당 타입 전달)
	response, err := h.mealService.GetRestaurantWeekMeals(restaurantType, date)
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
