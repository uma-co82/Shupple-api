package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
 * DBとのやり取りを担うUser構造体
 */
type User struct {
	gorm.Model
	UID      string    `json:"uid"`
	NickName string    `json:"nickName"`
	Sex      int       `json:"sex"`
	BirthDay time.Time `json:"birthDay"`
	Age      int       `json:"age"`
	ImageURL string    `json:"imageUrl"`
}

/*
 * DBとのやり取りを担うUserInformation構造体
 */
type UserInformation struct {
	gorm.Model
	UID         string `json:"uid"`
	OpponentAge int    `json:"opponentAge"`
	Hobby       string `json:"hobby"`
	Residence   int    `json:"residence"`
	Job         int    `json:"job"`
	Personality int    `json:"personality"`
}

/*
 * 1度マッチングしたか判定するための構造体
 */
type UserCombination struct {
	gorm.Model
	UID         string `json:"uid"`
	OpponentUID string `json:"opponentUid"`
}

/*
 * 相性が良いUserInformationを記録していくUserCompatible構造体
 */
type UserCompatible struct {
	gorm.Model
	InfoID  string `json:"infoID"`
	OtherID string `json:"otherID"`
}

/*
 * POSTされた値を受け取る為の構造体
 */
type PostUser struct {
	UID         string    `json:"uid"`
	NickName    string    `json:"nickName"`
	Sex         int       `json:"sex"`
	BirthDay    time.Time `json:"birthDay"`
	OpponentAge int       `json:"opponentAge"`
	Hobby       string    `json:"hobby"`
	Residence   int       `json:"residence"`
	Job         int       `json:"job"`
	Personality int       `json:"personality"`
}

/*
 * フロントへ返却するUserのProfile構造体
 */
type Profile struct {
	User        User            `json:"user"`
	Information UserInformation `json:"userInformation"`
}

/*
 * エラーが発生した場合にフロントへ返却するError構造体
 */
type Error struct {
	Code              int      `json:"code"`
	Message           string   `json:"message"`
	ValidationMessage []string `json:"validationMessage"`
}
