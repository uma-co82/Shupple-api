package server

import (
	"github.com/gin-gonic/gin"
	"github.com/holefillingco-ltd/Shupple-api/src/api/controller"
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
		// テスト用 確認取れれば消して良い
		u.GET("", ctrl.Index)
		u.POST("", ctrl.Create)
	}

	return r
}
