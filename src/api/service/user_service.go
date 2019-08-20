package service

import (
	"fmt"
	"math/rand"
	"time"

	"../db"
	"../model"
	"github.com/gin-gonic/gin"
)

type UserService struct{}

type PostUser model.PostUser
type User model.User
type UserInformation model.UserInformation
type UserCombination model.UserCombination
type Profile model.Profile

type UID struct {
	uid string `json:"uid"`
}

// ランダム取得
func getRandUser(u []User) User {
	if u == nil {
		return User{}
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(u))
	return u[i]
}

func (s UserService) GetOpponent(c *gin.Context) (User, error) {
	db := db.GetDB()
	var users []User
	var user User

	uid := c.Request.Header.Get("uid")

	user.UID = uid

	if err := db.First(&user).Error; err != nil {
		return user, err
	}

	opponentSex := user.opponentSex()

	if err := db.Find(&users, "sex=?", opponentSex).Error; err != nil {
		return user, err
	}
	fmt.Println(users)

	// WARN: usersがnullの時エラる
	opponent := getRandUser(users)

	return opponent, nil
}

func (s UserService) CreateUser(c *gin.Context) (User, error) {
	db := db.GetDB()
	var postUser PostUser
	var user User

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
		fmt.Printf("Binding Error %v", err)
		return user, err
	}

	// Firebase UID
	user.UID = postUser.UID
	user.calcAge(postUser.BirthDay)
	user.NickName = postUser.NickName
	user.Sex = postUser.Sex
	user.Hobby = postUser.Hobby
	user.BirthDay = postUser.BirthDay
	user.Residence = postUser.Residence

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return user, err
	}

	return user, nil
}
