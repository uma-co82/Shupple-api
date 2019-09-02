package service

import (
	"../db"
	"../structs"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type (
	UserService     struct{}
	User            structs.User
	UserInformation structs.UserInformation
	UserCombination structs.UserCombination
	InfoCompatible  structs.InfoCompatible
	Error           structs.Error
	PostUser        structs.PostUser
)

/*
 * 引数の[]Userからランダムに1件取得
 */
func getRandUser(u []User) User {
	var user User
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(u))
	user = u[i]
	return user
}

/*
 * 異性かつ希望の年齢層のUserをランダムに1件返す
 * TODO: マッチング状態かどうかの判定カラムで検索条件追加
 */
func (s UserService) GetOpponent(c *gin.Context) (User, error) {
	db := db.GetDB()
	var (
		candidateUsers []User
		user           User
		opponent       User
		uComb          UserCombination
	)

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return opponent, err
	}
	if err := db.Model(&user).Related(&user.UserInformation, "UserInformation").Error; err != nil {
		return opponent, err
	}

	opponentSex := user.opponentSex()

	// 条件に合うユーザを検索
	// 条件にあうかつ、UserCombinationのOpponentUIDにないと言う条件で絞る
	// select * from users where age BETWEEN 20 AND 30 AND sex=1 AND is_combination=false AND uid NOT IN (select opponent_uid from user_combinations where uid='自分のuid')
	if err := db.Where("age BETWEEN ? AND ? AND sex=? AND is_combination=? AND uid NOT IN (select opponent_uid from user_combinations where uid=?)", user.UserInformation.OpponentAgeLow, user.UserInformation.OpponentAgeUpper, opponentSex, false, uid).Find(&candidateUsers).Error; err != nil {
		return opponent, err
	}

	if len(candidateUsers) == 0 {
	}

	opponent = getRandUser(candidateUsers)
	opponentAfter := opponent
	opponentAfter.IsCombination = true
	userAfter := user
	userAfter.IsCombination = true

	if err := db.Model(&opponent).Update(&opponentAfter).Error; err != nil {
		return opponent, err
	}
	if err := db.Model(&user).Update(&userAfter).Error; err != nil {
		return opponent, err
	}

	if err := db.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
		return opponent, err
	}

	uComb.setUserCombination(user.UID, opponent.UID)
	if err := db.Create(&uComb).Error; err != nil {
		return opponent, err
	}

	return opponent, nil
}

/*
 * POSTされたjsonを元にUser, UserInformation, UserCombinationを作成
 */
func (s UserService) CreateUser(c *gin.Context) (User, error) {
	db := db.GetDB()
	var (
		postUser PostUser
		user     User
	)

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
		fmt.Printf("Binding Error %v", err)
		return user, err
	}

	if err := postUser.checkValidate(); err != nil {
		return user, err
	}

	user.setUser(postUser)
	err := user.calcAge(postUser.BirthDay)
	if err != nil {
		return user, err
	}

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return user, err
	}

	return user, nil
}

/*
 * UIDでユーザーを検索する
 */
func (s UserService) GetUser(c *gin.Context) (User, error) {
	db := db.GetDB()
	var (
		user User
	)

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return user, err
	}

	if err := db.Model(&user).Related(&user.UserInformation, "UserInformation").Error; err != nil {
		return user, err
	}

	return user, nil
}

/*
 * User情報の更新
 */
func (s UserService) Update(c *gin.Context) (User, error) {
	db := db.GetDB()
	var (
		postUser   PostUser
		userBefore User
		userAfter  User
	)

	uid := c.Request.Header.Get("Uid")

	if err := c.BindJSON(&postUser); err != nil {
		return userAfter, err
	}
	if err := postUser.checkValidate(); err != nil {
		return userAfter, err
	}

	userBefore.UID = uid

	if err := db.First(&userAfter, "uid=?", uid).Error; err != nil {
		return userAfter, err
	}

	if err := db.Model(&userAfter).Related(&userAfter.UserInformation, "UserInformation").Error; err != nil {
		return userAfter, err
	}

	userAfter.setUser(postUser)
	if err := userAfter.calcAge(postUser.BirthDay); err != nil {
		return userAfter, err
	}

	if err := db.Model(&userBefore).Update(&userAfter).Error; err != nil {
		return userAfter, err
	}

	return userAfter, nil
}

/*
 * n通以上メッセージのやり取りがあった場合に相性が良い組み合わせと考え
 * UserCompatibleに保存する
 * MEMO: そもそもこれフロントからinfoID送られないと思うので一旦放置
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
