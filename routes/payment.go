package routes

import (
	"github.com/bot-on-tapwater/cbcexams-backend/controllers"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine) {
	paymentCtrl := controllers.NewPaymentController()

	payment := r.Group("v1/api/payments")
	{
		payment.POST("/initiate", paymentCtrl.InitiatePayment)
		payment.GET("/status/:order_id", paymentCtrl.CheckPaymentStatus)
	}
}
