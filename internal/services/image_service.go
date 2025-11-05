package services

import (
	"time"

	"github.com/School-meal-lover/backend/internal/models"
)

type ImageService struct {
	current_image_name string
	current_image_date time.Time
}

func NewImageService() *ImageService {
	return &ImageService{
		current_image_name: "",
		current_image_date: time.Time{},
	}
}

// 이미지 이름 업로드 및 저장
func (s *ImageService) UploadImageName(imageName string) (*models.ImageInfoResponse, error) {
	s.current_image_name = imageName
	s.current_image_date = time.Now()
	return &models.ImageInfoResponse{
		Success:   true,
		ImageName: imageName,
		ImageDate: s.current_image_date.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *ImageService) GetCurrentImageName() (*models.ImageInfoResponse, error) {
	return &models.ImageInfoResponse{
		Success:   true,
		ImageName: s.current_image_name,
		ImageDate: s.current_image_date.Format("2006-01-02 15:04:05"),
	}, nil
}
