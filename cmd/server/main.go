package main

import (
	"log"

	docs "github.com/School-meal-lover/backend/docs"
	"github.com/School-meal-lover/backend/internal/database"
	"github.com/School-meal-lover/backend/internal/handlers"
	"github.com/School-meal-lover/backend/internal/middleware"
	"github.com/School-meal-lover/backend/internal/repository"
	"github.com/School-meal-lover/backend/internal/services"
	"github.com/joho/godotenv"

	gin "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Grrrrr API
// @version 1.0
// @host api.grrrr.me
// @description The server for Grrrrr application.
// @BasePath /api/v1
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 토큰 인증. 토큰은 환경변수 BEARER_TOKEN에서 설정.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found; continuing with environment variables")
	}
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	// DB 연결
	database.ConnectDatabase()
	db := database.Db

	// 의존성 주입
	mealRepo := repository.NewMealRepository(db)

	// 서비스 초기화
	mealService := services.NewMealService(mealRepo)
	excelService := services.NewExcelService(mealRepo)
	textService := services.NewTextService(mealRepo)
	imageService := services.NewImageService()

	// 핸들러 초기화
	mealHandler := handlers.NewMealHandler(mealService)
	excelHandler := handlers.NewExcelHandler(excelService)
	textHandler := handlers.NewTextHandler(textService)
	imageHandler := handlers.NewImageHandler(imageService)

	// CORS 미들웨어
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 라우트
	api := router.Group("/api/v1")
	{
		api.GET("/restaurants/:name", mealHandler.GetRestaurantMeals)

		api.POST("/upload/excel", excelHandler.UploadAndProcessExcel)

		// Bearer token 인증이 필요한 엔드포인트
		api.POST("/upload/text", middleware.BearerTokenAuth(), textHandler.UploadText)

		api.POST("/images/upload", imageHandler.UploadImageName)
		api.GET("/images/current", imageHandler.GetCurrentImageName)
	}
	// Set up Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	router.Run(":8080")
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
