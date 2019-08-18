package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DBとのやり取りを担うUser構造体
type User struct {
	gorm.Model
	UID           string    `json:"uid"`
	NickName      string    `json:"nickName"`
	Sex           int       `json:"sex"`
	BirthDay      time.Time `json:"birthDay"`
	Age           int       `json:"Age"`
	OpponentAge   int       `json:"opponentAge"`
	ImageURL      string    `json:"imageUrl"`
	InformationID int       `json:"informationID"`
}

// DBとのやり取りを担うUserInformation構造体
type UserInformation struct {
	gorm.Model
	Hobby       string `json:"hobby"`
	Residence   int    `json:"residence"`
	Job         int    `json:"job"`
	Personality int    `json:"personality"`
}

// 1度マッチングしたか判定するための構造体
type UserCombination struct {
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
	UID         string          `json:"uid"`
	NickName    string          `json:"nickName"`
	Sex         int             `json:"sex"`
	Age         int             `json:"Age"`
	OpponentAge int             `json:"opponentAge"`
	Information UserInformation `json:"userInformation"`
}
