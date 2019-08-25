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

	u := r.Group("/users")
	{
		ctrl := controller.UserController{}
		u.GET("", ctrl.Shupple)
		u.POST("", ctrl.Create)
		u.GET("/select", ctrl.GetUser)
		u.PUT("/self", ctrl.Update)
	}

	return r
}
