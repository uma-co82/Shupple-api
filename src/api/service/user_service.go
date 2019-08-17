package service

import (
	"fmt"
	"strconv"
	"time"

	"../db"
	"../model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService struct{}

type PostUser model.PostUser
type User model.User

func (user *User) getUUID() {
	u, _ := uuid.NewRandom()
	uu := u.String()
	user.UID = uu
}

// TODO: エラーハンドリングのやり方が分からん
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
		return user, err
	}

	user.getUUID()
	user.calcAge(postUser.BirthDay)
	user.NickName = postUser.NickName
	user.Sex = postUser.Sex
	user.Hobby = postUser.Hobby

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("DBError")
		fmt.Printf("%v", err)
		return user, err
	}

	return user, nil
}
