package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/uma-co82/Shupple-api/src/api/domain"
	"github.com/uma-co82/Shupple-api/src/api/domain/user"
)

type UserRepository interface {
	GetByUid(string) (user.User, error)
	GetUserInformationByRelatedUser(user.User) (user.UserInformation, error)
	GetUserCombinationByBothUid(string, string) (user.UserCombination, error)
	GetShupple(int, int, int, int, string) ([]user.User, error)
	Update(user.User, user.User) error
	CreateUserCombination(user.UserCombination) error
	CancelMatchingStatus(user.User) error
	CreateUser(user.User) error
	SoftDeleteUser(user.User) error
	CreateUnAuthorizeUser(user.UnauthorizedUser) error
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (u *userRepository) GetByUid(uid string) (user.User, error) {
	var person user.User
	if err := u.conn.First(&person, "uid=?", uid).Error; err != nil {
		return person, domain.RaiseDBError()
	}
	return person, nil
}

func (u *userRepository) GetUserInformationByRelatedUser(person user.User) (user.UserInformation, error) {
	var userInformation user.UserInformation
	if err := u.conn.Model(&person).Related(&userInformation, "UserInformation").Error; err != nil {
		return userInformation, domain.RaiseDBError()
	}
	return userInformation, nil
}

func (u *userRepository) GetUserCombinationByBothUid(personUid, opponentUid string) (user.UserCombination, error) {
	var userCombination user.UserCombination
	if err := u.conn.Where("uid IN (?, ?) AND opponent_uid IN (?, ?)", personUid, opponentUid, personUid, opponentUid).First(&userCombination).Error; err != nil {
		return userCombination, domain.RaiseDBError()
	}
	return userCombination, nil
}

// select * from users where age BETWEEN 20 AND 30 AND sex=1 AND is_combination=false AND uid NOT IN (select opponent_uid from user_combinations where uid='自分のuid')
func (u *userRepository) GetShupple(ageLow, ageUpper, sex, residence int, uid string) ([]user.User, error) {
	var candidateUsers []user.User
	if err := u.conn.Where(
		"age BETWEEN ? AND ? AND sex=? AND is_combination=? AND uid NOT IN (select opponent_uid from user_combinations where uid=?) AND uid IN (select uid from user_informations where residence=?)",
		ageLow, ageUpper, sex, false, uid, residence).Find(&candidateUsers).Error; err != nil {
		return candidateUsers, domain.RaiseDBError()
	}
	return candidateUsers, nil
}

func (u *userRepository) Update(before user.User, after user.User) error {
	if err := u.conn.Model(&before).Updates(&after).Error; err != nil {
		return domain.RaiseDBError()
	}
	return nil
}

func (u *userRepository) CreateUserCombination(userCombination user.UserCombination) error {
	if err := u.conn.Create(&userCombination).Error; err != nil {
		return domain.RaiseDBError()
	}
	return nil
}

func (u *userRepository) CancelMatchingStatus(person user.User) error {
	var updateTarget = map[string]interface{}{"is_combination": false, "opponent_uid": nil}

	if err := u.conn.Model(&person).Updates(updateTarget).Error; err != nil {
		return domain.RaiseDBError()
	}
	return nil
}

func (u *userRepository) CreateUser(person user.User) error {
	if err := u.conn.Create(&person).Error; err != nil {
		return domain.RaiseDBError()
	}
	return nil
}

func (u *userRepository) SoftDeleteUser(person user.User) error {
	if err := u.conn.Delete(&person).Error; err != nil {
		return domain.RaiseDBError()
	}
	return nil
}

func (u *userRepository) CreateUnAuthorizeUser(unauthorizedUser user.UnauthorizedUser) error {
	if err := u.conn.Create(&unauthorizedUser).Error; err != nil {
		return domain.RaiseDBError()
	}
	return nil
}
