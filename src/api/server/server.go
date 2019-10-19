package server

import (
	//"../controller"
	"github.com/gin-gonic/gin"
	"github.com/uma-co82/Shupple-api/src/api/controller"
)

func Init() {
	r := router()
	r.Run()
}

func router() *gin.Engine {
	r := gin.Default()

	/**
	 * ヘルスチェクに対するレスポンス
	 */
	hGroup := r.Group("/")
	{
		ctrl := controller.HealthCheckController{}
		hGroup.GET("", ctrl.HealthCheck)
	}

	uGroup := r.Group("/users")
	{
		ctrl := controller.UserController{}
		uGroup.GET("/shupple", ctrl.Shupple)
		uGroup.PUT("/shupple", ctrl.CancelOpponent)
		uGroup.POST("", ctrl.CreateUser)
		uGroup.PUT("", ctrl.UpdateUser)
		uGroup.DELETE("", ctrl.SoftDeleteUser)
		uGroup.GET("/isRegistered", ctrl.IsRegisteredUser)
		uGroup.GET("/isMatched", ctrl.IsMatchedUser)
		uGroup.GET("/select", ctrl.GetUser)
		uGroup.POST("/unauthorized", ctrl.UnauthorizedUser)
		uGroup.POST("/compatible", ctrl.CreateCompatible)
	}

	aGroup := r.Group("/admin")
	{
		ctrl := controller.AdminController{}
		aGroup.GET("/users", ctrl.GetAllUser)
	}

	return r
}
