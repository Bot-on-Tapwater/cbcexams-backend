package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func WebDevRoutes(r *gin.Engine, db *gorm.DB) {
	webdevCtrl := controllers.WebDevController{DB: db}
	v1 := r.Group("v1/api/webdev")

	{
		v1.POST("/requests", webdevCtrl.CreateRequest)
		v1.GET("/requests", webdevCtrl.GetRequests)
	}
}
