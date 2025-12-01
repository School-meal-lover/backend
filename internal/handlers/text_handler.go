package handlers

import (
	"net/http"
	"strings"

	"github.com/School-meal-lover/backend/internal/models"
	"github.com/School-meal-lover/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type TextHandler struct {
	textService *services.TextService
}

func NewTextHandler(textService *services.TextService) *TextHandler {
	return &TextHandler{textService: textService}
}

// @Summary 텍스트로 식단 데이터 업로드
// @Description json 한 주치 식단 데이터를 받아서 디비에 저장합니다. Bearer token 인증이 필요합니다.
// @Tags text
// @Accept json
// @Produce json
// @Security BearerAuth
// @query restaurant_name int true "레스토랑 이름" Enums(RESTAURANT_1, RESTAURANT_2)
// @Param text body string true "식단 텍스트 데이터" example:"RESTAURANT_1\n2025-05-26\nMonday 2025-05-26\nBreakfast\n밥\n국\n반찬\nLunch_1\n메인메뉴\nLunch_2\n밥\n국\n메인메뉴\n반찬\nDinner\n밥\n국\n메인메뉴\n반찬"
// @Success 200 {object} models.ExcelProcessResult "Text processed successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Failed to process text"
// @Router /upload/text [post]
func (h *TextHandler) UploadText(c *gin.Context) {
	restaurantName := strings.ToUpper(c.Query("restaurant_name"))
	if restaurantName != "RESTAURANT_1" && restaurantName != "RESTAURANT_2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant_name query parameter"})
		return
	}
	var requestBody models.TextUploadRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.textService.ProcessText(restaurantName, requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
