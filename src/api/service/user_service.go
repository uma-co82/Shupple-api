package service

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"../db"
	"../model"
	"github.com/gin-gonic/gin"
)

type UserService struct{}

type PostUser model.PostUser
type User model.User

type UID struct {
	uid string `json:"uid"`
}

func (user *User) opponentSex() int {
	switch user.Sex {
	case 0:
		return 1
	case 1:
		return 0
	default:
		return 0
	}
}

// TODO: エラーハンドリング
// 引数にとったTimeから年齢を算出
func (user *User) calcAge(t time.Time) {
	dateFormatOnlyNumber := "20060102" // YYYYMMDD

	now := time.Now().Format(dateFormatOnlyNumber)
	birthday := t.Format(dateFormatOnlyNumber)

	// 日付文字列をそのまま数値化
	nowInt, err := strconv.Atoi(now)
	if err != nil {
		fmt.Println(err)
		return
	}
	birthdayInt, err := strconv.Atoi(birthday)
	if err != nil {
		fmt.Println(err)
		return
	}

	age := (nowInt - birthdayInt) / 10000
	user.Age = age
}

// ランダム取得
func getRandUser(u []User) User {
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

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return user, err
	}

	return user, nil
}
