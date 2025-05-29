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

//	@title			Healthy Service API Document
//	@version		1.0
//	@description	List APIs of Healthy Management Service
//	@termsOfService	http://swagger.io/terms/

//	@host		127.0.0.1:9000
//	@BasePath	/api/v1
//	@schemes	http https
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
	authService := services.NewAuthServiceImpl(accountRepo, redis)
	authHandler := handlers.NewAuthHandler(authService)

	profileRepo := repositories.NewProfileRepoImpl(repositories.DB)
	profileService := services.NewProfileServiceImpl(profileRepo)
	profileHandler := handlers.NewProfileHandler(profileService)

	userService := services.NewUserServiceImpl(accountRepo)
	userHandler := handlers.NewUserHandler(userService)

	expertRepo := repositories.NewExpertRepositoryImpl(repositories.DB)
	expertService := services.NewExpertService(expertRepo)
	expertHandler := handlers.NewExpertHandler(expertService)
	// 5. Đăng ký các route
	registerRouter(router, authHandler, profileHandler, userHandler, expertHandler)

	// 6. Khởi động server
	if err := router.Run(":"+config.AppConfig.GinPort); err != nil {
		panic(err)
	}
}

func registerRouter(
	router *gin.Engine, 
	accountHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	userHandler *handlers.UserHandler,
	expertHandler *handlers.ExpertHandler,
	) {
	// Tạo một nhóm router cho API
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.Static("/uploads/avatar", "./uploads/avatar")
		// Đăng ký các route cho tài khoản
		authGroup := api.Group("/auth")
		{
			public := authGroup.Group("")
			{
				public.POST("/register", accountHandler.RegisterAccountHandler)
				public.POST("/verify-email", accountHandler.RegisterVerifyOTPHandler)
				public.POST("/login", accountHandler.LoginHandler)
				public.POST("/token/refresh", accountHandler.RefreshTokenHandler)
				public.POST("/password/forgot", accountHandler.ForgotPasswordHandler)
				public.POST("/password/verify-otp", accountHandler.VerifyOTPHandler)
				public.POST("/password/reset", accountHandler.ResetPasswordHandler)
			}
	
			protected := authGroup.Group("")
			{
				protected.Use(middleware.JWTAuthMiddleware(*utils.NewTokenService(config.AppConfig.SECRET_KEY)))
				protected.POST("/password/change", accountHandler.ChangePasswordHandler)
			}		
		}

		profileGroup := api.Group("/profile")
		{
			protected := profileGroup.Group("")
			{
				protected.Use(middleware.JWTAuthMiddleware(*utils.NewTokenService(config.AppConfig.SECRET_KEY)))
				protected.POST("",profileHandler.CreateProfileHandler)
				protected.PUT(":id", profileHandler.UpdateProfileHandler)
			}
		}

		adminGroup := api.Group("/admin")
		{
			adminGroup.Use(middleware.JWTAuthMiddleware(*utils.NewTokenService(config.AppConfig.SECRET_KEY), "admin"))
			userGroup := adminGroup.Group("")
			{
				userGroup.POST("/user", userHandler.CreateUserHandler)
				userGroup.POST("/user/reset-password", userHandler.ResetPasswordUserHandler)
				userGroup.GET("/users", userHandler.GetListUserHandler)
				userGroup.GET("/user/:id", userHandler.GetUserByIdHandler)
				userGroup.PATCH("/user/:id/lock", userHandler.LockUserAccountHandler)
				userGroup.PATCH("/user/:id/unlock", userHandler.UnlockUserAccountHandler)
			}

			expertGroup := adminGroup.Group("")
			{
				expertGroup.POST("/expert",expertHandler.CreateExpertHandler)
			}
		}
	}
}