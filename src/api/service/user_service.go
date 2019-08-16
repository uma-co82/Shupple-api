package service

import (
	"fmt"

	"../db"
	"../model"
	"github.com/gin-gonic/gin"
)

type UserService struct{}

type User model.User

// MEMO: テスト用確認できたら消して良い
func (s UserService) GetAll() ([]User, error) {
	db := db.GetDB()
	var users []User

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s UserService) CreateUser(c *gin.Context) (User, error) {
	db := db.GetDB()
	var user User

	if err := c.BindJSON(&user); err != nil {
		fmt.Println("BindJSONError")
		fmt.Printf("%v", err)
		return user, err
	}

	if err := db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
