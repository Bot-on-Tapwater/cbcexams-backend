package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func JobRoutes(r *gin.Engine, db *gorm.DB) {
	jobsCtrl := controllers.JobsController{DB: db}
	v1 := r.Group("v1/api/jobs")

	{
		/* School job listings */
		v1.POST("/schools", jobsCtrl.CreateSchoolJob)
		v1.GET("/schools", jobsCtrl.GetSchoolJobs)

		/* Teacher profiles */
		v1.POST("/teachers", jobsCtrl.CreateTeacherProfile)
		v1.GET("/teachers", jobsCtrl.GetTeacherProfiles)
	}
}
