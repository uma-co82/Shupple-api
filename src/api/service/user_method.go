package service

import (
	"strconv"
	"time"
)

/*
 * レシーバーで受け取ったUserの異性を返す
 * 0 - 男性
 * 1 - 女性
 */
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

/*
 * 引数にとったTimeから年齢を算出
 */
func (user *User) calcAge(t time.Time) error {
	dateFormatOnlyNumber := "20060102"

	now := time.Now().Format(dateFormatOnlyNumber)
	birthday := t.Format(dateFormatOnlyNumber)

	// 日付文字列をそのまま数値化
	// TODO: エラハン
	nowInt, err := strconv.Atoi(now)
	if err != nil {
		return err
	}
	birthdayInt, err := strconv.Atoi(birthday)
	if err != nil {
		return err
	}

	age := (nowInt-birthdayInt)/10000 + 1
	user.Age = age
	return nil
}

/*
 * Userの詰め替え
 */
func (user *User) setUser(postUser PostUser) {
	user.UID = postUser.UID
	user.NickName = postUser.NickName
	user.Sex = postUser.Sex
	user.BirthDay = postUser.BirthDay
}

/*
 * UserInformationの詰め替え
 */
func (uInfo *UserInformation) setUserInformation(postUser PostUser) {
	uInfo.UID = postUser.UID
	uInfo.OpponentAgeLow = postUser.OpponentAgeLow
	uInfo.OpponentAgeUpper = postUser.OpponentAgeUpper
	uInfo.Hobby = postUser.Hobby
	uInfo.Residence = postUser.Residence
	uInfo.Job = postUser.Job
	uInfo.Personality = postUser.Personality
}

/*
 * UserCombinationの詰め替え
 */
func (uCombi *UserCombination) setUserCombination(uid string, opponentUid string) {
	uCombi.UID = uid
	uCombi.OpponentUID = opponentUid
}
