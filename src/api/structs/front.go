package structs

import "time"

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
 * エラーが発生した場合にフロントへ返却するError構造体
 */
type Error struct {
	Code              int      `json:"code"`
	Message           string   `json:"message"`
	ValidationMessage []string `json:"validationMessage"`
}
