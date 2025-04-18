package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func FeedbackRoutes(r *gin.Engine, db *gorm.DB) {
	feedbackCtrl := controllers.FeedbackController{DB: db}
	v1 := r.Group("v1/api/feedback")

	{
		v1.POST("", feedbackCtrl.SubmitFeedback)
		v1.GET("", feedbackCtrl.GetFeedback)

	}
}
