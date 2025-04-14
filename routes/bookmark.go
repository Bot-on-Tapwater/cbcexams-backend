package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"github.com/bot-on-tapwater/cbcexams-backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BookmarkRoutes(r *gin.Engine, db *gorm.DB) {
	bookmarkCtrl := controllers.BookmarkController{DB: db}

	protected := r.Group("/v1/api/bookmarks")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/", bookmarkCtrl.CreateBookmark)
		protected.DELETE("/:resource_id", bookmarkCtrl.DeleteBookmark)
		protected.GET("/", bookmarkCtrl.GetUserBookmarks)
	}
}
