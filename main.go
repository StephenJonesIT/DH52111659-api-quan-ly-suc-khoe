package main

import (
	"DH52111659-api-quan-ly-suc-khoe/config"
	_ "DH52111659-api-quan-ly-suc-khoe/docs"
	"DH52111659-api-quan-ly-suc-khoe/internal/handlers"
	"DH52111659-api-quan-ly-suc-khoe/internal/middleware"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"DH52111659-api-quan-ly-suc-khoe/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Healthy Service API Document
// @version 1.0
// @description List APIs of Healthy Management Service
// @termsOfService http://swagger.io/terms/

// @host 127.0.0.1:9000
// @BasePath /api/v1
// @schemes http https
func main() {
	config.LoadConfig()
	repositories.ConnectDB()
	redis, err := repositories.NewRedisStore()
	if err != nil {
		panic(err)
	}
	
	// 3. Khởi tạo router
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(gin.Logger()) // Log requests
	router.Use(gin.Recovery()) // Recover from panics and log them

	// 4. Khởi tạo các service và handler
	accountRepo := repositories.NewAccountRepoImpl(repositories.DB)
	accountService := services.NewAccountServiceImpl(accountRepo, redis)
	accountHandler := handlers.NewAccountHandler(accountService)

	// 5. Đăng ký các route
	registerRouter(router, accountHandler)

	// 6. Khởi động server
	if err := router.Run(":"+config.AppConfig.GinPort); err != nil {
		panic(err)
	}
}

func registerRouter(router *gin.Engine, accountHandler *handlers.AccountHandler) {
	// Tạo một nhóm router cho API
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		// Đăng ký các route cho tài khoản
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/refresh-token", accountHandler.RefreshTokenHandler)
			authGroup.POST("/register", accountHandler.RegisterAccountHandler)
			authGroup.POST("/verify-email", accountHandler.RegisterVerifyOTPHandler)
			authGroup.POST("/login", accountHandler.LoginHandler)
			authGroup.POST("/forgot-password", accountHandler.ForgotPasswordHandler)
			authGroup.POST("/verify-otp", accountHandler.VerifyOTPHandler)
			authGroup.POST("/reset-password", accountHandler.ResetPasswordHandler)
			
			authGroup.Use(middleware.JWTAuthMiddleware(*utils.NewTokenService(config.AppConfig.SECRET_KEY)))
			authGroup.POST("/change-password", accountHandler.ChangePasswordHandler)
		}
	}
}