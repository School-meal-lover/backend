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
// @Param id path string true "Restaurant ID"
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
