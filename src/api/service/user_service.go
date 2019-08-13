package service

import (
	"github.com/gin-gonic/gin"
	"github.com/holefillingco-ltd/Shupple-api/src/api/db"
	"github.com/holefillingco-ltd/Shupple-api/src/api/model"
)

type Service struct{}

type User model.User

func (s Service) GetAll() ([]User, error) {
	db := db.GetDB()
	var users []User

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s Service) CreateUser(c *gin.Context) (User, error) {
	db := db.GetDB()
	var user User

	if err := c.BindJSON(&user); err != nil {
		return user, err
	}

	if err := db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
