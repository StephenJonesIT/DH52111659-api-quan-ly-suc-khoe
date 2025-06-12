package main

import (
	"DH52111659-api-quan-ly-suc-khoe/config"
	_ "DH52111659-api-quan-ly-suc-khoe/docs"
	"DH52111659-api-quan-ly-suc-khoe/internal/handlers"
	"DH52111659-api-quan-ly-suc-khoe/internal/middleware"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"DH52111659-api-quan-ly-suc-khoe/utils"
 	"github.com/gin-contrib/cors"
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
	expertService := services.NewExpertService(expertRepo, accountRepo)
	expertHandler := handlers.NewExpertHandler(expertService)

	programRepo := repositories.NewProgramRepository(repositories.DB)
	programService := services.NewProgramService(programRepo, expertRepo)
	programHandler := handlers.NewProgramHandler(programService)

	activityRepo := repositories.NewActivityRepositoryImpl(repositories.DB)
	activityService := services.NewActivityService(activityRepo, expertRepo)
	activityHandler := handlers.NewActivityHandler(activityService)

	scheduleRepo := repositories.NewScheduleRepositoryImpl(repositories.DB)
	scheduleService := services.NewScheduleServiceImpl(scheduleRepo, programRepo, activityRepo)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)

	// 5. Đăng ký các route
	registerRouter(router, authHandler, profileHandler, userHandler, expertHandler, programHandler, activityHandler, scheduleHandler)

	// 6. Khởi động server
	if err := router.Run("0.0.0.0:"+config.AppConfig.GinPort); err != nil {
		panic(err)
	}
}

func registerRouter(
	router *gin.Engine, 
	accountHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	userHandler *handlers.UserHandler,
	expertHandler *handlers.ExpertHandler,
	programHandler *handlers.ProgramHandler,
	activityHandler *handlers.ActivityHandler,
	scheduleHandler *handlers.ScheduleHandler,
	) {
	// Tạo một nhóm router cho API
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Change "*" to specific domains for security
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))


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
				protected.PUT("/:id", profileHandler.UpdateProfileHandler)
			}
		}

		adminGroup := api.Group("/admin")
		{
			adminGroup.Use(middleware.JWTAuthMiddleware(*utils.NewTokenService(config.AppConfig.SECRET_KEY), "admin"))
			userGroup := adminGroup.Group("/users")
			{
				userGroup.POST("", userHandler.CreateUserHandler)
				userGroup.POST("/reset-password", userHandler.ResetPasswordUserHandler)
				userGroup.GET("", userHandler.GetUsersHandler)
				userGroup.GET("/:id", userHandler.GetUserByIdHandler)
				userGroup.PATCH("/:id/lock", userHandler.LockUserAccountHandler)
				userGroup.PATCH("/:id/unlock", userHandler.UnlockUserAccountHandler)
			}

			expertGroup := adminGroup.Group("experts")
			{
				expertGroup.POST("",expertHandler.CreateExpertHandler)
				expertGroup.GET("", expertHandler.GetExpertsHandler)
				expertGroup.PUT("/:id", expertHandler.UpdateExpertHandler)
				expertGroup.DELETE("/:id", expertHandler.DeleteExpertHandler)
				
				accountGroup := expertGroup.Group("/accounts")
				{
					accountGroup.POST("", userHandler.CreateExpertAccountHandler)
					accountGroup.POST("/reset-password", userHandler.ResetPasswordExpertHandler)
					accountGroup.PATCH("/:id/lock", userHandler.LockExpertAccountHandler)
					accountGroup.PATCH("/:id/unlock", userHandler.UnlockExpertAccountHandler)
				}
			}
		}

		expertGroup := api.Group("/expert")
		{
			expertGroup.Use(middleware.JWTAuthMiddleware(*utils.NewTokenService(config.AppConfig.SECRET_KEY), "expert"))
			program := expertGroup.Group("/programs")
			{
				program.POST("", programHandler.CreateProgramHandler)

				schedule := program.Group("/:program_id/schedules")
				{
					schedule.POST("", scheduleHandler.CreateScheduleHandler)
				}
			}

			activity := expertGroup.Group("/activities")
			{
				activity.POST("", activityHandler.CreateActivityHandler)
				activity.GET("", activityHandler.GetAtivitiesExpertHandler)
			}
		}
	}
}