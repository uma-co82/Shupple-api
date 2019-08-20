package service

import (
	"fmt"
	"strconv"
	"time"
)

// レシーバーで受け取ったUserの異性を返す
// 0 - 男性
// 1 - 女性
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
