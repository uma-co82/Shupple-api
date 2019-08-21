package controller

import (
	"fmt"

	"../service"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// TODO: エラハンでError構造体をjsonで返す
func (usercontroller UserController) Shapple(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetOpponent(c)

	if err != nil {
		// TODO: ここ！
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}

// TODO: エラハンでError構造体をjsonで返す
func (usercontroller UserController) Create(c *gin.Context) {
	var userService service.UserService
	p, err := userService.CreateUser(c)

	if err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}

func (userController UserController) Self(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetSelfUser(c)

	if err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}
