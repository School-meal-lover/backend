package handlers

import (
	"net/http"

	"github.com/School-meal-lover/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageService *services.ImageService
}

func NewImageHandler(imageService *services.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

// @Summary      이미지 이름 업로드
// @Description  이미지 이름을 업로드합니다.
// @Tags         Images
// @Accept	   json
// @Produce      json
// @Param  image_name body string true "업로드할 이미지 이름"
// @Success      200 {object} models.ImageInfoResponse "성공적으로 이미지 이름 업로드"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /images/upload [post]
func (h *ImageHandler) UploadImageName(c *gin.Context) {
	var requestBody struct {
		ImageName string `json:"image_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	imageName := requestBody.ImageName
	response, err := h.imageService.UploadImageName(imageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Summary      현재 이미지 이름 조회
// @Description  현재 저장된 이미지 이름을 조회합니다.
// @Tags         Images
// @Accept	   json
// @Produce      json
// @Success      200 {object} models.ImageInfoResponse "성공적으로 현재 이미지 이름 조회"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /images/current [get]
func (h *ImageHandler) GetCurrentImageName(c *gin.Context) {
	response, err := h.imageService.GetCurrentImageName()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
