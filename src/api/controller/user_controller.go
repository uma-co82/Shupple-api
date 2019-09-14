package controller

import (
	"../service"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

/*
 * TODO: エラハンでError構造体をjsonで返す
 */
func (userController UserController) Shupple(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetOpponent(c)

	if err != nil {
		c.JSON(404, err)
	} else {
		c.JSON(200, p)
	}
}

/*
 * TODO: エラハンでError構造体をjsonで返す
 */
func (userController UserController) CreateUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.CreateUser(c)

	if err != nil {
		c.JSON(400, err)
	} else {
		c.JSON(200, p)
	}
}

/*
 * User情報更新
 */
func (userController UserController) UpdateUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.UpdateUser(c)

	if err != nil {
		c.JSON(400, err)
	} else {
		c.JSON(200, p)
	}
}

/**
 * Userが登録済みか調べる
 */
func (userController UserController) IsRegisteredUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.IsRegisterdUser(c)

	if err != nil {
		// TODO: エラハン
	} else {
		c.JSON(200, p)
	}
}

/**
 * Userがマッチング済みか判定
 * 済みの場合はマッチング相手を返す
 */
func (userController UserController) IsMatchedUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.IsMatchedUser(c)

	if err != nil {
		// TODO: エラハン
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
	} else {
		c.JSON(200, p)
	}
}

/**
 * 相性が良い条件の組み合わせを保存
 */
func (userController UserController) CreateCompatible(c *gin.Context) {
	var userService service.UserService
	p, err := userService.CreateCompatible(c)

	if err != nil {
		c.AbortWithStatus(400)
	} else {
		c.JSON(200, p)
	}
}
