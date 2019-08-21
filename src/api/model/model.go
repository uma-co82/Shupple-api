package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DBとのやり取りを担うUser構造体
type User struct {
	gorm.Model
	// Firebase UID
	UID         string    `json:"uid"`
	NickName    string    `json:"nickName"`
	Sex         int       `json:"sex"`
	BirthDay    time.Time `json:"birthDay"`
	Age         int       `json:"Age"`
	OpponentAge int       `json:"opponentAge"`
	ImageURL    string    `json:"imageUrl"`
}

// DBとのやり取りを担うUserInformation構造体
type UserInformation struct {
	gorm.Model
	UID         string `json:"uid"`
	Hobby       string `json:"hobby"`
	Residence   int    `json:"residence"`
	Job         int    `json:"job"`
	Personality int    `json:"personality"`
}

// 1度マッチングしたか判定するための構造体
type UserCombination struct {
	gorm.Model
	UID         string `json:"uid"`
	OpponentUID string `json:"opponentUid"`
}

// POSTされた値を受け取る為の構造体
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

// アプリ側へ返却するUserのProfile構造体
type Profile struct {
	User        User            `json:"user"`
	Information UserInformation `json:"userInformation"`
}

type Error struct {
	Code              int      `json:"code"`
	Message           string   `json:"message"`
	ValidationMessage []string `json:"validationMessage"`
}
