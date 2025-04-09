package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"github.com/bot-on-tapwater/cbcexams-backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UsersRoutes(r *gin.Engine, db *gorm.DB) {
	usersController := controllers.UsersController{DB: db}

	protected := r.Group("/v1/api/users/")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", usersController.Profile)
		protected.GET("/", usersController.GetUsers)
		protected.PATCH("/update-profile", usersController.UpdateProfile)
	}
}