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

// @Summary 텍스트로 식단 데이터 업로드
// @Description plain text로 한 주치 식단 데이터를 받아서 디비에 저장합니다. Bearer token 인증이 필요합니다.
// @Tags text
// @Accept text/plain
// @Produce json
// @Security BearerAuth
// @Param text body string true "식단 텍스트 데이터" example:"RESTAURANT_1\n2025-05-26\nMonday 2025-05-26\nBreakfast\n밥\n국\n반찬\nLunch_1\n메인메뉴\nLunch_2\n밥\n국\n메인메뉴\n반찬\nDinner\n밥\n국\n메인메뉴\n반찬"
// @Success 200 {object} models.ExcelProcessResult "Text processed successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Failed to process text"
// @Router /upload/text [post]
func (h *TextHandler) UploadText(c *gin.Context) {
	// 텍스트 데이터 읽기
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Failed to read request body",
		})
		return
	}

	text := string(body)
	if text == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Text body is required",
		})
		return
	}

	// 텍스트 처리
	result, err := h.textService.ProcessText(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to process text: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

