package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

/**
 * DBとのやり取りを担うUser構造体
 */
type User struct {
	gorm.Model
	UID              string            `json:"uid"`
	NickName         string            `json:"nickName"`
	Sex              int               `json:"sex"`
	BirthDay         time.Time         `json:"birthDay"`
	Age              int               `json:"age"`
	ImageURL         string            `json:"imageUrl"`
	UserInformation  UserInformation   `gorm:"forignkey:UID"`
	UserCombinations []UserCombination `gorm:"foreignkey:UID;association_foreignkey:UID"`
}

/*
 * DBとのやり取りを担うUserInformation構造体
 */
type UserInformation struct {
	gorm.Model
	UID              string `json:"uid"`
	OpponentAgeLow   int    `json:"opponentAgeLow"`
	OpponentAgeUpper int    `json:"opponentAgeUpper"`
	Hobby            string `json:"hobby"`
	Residence        int    `json:"residence"`
	Job              int    `json:"job"`
	Personality      int    `json:"personality"`
}

/**
 * 任意のUIDの組み合わせを表す構造体
 */
type UserCombination struct {
	gorm.Model
	UID         string `json:"uid"`
	OpponentUID string `json:"otherID"`
}

/**
 * 相性が良いUserInformationを記録していくUserCompatible構造体
 */
type InfoCompatible struct {
	gorm.Model
	InfoID  string `json:"infoID"`
	OtherID string `json:"otherID"`
}

/**
 * POSTされた値を受け取る為の構造体
 */
type PostUser struct {
	UID              string    `json:"uid"`
	NickName         string    `json:"nickName"`
	Sex              int       `json:"sex"`
	BirthDay         time.Time `json:"birthDay"`
	OpponentAgeLow   int       `json:"opponentAgeLow"`
	OpponentAgeUpper int       `json:"opponentAgeUpper"`
	Hobby            string    `json:"hobby"`
	Residence        int       `json:"residence"`
	Job              int       `json:"job"`
	Personality      int       `json:"personality"`
}

/**
 * フロントへ返却するUserのProfile構造体
 */
type Profile struct {
	User        User            `json:"user"`
	Information UserInformation `json:"userInformation"`
}

/**
 * エラーが発生した場合にフロントへ返却するError構造体
 */
type Error struct {
	Code              int      `json:"code"`
	Message           string   `json:"message"`
	ValidationMessage []string `json:"validationMessage"`
}
