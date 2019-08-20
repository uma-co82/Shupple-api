package controller

import (
	"fmt"

	"../service"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (usercontroller UserController) Shapple(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetOpponent(c)

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}

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
