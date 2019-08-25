package controller

import (
	"fmt"

	"../service"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

/*
 * TODO: エラハンでError構造体をjsonで返す
 */
func (usercontroller UserController) Shupple(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetOpponent(c)

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}
/*
 * TODO: エラハンでError構造体をjsonで返す
 */
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
/*
 * User情報更新
 */
func (userController UserController) Update(c *gin.Context)  {
	var userService service.UserService
	p, err := userService.Update(c)

	if err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}
/*
 * User取得
 */
func (userController UserController) GetUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetUser(c)

	if err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}
