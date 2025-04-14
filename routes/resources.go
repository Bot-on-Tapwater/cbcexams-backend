package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResourceRoutes(r *gin.Engine, db *gorm.DB) {
	resourceCtrl := controllers.ResourceController{DB: db}

	resources := r.Group("v1/api/resources")
	{
		resources.GET("/", resourceCtrl.GetResources)
		resources.GET("/parent-directories", resourceCtrl.GetUniqeParentDirectories)
	}
}
