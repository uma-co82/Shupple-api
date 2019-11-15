package service

import (
	"github.com/uma-co82/Shupple-api/src/api/domain/repository"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uma-co82/Shupple-api/src/api/domain"
	"github.com/uma-co82/Shupple-api/src/api/domain/user"
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/db"
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/s3"
)

type (
	UserService struct{}
)

/**
 * 引数の[]Userからランダムに1件取得
 */
func getRandUser(u []user.User) user.User {
	var person user.User
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(u))
	person = u[i]
	return person
}

/**
 * UIDからユーザーが登録済みかを判定する
 */
func (s UserService) IsRegisteredUser(c *gin.Context) (user.IsRegistered, error) {
	db := db.Init()
	defer db.Close()
	var (
		isRegistered user.IsRegistered
	)

	uid := c.Request.Header.Get("Uid")

	// Transaction
	tx := db.Begin()
	userRepository := repository.NewUserRepository(tx)
	_, err := userRepository.GetByUid(uid)
	if err != nil {
		isRegistered.IsRegistered = false
		return isRegistered, nil
	}
	tx.Commit()
	// Transaction

	isRegistered.IsRegistered = true

	return isRegistered, nil
}

/**
 * UIDからマッチング済みかを判定する
 * マッチング済みの場合、マッチング相手を返す
 */
func (s UserService) IsMatchedUser(c *gin.Context) (user.IsMatched, error) {
	db := db.Init()
	defer db.Close()
	var (
		isMatched user.IsMatched
	)

	uid := c.Request.Header.Get("uid")

	// Transaction
	tx := db.Begin()
	userRepository := repository.NewUserRepository(tx)
	person, err := userRepository.GetByUid(uid)
	if err != nil {
		return isMatched, domain.RaiseDBError()
	}

	if person.IsCombination == true {
		opponent, err := userRepository.GetByUid(person.OpponentUid)
		if err != nil {
			return isMatched, err
		}
		// TODO: ポインタを引数にしてUserInformationを直接詰めたい。。
		// MEMO: ポインタを引数に渡すとgormに怒られる...
		opponent.UserInformation, err = userRepository.GetUserInformationByRelatedUser(opponent)
		if err != nil {
			return isMatched, err
		}
		// TODO: ここも同じ！(メソッド使えば行けるかも!)
		opponent.UserCombination, err = userRepository.GetUserCombinationByBothUid(person.UID, opponent.UID)
		if err != nil {
			return isMatched, err
		}

		isMatched.IsMatched = true
		isMatched.User = &opponent
		return isMatched, nil
	}
	tx.Commit()
	// Transaction

	isMatched.IsMatched = false
	isMatched.User = nil
	return isMatched, nil
}

/**
 * 異性かつ希望の年齢層のUserをランダムに1件返す
 * マッチング済みの場合はマッチング相手を返す
 */
func (s UserService) GetOpponent(c *gin.Context) (user.User, error) {
	db := db.Init()
	defer db.Close()
	var (
		candidateUsers []user.User
		opponent       user.User
		uComb          user.UserCombination
	)

	uid := c.Request.Header.Get("Uid")

	tx := db.Begin()
	userRepository := repository.NewUserRepository(tx)
	person, err := userRepository.GetByUid(uid)
	if err != nil {
		return opponent, domain.RaiseDBError()
	}

	// 当事者(person)が既にマッチング済みの場合、userCombinationを含めて返す(フロントで時間が必要な為)
	if person.IsCombination == true {
		opponent, err := userRepository.GetByUid(person.OpponentUid)
		if err != nil {
			return opponent, err
		}
		// TODO: ポインタを引数にしてUserInformationを直接詰めたい。。
		// MEMO: ポインタを引数に渡すとgormに怒られる...
		opponent.UserInformation, err = userRepository.GetUserInformationByRelatedUser(opponent)
		if err != nil {
			return opponent, err
		}
		// TODO: ここも同じ！(メソッド使えば行けるかも!)
		opponent.UserCombination, err = userRepository.GetUserCombinationByBothUid(person.UID, opponent.UID)
		if err != nil {
			return opponent, err
		}
		return opponent, nil
	}

	// 当事者のUserInformationを取得
	person.UserInformation, err = userRepository.GetUserInformationByRelatedUser(person)
	if err != nil {
		return opponent, err
	}

	// 相手の性別をゲット
	opponentSex := person.OpponentSex()

	// 条件に合うユーザを検索
	// 条件にあうかつ、UserCombinationのOpponentUIDにないと言う条件で絞る
	// select * from users where age BETWEEN 20 AND 30 AND sex=1 AND is_combination=false AND uid NOT IN (select opponent_uid from user_combinations where uid='自分のuid')
	candidateUsers, err = userRepository.GetShupple(
		person.UserInformation.OpponentAgeLow,
		person.UserInformation.OpponentAgeUpper,
		opponentSex,
		person.UserInformation.OpponentResidence, uid,
	)
	if err != nil {
		return opponent, err
	}

	if len(candidateUsers) == 0 {
		return opponent, domain.RaiseError(404, "Opponent Not Found", nil)
	}

	opponent = getRandUser(candidateUsers)
	opponentAfter := opponent
	opponentAfter.IsCombination = true
	opponentAfter.OpponentUid = person.UID
	userAfter := person
	userAfter.OpponentUid = opponent.UID
	userAfter.IsCombination = true

	if err := userRepository.Update(opponent, opponentAfter); err != nil {
		tx.Rollback()
		return opponent, err
	}

	if err := userRepository.Update(person, userAfter); err != nil {
		tx.Rollback()
		return opponent, err
	}

	opponent.UserInformation, err = userRepository.GetUserInformationByRelatedUser(opponent)
	if err != nil {
		return opponent, err
	}

	uComb.SetUserCombination(person.UID, opponent.UID)
	if err := userRepository.CreateUserCombination(uComb); err != nil {
		tx.Rollback()
		return opponent, err
	}

	// MEMO: CreatedAtとか時間系がレスポンスに入って無いけど良いんだっけ？？
	opponent.UserCombination = user.UserCombination(uComb)
	tx.Commit()

	return opponent, nil
}

/**
 * マッチング後48時間経過していた場合、マッチング解除
 */
func (s UserService) CancelOpponent(c *gin.Context) (bool, error) {
	db := db.Init()
	defer db.Close()

	uid := c.Request.Header.Get("Uid")

	// Transaction
	tx := db.Begin()
	userRepository := repository.NewUserRepository(tx)

	person, err := userRepository.GetByUid(uid)
	if err != nil {
		return false, err
	}

	opponent, err := userRepository.GetByUid(person.OpponentUid)
	if err != nil {
		return false, err
	}

	if err := userRepository.CancelMatchingStatus(person); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := userRepository.CancelMatchingStatus(opponent); err != nil {
		tx.Rollback()
		return false, err
	}
	tx.Commit()
	// Transaction

	return true, nil
}

/**
 * POSTされたjsonを元にUser, UserInformation, UserCombinationを作成
 * TODO: チャネルがinterface{}で微妙。。
 */
func (s UserService) CreateUser(c *gin.Context) (user.User, error) {
	db := db.Init()
	defer db.Close()

	var (
		postUser  user.PostUser
		person    user.User
		s3Service s3.S3Service
		errChan   = make(chan interface{}, 1)
	)

	// TODO: Bind出来なかった時のエラーハンドリング
	if err := c.BindJSON(&postUser); err != nil {
		return person, err
	}

	// 並行処理でs3へアップロード
	// エラーがあった場合はerrChanに流す
	if postUser.Image != "" {
		person.ImageURL = postUser.UID + ".png"
		go func() {
			if err := s3Service.UploadToS3(postUser.Image, postUser.UID); err != nil {
				// エラーが起きた場合にチャネルに流す
				errChan <- err
			} else {
				errChan <- "noError"
			}
		}()
	} else {
		errChan <- "noError"
	}

	// Transaction
	tx := db.Begin()
	userRepository := repository.NewUserRepository(tx)

	if err := postUser.CheckPostUserValidate(); err != nil {
		return person, err
	}

	person.SetUserFromPost(postUser)

	err := person.CalcAge(postUser.BirthDay)
	if err != nil {
		return person, err
	}

	if err := userRepository.CreateUser(person); err != nil {
		return person, err
	}
	tx.Commit()
	// Transaction

	if err := <-errChan; err != "noError" {
		return person, err.(error)
	}
	defer close(errChan)

	return person, nil
}

/**
 * UIDでユーザーを検索する
 */
func (s UserService) GetUser(c *gin.Context) (user.User, error) {
	db := db.Init()
	defer db.Close()

	uid := c.Request.Header.Get("Uid")

	userRepository := repository.NewUserRepository(db)
	person, err := userRepository.GetByUid(uid)
	if err != nil {
		return person, err
	}

	person.UserInformation, err = userRepository.GetUserInformationByRelatedUser(person)

	return person, nil
}

/**
 * User情報の更新
 * TODO: 飛んできたプロパティーだけ更新したい。。
 */
func (s UserService) UpdateUser(c *gin.Context) (user.User, error) {
	db := db.Init()
	defer db.Close()
	var (
		putUser     user.PutUser
		afterPerson user.User
		s3Service   s3.S3Service
		errChan     = make(chan interface{}, 1)
	)

	uid := c.Request.Header.Get("Uid")

	if err := c.BindJSON(&putUser); err != nil {
		return afterPerson, err
	}

	if putUser.Image != "" {
		go func() {
			if err := s3Service.UploadToS3(putUser.Image, uid); err != nil {
				errChan <- err
			} else {
				errChan <- "noError"
			}
		}()
	} else {
		errChan <- "noError"
	}

	// Transaction
	tx := db.Begin()
	userRepository := repository.NewUserRepository(tx)

	if err := putUser.CheckPutUserValidate(); err != nil {
		return afterPerson, err
	}

	afterPerson.SetUserFromPut(putUser)

	beforePerson, err := userRepository.GetByUid(uid)
	if err != nil {
		return afterPerson, err
	}

	if err := userRepository.Update(beforePerson, afterPerson); err != nil {
		return afterPerson, err
	}

	if err := <-errChan; err != "noError" {
		return afterPerson, err.(error)
	}
	defer close(errChan)

	tx.Commit()
	// Transaction

	return beforePerson, nil
}

func (s UserService) SoftDeleteUser(c *gin.Context) error {
	db := db.Init()
	defer db.Close()

	uid := c.Request.Header.Get("Uid")

	// Transaction
	tx := db.Begin()
	userRepository := repository.NewUserRepository(db)
	person, err := userRepository.GetByUid(uid)
	if err != nil {
		return err
	}

	if person.IsCombination == true {
		opponent, err := userRepository.GetByUid(person.OpponentUid)
		if err != nil {
			return err
		}

		err = userRepository.CancelMatchingStatus(opponent)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = userRepository.SoftDeleteUser(person)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// Transaction

	return nil
}

/**
 * ユーザーからの報告, ブロック機能
 */
func (s UserService) CreateUnauthorizedUser(c *gin.Context) error {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		unauthorizedUid user.UnauthorizedUser
	)

	uid := c.Request.Header.Get("Uid")

	// Transaction
	userRepository := repository.NewUserRepository(db)
	person, err := userRepository.GetByUid(uid)
	if err != nil {
		return err
	}

	opponent, err := userRepository.GetByUid(person.OpponentUid)
	if err != nil {
		return err
	}

	if err = userRepository.CancelMatchingStatus(person); err != nil {
		tx.Rollback()
		return err
	}

	if err = userRepository.CancelMatchingStatus(opponent); err != nil {
		tx.Rollback()
		return err
	}

	unauthorizedUid.UID = opponent.UID
	if c.Request.Header.Get("Block") == "true" {
		unauthorizedUid.Block = true
	} else {
		unauthorizedUid.Block = false
	}

	if err := userRepository.CreateUnAuthorizeUser(unauthorizedUid); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	// Transaction

	return nil
}

/**
 * n通以上メッセージのやり取りがあった場合に相性が良い組み合わせと考え
 * UserCompatibleに保存する
 */
func (s UserService) CreateCompatible(c *gin.Context) (user.InfoCompatible, error) {
	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		infoCompatible user.InfoCompatible
		uComb          user.UserCombination
		uInfo          user.UserInformation
		otherUinfo     user.UserInformation
	)

	if err := c.BindJSON(&uComb); err != nil {
		return infoCompatible, err
	}

	if err := tx.First(&uInfo, "uid=?", uComb.UID).Error; err != nil {
		tx.Rollback()
		return infoCompatible, domain.RaiseDBError()
	}
	if err := tx.First(&otherUinfo, "uid=?", uComb.OpponentUID).Error; err != nil {
		tx.Rollback()
		return infoCompatible, domain.RaiseDBError()
	}

	infoCompatible.SetInfoCompatible(uInfo.UID, otherUinfo.UID)

	if err := tx.Create(&infoCompatible).Error; err != nil {
		tx.Rollback()
		return infoCompatible, domain.RaiseDBError()
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return infoCompatible, domain.RaiseDBError()
	}

	return infoCompatible, nil
}
