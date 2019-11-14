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
