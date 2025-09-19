package main

import (
	"github.com/School-meal-lover/backend/app/internal/database"
	"github.com/School-meal-lover/backend/app/internal/handlers"
	"github.com/School-meal-lover/backend/app/internal/repository"
	"github.com/School-meal-lover/backend/app/internal/services"

	_ "github.com/School-meal-lover/backend/docs" // Swagger 문서
	gin "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Grrrrr API
// @version 1.0
// @host api.grrrr.me
// @description The server for Grrrrr application.
// @BasePath /api/v1
// @schemes http
func main() {
	router := gin.Default()

	// DB 연결
	database.ConnectDatabase()
	db := database.Db

	// 의존성 주입
	mealRepo := repository.NewMealRepository(db)

	// 서비스 초기화
	mealService := services.NewMealService(mealRepo)
	excelService := services.NewExcelService(mealRepo)

	// 핸들러 초기화
	mealHandler := handlers.NewMealHandler(mealService)
	excelHandler := handlers.NewExcelHandler(excelService)

	// CORS 미들웨어
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 라우트
	api := router.Group("/api/v1")
	{
		api.GET("/restaurants/:id", mealHandler.GetRestaurantMeals)

		api.POST("/upload/excel", excelHandler.UploadAndProcessExcel)

		api.GET("/process/excel/local", excelHandler.ProcessLocalExcel) // 개발용
	}
	// Set up Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	router.Run(":8080")
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
