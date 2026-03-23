package main

import (
	"auth-system/config"
	"auth-system/database"
	"auth-system/handlers"
	"auth-system/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	database.ConnectDB(cfg)

	// Инициализация сервисов
	smsService := handlers.NewSMSService(cfg)
	emailService := handlers.NewEmailService(cfg)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(cfg, smsService, emailService)

	// Настройка роутера
	r := gin.Default()

	// Публичные маршруты
	auth := r.Group("/api/auth")
	{
		auth.POST("/send-sms", authHandler.SendSMS)
		auth.POST("/verify-sms", authHandler.VerifySMS)
		auth.POST("/send-email", authHandler.SendEmailCode)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Защищенные маршруты
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/logout", authHandler.Logout)
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
