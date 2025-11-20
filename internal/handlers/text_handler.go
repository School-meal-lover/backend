package handlers

import (
	"net/http"

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

// @Summary Plain text로 한 주치 식단 입력 API
// @Description Plain text 형식으로 한 주치 식단 데이터를 입력받아 데이터베이스에 저장합니다.
// @Tags text
// @Accept text/plain
// @Produce json
// @Param text body string true "식단 텍스트 데이터" example:"식당: RESTAURANT_1\n주 시작일: 2025-06-28\n\n월요일 (2025-06-28)\n아침:\n밥: 쌀밥\n국: 된장국\n반찬: 김치\n\n점심1:\n메인메뉴: 제육볶음\n\n점심2:\n밥: 쌀밥\n국: 미역국\n메인메뉴: 돈까스\n\n저녁:\n밥: 쌀밥\n국: 콩나물국\n메인메뉴: 불고기"
// @Success 200 {object} models.ExcelProcessResult "Text data processed successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid text format"
// @Failure 500 {object} models.ErrorResponse "Failed to process text data"
// @Router /upload/text [post]
func (h *TextHandler) UploadAndProcessText(c *gin.Context) {
	// Request body에서 plain text 읽기
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Failed to read request body: " + err.Error(),
		})
		return
	}

	text := string(body)
	if text == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Text data is empty",
		})
		return
	}

	// 텍스트 처리
	result, err := h.textService.ProcessTextData(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to process text data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

