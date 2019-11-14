package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/uma-co82/Shupple-api/src/api/domain"
	"github.com/uma-co82/Shupple-api/src/api/domain/service"
)

type UserController struct{}

/*
 * マッチング！
 */
func (userController UserController) Shupple(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetOpponent(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, p)
	}
}

/**
 * マッチング解除
 */
func (userController UserController) CancelOpponent(c *gin.Context) {
	var userService service.UserService
	p, err := userService.CancelOpponent(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, p)
	}
}

/*
 * User登録
 */
func (userController UserController) CreateUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.CreateUser(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
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
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, p)
	}
}

/**
 * Userが登録済みか調べる
 */
func (userController UserController) IsRegisteredUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.IsRegisteredUser(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
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
		c.JSON(err.(*domain.Error).Code, err)
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
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, p)
	}
}

/**
 * User論理削除
 */
func (userController UserController) SoftDeleteUser(c *gin.Context) {
	var userService service.UserService
	err := userService.SoftDeleteUser(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, "ok")
	}
}

/**
 * Userからのブロック、報告機能
 */
func (userController UserController) UnauthorizedUser(c *gin.Context) {
	var userService service.UserService
	err := userService.CreateUnauthorizedUser(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, "ok")
	}
}

/**
 * 相性が良い条件の組み合わせを保存
 */
func (userController UserController) CreateCompatible(c *gin.Context) {
	var userService service.UserService
	p, err := userService.CreateCompatible(c)

	if err != nil {
		c.JSON(err.(*domain.Error).Code, err)
	} else {
		c.JSON(200, p)
	}
}
