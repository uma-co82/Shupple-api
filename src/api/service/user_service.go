package service

import (
	"fmt"

	"../db"
	"../model"
	"github.com/gin-gonic/gin"
)

type UserService struct{}

type PostUser model.PostUser
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
	var postUser PostUser
	var user User

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
	}

	// TODO: メソッドにする
	user.UID = model.GetUUID()
	user.NickName = postUser.NickName
	user.Sex = postUser.Sex
	user.Hobby = postUser.Hobby
	user.Age, _ = model.CalcAge(postUser.BirthDay)

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("DBError")
		fmt.Printf("%v", err)
		return user, err
	}

	return user, nil
}
