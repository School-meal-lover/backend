package handlers

import (
	"net/http"
	"strings"
	"time"

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

// @Summary      특정 식당의 주간 식단 조회
// @Description  경로 파라미터로 받은 식당 이름과 쿼리로 받은 날짜를 기준으로 주간 식단을 조회합니다. Restaurant1은 평일만(월~금, 5일), Restaurant2는 주말 포함(월~일, 7일) 조회됩니다.
// @Tags         Meals
// @Accept       json
// @Produce      json
// @Param        name path string true "레스토랑 이름 (RESTAURANT_1: 평일만, RESTAURANT_2: 주말 포함)" example:"RESTAURANT_1 대소문자 관계없음"
// @Param        date query string true "조회할 날짜 (YYYY-MM-DD 형식)" example:"2025-06-28"
// @Success      200 {object} models.RestaurantMealsResponse "성공적으로 식단 정보 조회"
// @Failure      400 {object} models.ErrorResponse "잘못된 요청 파라미터 (식당 이름 또는 날짜 형식 오류)"
// @Failure      404 {object} models.ErrorResponse "해당 식당 또는 해당 날짜의 식단 정보를 찾을 수 없음"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /restaurants/{name} [get]
func (h *MealHandler) GetRestaurantMeals(c *gin.Context) {
	restaurantName := strings.ToUpper(c.Param("name"))
	date := c.Query("date")

	// 파라미터 검증
	if restaurantName == "" {
		c.JSON(http.StatusBadRequest, models.RestaurantMealsResponse{
			Success: false,
			Error:   "restaurant name is required in path",
			Code:    "MISSING_RESTAURANT_NAME",
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
	response, err := h.mealService.GetRestaurantWeekMeals(restaurantName, date)
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

// @Summary      개별 메뉴 품절 정보 조회
// @Description  개별 메뉴의 품절 정보를 조회합니다. 메뉴가 품절되었는지 여부와 품절 시각을 반환합니다.
// @Tags         Meals
// @Accept       json
// @Produce      json
// @Param        meal_id path string true "메뉴 ID" example:"123e4567-e89b-12d3-a456-426614174000"
// @Success      200 {object} models.IndMenuSoldResponse "성공적으로 개별 메뉴 품절 정보 조회"
// @Failure      400 {object} models.ErrorResponse "잘못된 요청 파라미터 (메뉴 ID 누락)"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /meals/{meal_id}/ind_menu_sold [get]
func (h *MealHandler) GetIndMenuSold(c *gin.Context) {
	mealID := c.Param("meal_id")
	if mealID == "" {
		c.JSON(http.StatusBadRequest, models.IndMenuSoldResponse{
			Success: false,
			Error:   "meal_id parameter is required",
			Code:    "MISSING_MEAL_ID",
		})
		return
	}
	response, err := h.mealService.GetIndMenuSold(mealID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.IndMenuSoldResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary      개별 메뉴 품절 정보 업데이트
// @Description  개별 메뉴의 품절 정보를 업데이트합니다. 품절 시각을 설정하거나 품절 해제를 할 수 있습니다.
// @Tags         Meals
// @Accept       json
// @Produce      json
// @Param        meal_id path string true "메뉴 ID" example:"123e4567-e89b-12d3-a456-426614174000"
// @Param        sold_out_at body string false "품절 시각 (ISO 8601 형식, null이면 품절 해제)" example:"2024-06-30T15:04:05Z07:00"
// @Success      200 {object} models.IndMenuSoldResponse "성공적으로 개별 메뉴 품절 정보 업데이트"
// @Failure      400 {object} models.ErrorResponse "잘못된 요청 파라미터 (메뉴 ID 누락 또는 잘못된 품절 시각 형식)"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /meals/{meal_id}/ind_menu_sold [put]
func (h *MealHandler) UpdateIndMenuSold(c *gin.Context) {
	mealID := c.Param("meal_id")
	if mealID == "" {
		c.JSON(http.StatusBadRequest, models.IndMenuSoldResponse{
			Success: false,
			Error:   "meal_id parameter is required",
			Code:    "MISSING_MEAL_ID",
		})
		return
	}
	var reqBody struct {
		SoldOutAt *string `json:"sold_out_at"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.IndMenuSoldResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST_BODY",
		})
		return
	}
	var soldOutAt *time.Time
	if reqBody.SoldOutAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *reqBody.SoldOutAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.IndMenuSoldResponse{
				Success: false,
				Error:   "Invalid sold_out_at format, must be ISO 8601",
				Code:    "INVALID_SOLD_OUT_AT_FORMAT",
			})
			return
		}
		soldOutAt = &parsedTime
	}
	response, err := h.mealService.MarkIndMenuSoldOut(mealID, soldOutAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.IndMenuSoldResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    "INTERNAL_ERROR",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
