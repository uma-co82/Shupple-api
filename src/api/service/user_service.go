package service

import (
	"../db"
	"../structs"
	"fmt"
	"github.com/gin-gonic/gin"
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
	defer db.Close()
	var (
		user         User
		isRegistered IsRegistered
	)

	uid := c.Request.Header.Get("Uid")

	if db.First(&user, "uid=?", uid).RecordNotFound() {
		isRegistered.IsRegistered = false
		return isRegistered, nil
	}

	isRegistered.IsRegistered = true

	return isRegistered, nil
}

/**
 * UIDからマッチング済みかを判定する
 * マッチング済みの場合、マッチング相手を返す
 */
func (s UserService) IsMatchedUser(c *gin.Context) (IsMatched, error) {
	db := db.Init()
	defer db.Close()
	var (
		user      User
		opponent  User
		isMatched IsMatched
	)

	uid := c.Request.Header.Get("uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return isMatched, err
	}

	if user.IsCombination == true {
		if err := db.First(&opponent, "uid=?", user.OpponentUid).Error; err != nil {
			return isMatched, err
		}
		if err := db.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
			return isMatched, err
		}
		isMatched.IsMatched = true
		tmp := structs.User(opponent)
		if err := db.Where("uid IN (?, ?) AND opponent_uid IN (?, ?)", uid, opponent.UID, uid, opponent.UID).First(&tmp.UserCombination).Error; err != nil {
			return isMatched, RaiseDBError()
		}
		isMatched.User = &tmp
		return isMatched, nil
	}

	isMatched.IsMatched = false
	isMatched.User = nil
	return isMatched, nil
}

/*
*
 * 異性かつ希望の年齢層のUserをランダムに1件返す
 * マッチング済みの場合はマッチング相手を返す
*/
func (s UserService) GetOpponent(c *gin.Context) (User, error) {
	db := db.Init()
	defer db.Close()
	var (
		candidateUsers []User
		user           User
		opponent       User
		uComb          UserCombination
	)

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return opponent, err
	}

	// Userが既にマッチング済みの場合、userCombinationを含めて返す(フロントで時間が必要な為)
	if user.IsCombination == true {
		if err := db.First(&opponent, "uid=?", user.OpponentUid).Error; err != nil {
			return opponent, err
		}
		if err := db.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
			return opponent, err
		}
		if err := db.Where("uid IN (?, ?) AND opponent_uid IN (?, ?)", uid, opponent.UID, uid, opponent.UID).First(&opponent.UserCombination).Error; err != nil {
			return opponent, RaiseDBError()
		}
		return opponent, nil
	}

	if err := db.Model(&user).Related(&user.UserInformation, "UserInformation").Error; err != nil {
		return opponent, err
	}

	opponentSex := user.opponentSex()

	// 条件に合うユーザを検索
	// 条件にあうかつ、UserCombinationのOpponentUIDにないと言う条件で絞る
	// select * from users where age BETWEEN 20 AND 30 AND sex=1 AND is_combination=false AND uid NOT IN (select opponent_uid from user_combinations where uid='自分のuid')
	if err := db.Where("age BETWEEN ? AND ? AND sex=? AND is_combination=? AND uid NOT IN (select opponent_uid from user_combinations where uid=?) AND uid IN (select uid from user_informations where residence=?)", user.UserInformation.OpponentAgeLow, user.UserInformation.OpponentAgeUpper, opponentSex, false, uid, user.UserInformation.OpponentResidence).Find(&candidateUsers).Error; err != nil {
		return opponent, err
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

	if err := db.Model(&opponent).Update(&opponentAfter).Error; err != nil {
		return opponent, err
	}
	if err := db.Model(&user).Update(&userAfter).Error; err != nil {
		return opponent, err
	}

	if err := db.Model(&opponent).Related(&opponent.UserInformation, "UserInformation").Error; err != nil {
		return opponent, err
	}

	uComb.setUserCombination(user.UID, opponent.UID)
	if err := db.Create(&uComb).Error; err != nil {
		return opponent, err
	}

	return opponent, nil
}

/**
 * マッチング後48時間経過していた場合、マッチング解除
 */
func (s UserService) CancelOpponent(c *gin.Context) (bool, error) {
	db := db.Init()
	defer db.Close()
	var (
		user         User
		opponent     User
		updateTarget = map[string]interface{}{"is_combinaton": false, "opponent_uid": nil}
	)

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return false, RaiseDBError()
	}
	if err := db.First(&opponent, "uid=?", user.OpponentUid).Error; err != nil {
		return false, RaiseDBError()
	}

	if err := db.Model(&user).Updates(updateTarget).Error; err != nil {
		return false, RaiseDBError()
	}
	if err := db.Model(&opponent).Updates(updateTarget).Error; err != nil {
		return false, RaiseDBError()
	}

	return true, nil
}

/**
 * POSTされたjsonを元にUser, UserInformation, UserCombinationを作成
 */
func (s UserService) CreateUser(c *gin.Context) (User, error) {
	db := db.Init()
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

	if err := s3Service.UploadToS3(postUser.Image, postUser.NickName); err != nil {
		return user, err
	}

	if err := postUser.checkPostUserValidate(); err != nil {
		return user, err
	}

	user.setUserFromPost(postUser)
	user.ImageURL = postUser.NickName + ".png"
	err := user.calcAge(postUser.BirthDay)
	if err != nil {
		return user, err
	}

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("DB Error %v", err)
		return user, err
	}

	return user, nil
}

/**
 * UIDでユーザーを検索する
 */
func (s UserService) GetUser(c *gin.Context) (User, error) {
	db := db.Init()
	defer db.Close()
	var (
		user User
	)

	uid := c.Request.Header.Get("Uid")

	if err := db.First(&user, "uid=?", uid).Error; err != nil {
		return user, err
	}

	if err := db.Model(&user).Related(&user.UserInformation, "UserInformation").Error; err != nil {
		return user, err
	}

	return user, nil
}

/**
 * User情報の更新
 * TODO: 飛んできたプロパティーだけ更新したい。。
 * TODO: userBefore要るのか？？
 */
func (s UserService) UpdateUser(c *gin.Context) (User, error) {
	db := db.Init()
	defer db.Close()
	var (
		putUser    PutUser
		userBefore User
		userAfter  User
	)

	uid := c.Request.Header.Get("Uid")

	if err := c.BindJSON(&putUser); err != nil {
		return userAfter, err
	}
	if err := putUser.checkPutUserValidate(); err != nil {
		return userAfter, err
	}

	if err := db.First(&userAfter, "uid=?", uid).Error; err != nil {
		return userAfter, err
	}

	userAfter.setUserFromPut(putUser)
	if err := db.Model(&userBefore).Update(&userAfter).Error; err != nil {
		return userAfter, err
	}

	return userBefore, nil
}

/**
 * n通以上メッセージのやり取りがあった場合に相性が良い組み合わせと考え
 * UserCompatibleに保存する
 */
func (s UserService) CreateCompatible(c *gin.Context) (InfoCompatible, error) {
	db := db.Init()
	db.Close()
	var (
		infoCompatible InfoCompatible
		uComb          UserCombination
		uInfo          UserInformation
		otherUinfo     UserInformation
	)

	if err := c.BindJSON(&uComb); err != nil {
		return infoCompatible, err
	}

	if err := db.First(&uInfo, "uid=?", uComb.UID).Error; err != nil {
		return infoCompatible, err
	}
	if err := db.First(&otherUinfo, "uid=?", uComb.OpponentUID).Error; err != nil {
		return infoCompatible, err
	}

	//infoCompatible.setInfoCompatible(uInfo.ID, otherUinfo.ID)

	if err := db.Create(&infoCompatible).Error; err != nil {
		return infoCompatible, err
	}

	return infoCompatible, nil
}
