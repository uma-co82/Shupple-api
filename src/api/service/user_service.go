package service

import (
	//"../db"
	//"../structs"
	"github.com/gin-gonic/gin"
	"github.com/uma-co82/Shupple-api/src/api/db"
	"github.com/uma-co82/Shupple-api/src/api/structs"
	"math/rand"
	"time"
)

type (
	UserService     struct{}
	User            structs.User
	UserInformation structs.UserInformation
	UserCombination structs.UserCombination
	InfoCompatible  structs.InfoCompatible
	Error           structs.Error
	PostUser        structs.PostUser
	PutUser         structs.PutUser
	IsRegistered    structs.IsRegistered
	IsMatched       structs.IsMatched
)

/**
 * 引数の[]Userからランダムに1件取得
 */
func getRandUser(u []User) User {
	var user User
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(u))
	user = u[i]
	return user
}

/**
 * UIDからユーザーが登録済みかを判定する
 * TODO: RecordNotFound以外のエラーハンドリング
 */
func (s UserService) IsRegisterdUser(c *gin.Context) (IsRegistered, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		user         User
		isRegistered IsRegistered
	)

	uid := c.Request.Header.Get("Uid")

	if tx.First(&user, "uid=?", uid).RecordNotFound() {
		tx.Rollback()
		isRegistered.IsRegistered = false
		return isRegistered, nil
	}

	isRegistered.IsRegistered = true

	return isRegistered, tx.Commit().Error
}

/**
 * UIDからマッチング済みかを判定する
 * マッチング済みの場合、マッチング相手を返す
 */
func (s UserService) IsMatchedUser(c *gin.Context) (IsMatched, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		user      User
		opponent  User
		isMatched IsMatched
	)

	uid := c.Request.Header.Get("uid")

	if err := tx.First(&user, "uid=?", uid).Error; err != nil {
		tx.Rollback()
		return isMatched, RaiseDBError()
	}

	if user.IsCombination == true {
		if err := tx.First(&opponent, "uid=?", user.OpponentUid).Error; err != nil {
			tx.Rollback()
			return isMatched, RaiseDBError()
		}
		if err := tx.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
			tx.Rollback()
			return isMatched, RaiseDBError()
		}
		isMatched.IsMatched = true
		tmp := structs.User(opponent)
		if err := tx.Where("uid IN (?, ?) AND opponent_uid IN (?, ?)", uid, opponent.UID, uid, opponent.UID).First(&tmp.UserCombination).Error; err != nil {
			tx.Rollback()
			return isMatched, RaiseDBError()
		}
		isMatched.User = &tmp
		return isMatched, nil
	}

	isMatched.IsMatched = false
	isMatched.User = nil
	return isMatched, tx.Commit().Error
}

/**
 * 異性かつ希望の年齢層のUserをランダムに1件返す
 * マッチング済みの場合はマッチング相手を返す
 */
func (s UserService) GetOpponent(c *gin.Context) (User, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		candidateUsers []User
		user           User
		opponent       User
		uComb          UserCombination
	)

	uid := c.Request.Header.Get("Uid")

	if err := tx.First(&user, "uid=?", uid).Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}

	// Userが既にマッチング済みの場合、userCombinationを含めて返す(フロントで時間が必要な為)
	if user.IsCombination == true {
		if err := tx.First(&opponent, "uid=?", user.OpponentUid).Error; err != nil {
			tx.Rollback()
			return opponent, RaiseDBError()
		}
		if err := tx.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
			tx.Rollback()
			return opponent, RaiseDBError()
		}
		if err := tx.Where("uid IN (?, ?) AND opponent_uid IN (?, ?)", uid, opponent.UID, uid, opponent.UID).First(&opponent.UserCombination).Error; err != nil {
			tx.Rollback()
			return opponent, RaiseDBError()
		}
		return opponent, nil
	}

	if err := tx.Model(&user).Related(&user.UserInformation, "UserInformation").Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}

	opponentSex := user.opponentSex()

	// 条件に合うユーザを検索
	// 条件にあうかつ、UserCombinationのOpponentUIDにないと言う条件で絞る
	// select * from users where age BETWEEN 20 AND 30 AND sex=1 AND is_combination=false AND uid NOT IN (select opponent_uid from user_combinations where uid='自分のuid')
	if err := tx.Where("age BETWEEN ? AND ? AND sex=? AND is_combination=? AND uid NOT IN (select opponent_uid from user_combinations where uid=?) AND uid IN (select uid from user_informations where residence=?)", user.UserInformation.OpponentAgeLow, user.UserInformation.OpponentAgeUpper, opponentSex, false, uid, user.UserInformation.OpponentResidence).Find(&candidateUsers).Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}

	if len(candidateUsers) == 0 {
		return opponent, RaiseError(404, "Opponent Not Found", nil)
	}

	opponent = getRandUser(candidateUsers)
	opponentAfter := opponent
	opponentAfter.IsCombination = true
	opponentAfter.OpponentUid = user.UID
	userAfter := user
	userAfter.OpponentUid = opponent.UID
	userAfter.IsCombination = true

	if err := tx.Model(&opponent).Update(&opponentAfter).Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}
	if err := tx.Model(&user).Update(&userAfter).Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}

	if err := tx.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}

	uComb.setUserCombination(user.UID, opponent.UID)
	if err := tx.Create(&uComb).Error; err != nil {
		tx.Rollback()
		return opponent, RaiseDBError()
	}

	opponent.UserCombination = structs.UserCombination(uComb)

	return opponent, tx.Commit().Error
}

/**
 * マッチング後48時間経過していた場合、マッチング解除
 */
func (s UserService) CancelOpponent(c *gin.Context) (bool, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		user         User
		opponent     User
		updateTarget = map[string]interface{}{"is_combination": false, "opponent_uid": nil}
	)

	uid := c.Request.Header.Get("Uid")

	if err := tx.First(&user, "uid=?", uid).Error; err != nil {
		tx.Rollback()
		return false, RaiseDBError()
	}
	if err := tx.First(&opponent, "uid=?", user.OpponentUid).Error; err != nil {
		tx.Rollback()
		return false, RaiseDBError()
	}

	if err := tx.Model(&user).Updates(updateTarget).Error; err != nil {
		tx.Rollback()
		return false, RaiseDBError()
	}
	if err := tx.Model(&opponent).Updates(updateTarget).Error; err != nil {
		tx.Rollback()
		return false, RaiseDBError()
	}

	return true, tx.Commit().Error
}

/**
 * POSTされたjsonを元にUser, UserInformation, UserCombinationを作成
 */
func (s UserService) CreateUser(c *gin.Context) (User, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		postUser  PostUser
		user      User
		s3Service S3Service
	)

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
		return user, err
	}

	if err := s3Service.UploadToS3(postUser.Image, postUser.UID); err != nil {
		return user, err
	}

	if err := postUser.checkPostUserValidate(); err != nil {
		return user, err
	}

	user.setUserFromPost(postUser)
	user.ImageURL = postUser.UID + ".png"
	err := user.calcAge(postUser.BirthDay)
	if err != nil {
		return user, err
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return user, RaiseDBError()
	}

	return user, tx.Commit().Error
}

/**
 * UIDでユーザーを検索する
 */
func (s UserService) GetUser(c *gin.Context) (User, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		user User
	)

	uid := c.Request.Header.Get("Uid")

	if err := tx.First(&user, "uid=?", uid).Error; err != nil {
		tx.Rollback()
		return user, RaiseDBError()
	}

	if err := tx.Model(&user).Related(&user.UserInformation, "UserInformation").Error; err != nil {
		tx.Rollback()
		return user, RaiseDBError()
	}

	return user, tx.Commit().Error
}

/**
 * User情報の更新
 * TODO: 飛んできたプロパティーだけ更新したい。。
 * TODO: userBefore要るのか？？
 */
func (s UserService) UpdateUser(c *gin.Context) (User, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		putUser    PutUser
		userBefore User
		userAfter  User
		s3Service  S3Service
	)

	uid := c.Request.Header.Get("Uid")

	if err := c.BindJSON(&putUser); err != nil {
		return userAfter, err
	}
	if err := putUser.checkPutUserValidate(); err != nil {
		return userAfter, err
	}

	if err := tx.First(&userAfter, "uid=?", uid).Error; err != nil {
		tx.Rollback()
		return userAfter, RaiseDBError()
	}

	if err := s3Service.UploadToS3(putUser.Image, uid); err != nil {
		return userAfter, err
	}

	userAfter.setUserFromPut(putUser)
	if err := tx.Model(&userBefore).Update(&userAfter).Error; err != nil {
		tx.Rollback()
		return userAfter, RaiseDBError()
	}

	return userBefore, tx.Commit().Error
}

func (s UserService) SoftDeleteUser(c *gin.Context) error {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		user User
	)

	uid := c.Request.Header.Get("Uid")

	if err := tx.First(&user, "uid=?", uid).Error; err != nil {
		tx.Rollback()
		return RaiseDBError()
	}

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return RaiseDBError()
	}

	return tx.Commit().Error
}

/**
 * n通以上メッセージのやり取りがあった場合に相性が良い組み合わせと考え
 * UserCompatibleに保存する
 */
func (s UserService) CreateCompatible(c *gin.Context) (InfoCompatible, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		infoCompatible InfoCompatible
		uComb          UserCombination
		uInfo          UserInformation
		otherUinfo     UserInformation
	)

	if err := c.BindJSON(&uComb); err != nil {
		return infoCompatible, err
	}

	if err := tx.First(&uInfo, "uid=?", uComb.UID).Error; err != nil {
		tx.Rollback()
		return infoCompatible, RaiseDBError()
	}
	if err := tx.First(&otherUinfo, "uid=?", uComb.OpponentUID).Error; err != nil {
		tx.Rollback()
		return infoCompatible, RaiseDBError()
	}

	infoCompatible.setInfoCompatible(uInfo.ID, otherUinfo.ID)

	if err := tx.Create(&infoCompatible).Error; err != nil {
		tx.Rollback()
		return infoCompatible, RaiseDBError()
	}
	return infoCompatible, tx.Commit().Error
}
