package service

import (
	"../structs"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"time"
)

/*
 * レシーバーで受け取ったUserの異性を返す
 * 1 - 男性
 * 2 - 女性
 */
func (user *User) opponentSex() int {
	switch user.Sex {
	case 1:
		return 2
	case 2:
		return 1
	default:
		return 1
	}
}

/*
 * 引数にとったTimeから年齢を算出
 */
func (user *User) calcAge(birthDay time.Time) error {
	dateFormatOnlyNumber := "20060102"

	now := time.Now().Format(dateFormatOnlyNumber)
	birthday := birthDay.Format(dateFormatOnlyNumber)

	// 日付文字列をそのまま数値化
	nowInt, err := strconv.Atoi(now)
	if err != nil {
		return err
	}
	birthdayInt, err := strconv.Atoi(birthday)
	if err != nil {
		return err
	}

	age := (nowInt - birthdayInt) / 10000
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
	user.UserInformation = structs.UserInformation{UID: postUser.UID,
		OpponentAgeLow:    postUser.OpponentAgeLow,
		OpponentAgeUpper:  postUser.OpponentAgeUpper,
		OpponentResidence: postUser.OpponentResidence,
		Hobby:             postUser.Hobby,
		Residence:         postUser.Residence,
		Job:               postUser.Job, Personality: postUser.Personality}
}

/*
 * UserCombinationの詰め替え
 */
func (uCombi *UserCombination) setUserCombination(uid, opponentUid string) {
	uCombi.UID = uid
	uCombi.OpponentUID = opponentUid
}

/*
 * InfoCompatibleの詰め替え
 */
func (infoCompatible *InfoCompatible) setInfoCompatible(infoID, otherID uint) {
	infoCompatible.InfoID = infoID
	infoCompatible.OtherID = otherID
}

/*
 * PostUserのvalidationチェック
 * エラーがあった場合はError構造体を返す
 */
func (postUser *PostUser) checkValidate() error {
	validate := validator.New()
	var errMsges []string

	if err := validate.Struct(postUser); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errMsg string
			fieldName := err.Field()

			switch fieldName {
			case "NickName":
				var tag = err.Tag()
				switch tag {
				case "required":
					errMsg = "NickName is required"
				case "gte", "lt":
					errMsg = "NickName is between 1 and 10 characters"
				}
			case "Sex":
				errMsg = "Sex is required"
			case "BirthDay":
				errMsg = "BirthDay is required"
			case "OpponentAgeLow":
				errMsg = "OpponentAgeLow is required"
			case "OpponentAgeUpper":
				errMsg = "OpponentAgeUpper is required"
			case "OpponentResidence":
				errMsg = "OpponentResidence is required"
			case "Hobby":
				var tag = err.Tag()
				switch tag {
				case "required":
					errMsg = "Hobby is required"
				case "gte", "lt":
					errMsg = "Hobby is between 1 and 10 characters"
				}
			case "Residence":
				errMsg = "Residence is required"
			case "Job":
				errMsg = "Job is required"
			case "Personality":
				errMsg = "Personality is required"
			}
			errMsges = append(errMsges, errMsg)
		}
		return RaiseError(400, "validation failed", errMsges)
	}

	return nil
}
