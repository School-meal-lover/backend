package services

import (
	"time"

	"github.com/School-meal-lover/backend/internal/models"
)

type ImageService struct {
	current_image_name [3]string
	current_image_date [3]time.Time
}

func NewImageService() *ImageService {
	return &ImageService{
		current_image_name: [3]string{"", "", ""},
		current_image_date: [3]time.Time{},
	}
}

// 이미지 이름 업로드 및 저장
func (s *ImageService) UploadImageName(imageName string, restaurant_num int) (*models.ImageInfoResponse, error) {
	s.current_image_name[restaurant_num] = imageName
	s.current_image_date[restaurant_num] = time.Now()
	return &models.ImageInfoResponse{
		Success:   true,
		ImageName: s.current_image_name[restaurant_num],
		ImageDate: s.current_image_date[restaurant_num].Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *ImageService) GetCurrentImageName(restaurant_num int) (*models.ImageInfoResponse, error) {
	return &models.ImageInfoResponse{
		Success:   true,
		ImageName: s.current_image_name[restaurant_num],
		ImageDate: s.current_image_date[restaurant_num].Format("2006-01-02 15:04:05"),
	}, nil
}
