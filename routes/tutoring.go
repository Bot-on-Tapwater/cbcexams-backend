package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func TutoringRoutes(r *gin.Engine, db *gorm.DB) {
	tutoringCtrl := controllers.TutoringController{DB: db}
	v1 := r.Group("v1/api/tutoring")

	{
		/* Tutor Requests (Students/Parents) */
		v1.POST("/requests", tutoringCtrl.CreateTutorRequest)
		v1.GET("/requests", tutoringCtrl.GetTutorRequests)

		/* Tutor Applications (Tutors) */
		v1.POST("/applications", tutoringCtrl.CreateTutorApplication)
		v1.GET("/applications", tutoringCtrl.GetTutorApplications)
	}
}
