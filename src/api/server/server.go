package server

import (
	"../controller"
	"github.com/gin-gonic/gin"
)

func Init() {
	r := router()
	r.Run()
}

func router() *gin.Engine {
	r := gin.Default()

	uGroup := r.Group("/users")
	{
		ctrl := controller.UserController{}
		uGroup.GET("/shupple", ctrl.Shupple)
		uGroup.PUT("/shupple", ctrl.CancelOpponent)
		uGroup.POST("", ctrl.CreateUser)
		uGroup.PUT("", ctrl.UpdateUser)
		uGroup.GET("/isRegistered", ctrl.IsRegisteredUser)
		uGroup.GET("/isMatched", ctrl.IsMatchedUser)
		uGroup.GET("/select", ctrl.GetUser)
		uGroup.POST("/compatible", ctrl.CreateCompatible)
	}

	return r
}
