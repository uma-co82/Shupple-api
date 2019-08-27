package service

import (
	"../db"
	"../model"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type (
	UserService     struct{}
	PostUser        model.PostUser
	User            model.User
	UserInformation model.UserInformation
	UserCombination model.UserCombination
	InfoCompatible  model.InfoCompatible
	Profile         model.Profile
	Error           model.Error
)

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
 * 異性かつ希望の年齢層のUserをランダムに1件返す
 */
func (s UserService) GetOpponent(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var (
		users         []User
		user          User
		opponent      User
		profile       Profile
		uComb         UserCombination
		uCombinations []UserCombination
		uInfo         UserInformation
		opponentInfo  UserInformation
	)

	// 自分のUID
	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return profile, err
	}
	if err := db.First(&uInfo, "uid=?", uid).Error; err != nil {
		return profile, err
	}

	opponentSex := user.opponentSex()

	if err := db.Where("age BETWEEN ? AND ? AND sex=?", uInfo.OpponentAgeLow, uInfo.OpponentAgeUpper, opponentSex).Find(&users).Error; err != nil {
		return profile, err
	}

	// TODO: 新規のユーザーが見つからなかったら無限ループしちゃう
	for {
		var err error
		opponent, err = getRandUser(users)
		if err != nil {
			return profile, err
		}
		if err := db.Where("uid=? AND opponent_uid=?", user.UID, opponent.UID).Find(&uCombinations).Error; err != nil {
			return profile, err
		}
		if len(uCombinations) == 0 {
			break
		}
	}

	if err := db.First(&opponentInfo, "uid=?", opponent.UID).Error; err != nil {
		return profile, err
	}

	uComb.setUserCombination(user.UID, opponent.UID)
	if err := db.Create(&uComb).Error; err != nil {
		return profile, err
	}

	profile = Profile{User: model.User(opponent),
		Information: model.UserInformation(opponentInfo)}

	return profile, nil
}

/*
 * POSTされたjsonを元にUser, UserInformation, UserCombinationを作成
 */
func (s UserService) CreateUser(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var (
		postUser PostUser
		user     User
		uInfo    UserInformation
		profile  Profile
	)

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
	var (
		user         User
		uInformation UserInformation
		profile      Profile
	)

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
func (s UserService) Update(c *gin.Context) (Profile, error) {
	db := db.GetDB()
	var (
		postUser           PostUser
		userBefore         User
		userAfter          User
		uInformationBefore UserInformation
		uInformationAfter  UserInformation
		profile            Profile
	)

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

/*
 * n通以上メッセージのやり取りがあった場合に相性が良い組み合わせと考え
 * UserCompatibleに保存する
 */
func (s UserService) CreateCompatible(c *gin.Context) (InfoCompatible, error) {
	db := db.GetDB()
	var infoCompatible InfoCompatible

	if err := c.BindJSON(&infoCompatible); err != nil {
		return infoCompatible, err
	}

	if err := db.Create(&infoCompatible).Error; err != nil {
		return infoCompatible, err
	}

	return infoCompatible, nil
}
