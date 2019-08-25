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

/*
 * 引数の[]Userからランダムに1件取得
 */
func getRandUser(u []User) (User, error) {
	var user User
	if u == nil {
		err := RaiseError(404, "Opponent Not Found", nil)
		return user, err
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(u))
	user = u[i]
	return user, nil
}
/*
 * 異性のUserを返す
 */
func (s UserService) GetOpponent(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var users []User
	var user User
	var opponent User
	var profile Profile
	var uCombi UserCombination
	var uCombinations []UserCombination
	var uInfo UserInformation

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	opponentSex := user.opponentSex()

	if err := db.Find(&users, "sex=?", opponentSex).Error; err != nil {
		return profile, err
	}

	for {
		var err error
		opponent, err = getRandUser(users)
		if err != nil {
			return profile, err
		}
		if err := db.Where("uid=? AND opponent_uid=?", user.UID, opponent.UID).Find(&uCombinations).Error; err != nil {
			return profile, err
		}
		fmt.Printf("*************************** %v * \n", uCombinations)
		if len(uCombinations) == 0 {
			break
		}
	}

	if err := db.First(&uInfo, "uid=?", opponent.UID).Error; err != nil {
		return profile, err
	}

	uCombi.setUserCombination(user.UID, opponent.UID)
	if err := db.Create(&uCombi).Error; err != nil {
		return profile, err
	}

	profile = Profile{User: model.User(opponent),
		Information: model.UserInformation(uInfo)}

	return profile, nil
}

/*
 * POSTされたjsonを元にUser, UserInformation, UserCombinationを作成
 */
func (s UserService) CreateUser(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var postUser PostUser
	var user User
	var uInfo UserInformation
	var profile Profile

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
		fmt.Printf("Binding Error %v", err)
		return profile, err
	}

	user.setUser(postUser)
	err := user.calcAge(postUser.BirthDay)
	if err != nil {
		return profile, err
	}
	uInfo.setUserInformation(postUser)

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return profile, err
	}

	if err := db.Create(&uInfo).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return profile, err
	}

	profile = Profile{User: model.User(user),
		Information: model.UserInformation(uInfo)}

	return profile, nil
}
/*
 * UIDでユーザーを検索する
 */
func (s UserService) GetUser(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var user User
	var uInformation UserInformation
	var profile Profile

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	if err := db.First(&uInformation, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	profile = Profile{User: model.User(user),
		Information: model.UserInformation(uInformation)}

	return profile, nil
}
/*
 * User情報の更新
 */
func (s UserService) Update(c *gin.Context) (Profile, error)  {
	db := db.GetDB()
	var postUser PostUser
	var userBefore User
	var userAfter User
	var uInformationBefore UserInformation
	var uInformationAfter UserInformation
	var profile Profile

	uid := c.Request.Header.Get("Uid")

	if err := c.BindJSON(&postUser); err != nil {
		return profile, err
	}

	userBefore.UID = uid
	uInformationBefore.UID = uid

	if err := db.First(&userAfter, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	if err := db.First(&uInformationAfter, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	userAfter.setUser(postUser)
	if err := userAfter.calcAge(postUser.BirthDay); err != nil {
		return profile, err
	}
	uInformationAfter.setUserInformation(postUser)

	if err := db.Model(&userBefore).Update(&userAfter).Error; err != nil {
		return profile, err
	}

	if err := db.Model(&uInformationBefore).Update(&uInformationAfter).Error; err != nil {
		return profile, err
	}

	profile = Profile{User: model.User(userAfter),
		Information: model.UserInformation(uInformationAfter)}

	return profile, nil
}
