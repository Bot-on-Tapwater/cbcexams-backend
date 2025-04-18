package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"github.com/bot-on-tapwater/cbcexams-backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.AuthController{DB: db}

	auth := r.Group("/v1/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password", authController.ResetPassword)
	}

	protected := r.Group("/v1/auth")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/logout", authController.Logout)
		protected.GET("/check", authController.CheckAuth)
	}
}
