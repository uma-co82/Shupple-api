package model

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// DBへ保存する為の構造体
type User struct {
	gorm.Model
	UID      string `json:"uid"`
	NickName string `json:"nickName"`
	Sex      string `json:"sex"`
	Hobby    string `json:"hobby"`
	Age      int    `json:"Age"`
}

// POSTされた値を受け取る為の構造体
type PostUser struct {
	NickName string    `json:"nickName"`
	Sex      string    `json:"sex"`
	Hobby    string    `json:"hobby"`
	BirthDay time.Time `json:"birthDay"`
}

// UID作成
func GetUUID() string {
	u, _ := uuid.NewRandom()
	uu := u.String()
	return uu
}

// 引数にとったTimeから年齢を算出
func CalcAge(t time.Time) (int, error) {
	dateFormatOnlyNumber := "20060102" // YYYYMMDD

	now := time.Now().Format(dateFormatOnlyNumber)
	birthday := t.Format(dateFormatOnlyNumber)

	// 日付文字列をそのまま数値化
	nowInt, err := strconv.Atoi(now)
	if err != nil {
		return 0, err
	}
	birthdayInt, err := strconv.Atoi(birthday)
	if err != nil {
		return 0, err
	}

	age := (nowInt - birthdayInt) / 10000
	return age, nil
}
