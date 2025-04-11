package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"

	"github.com/gin-gonic/gin"
)

func CategoriesRoutes(r *gin.Engine) {
	v1 := r.Group("v1/api/categories")

	{
		v1.GET("/", controllers.GetCategories)
	}
}
