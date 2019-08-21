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
type Error model.Error

type UID struct {
	uid string `json:"uid"`
}

// 引数の[]Userからランダムに1件取得
func getRandUser(u []User) (User, error) {
	user := User{}
	if u == nil {
		err := RaiseError(404, "Opponent Not Found", nil)
		return user, err
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(u))
	user = u[i]
	return user, nil
}

// 異性のUserを返す
func (s UserService) GetOpponent(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var users []User
	var user User
	var profile Profile
	var uInfo UserInformation

	uid := c.Request.Header.Get("uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	opponentSex := user.opponentSex()

	if err := db.Find(&users, "sex=?", opponentSex).Error; err != nil {
		return profile, err
	}

	opponent, err := getRandUser(users)
	if err != nil {
		return profile, err
	}

	if err := db.First(&uInfo, "uid=?", opponent.UID).Error; err != nil {
		return profile, err
	}

	profile.User = model.User(opponent)
	profile.Information = model.UserInformation(uInfo)

	return profile, nil
}

// POSTされたjsonを元にUser, UserInformationを作成
// HACK: どうしても詰め替えの作業が冗長になってる。。ここだけメソッドに任せよう！
func (s UserService) CreateUser(c *gin.Context) (User, error) {
	db := db.GetDB()
	var postUser PostUser
	var user User
	var uInformation UserInformation

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
		fmt.Printf("Binding Error %v", err)
		return user, err
	}

	user.setUser(postUser)
	err := user.calcAge(postUser.BirthDay)
	if err != nil {
		return user, err
	}
	uInformation.setUserInformation(postUser)

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return user, err
	}

	if err := db.Create(&uInformation).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return user, err
	}

	return user, nil
}
