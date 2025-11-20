package handlers

import (
	"net/http"
	"strconv"

	"github.com/School-meal-lover/backend/internal/models"
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
// @Param restaurant_name query int true "레스토랑 번호" Enums(1, 2, 3)
// @Param  data body models.ImageUploadRequest true "업로드할 이미지 이름"
// @Success      200 {object} models.ImageInfoResponse "성공적으로 이미지 이름 업로드"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /images/upload [post]
func (h *ImageHandler) UploadImageName(c *gin.Context) {
	var requestBody models.ImageUploadRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	restaurantNumberString := c.Query("restaurant_name")
	if restaurantNumberString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_name query parameter is required"})
		return
	}
	restaurantNumber, err := strconv.Atoi(restaurantNumberString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant_name query parameter"})
		return
	}

	imageName := requestBody.ImageName
	response, err := h.imageService.UploadImageName(imageName, restaurantNumber-1)
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
// @Param restaurant_name query int true "레스토랑 번호" Enums(1, 2, 3)
// @Success      200 {object} models.ImageInfoResponse "성공적으로 현재 이미지 이름 조회"
// @Failure      500 {object} models.ErrorResponse "서버 내부 오류 발생"
// @Router       /images/current [get]
func (h *ImageHandler) GetCurrentImageName(c *gin.Context) {
	restaurantNumberString := c.Query("restaurant_name")
	if restaurantNumberString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_name query parameter is required"})
		return
	}
	restaurantNumber, err := strconv.Atoi(restaurantNumberString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant_name query parameter"})
		return
	}

	response, err := h.imageService.GetCurrentImageName(restaurantNumber - 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
