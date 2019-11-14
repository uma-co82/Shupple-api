package service

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/uma-co82/Shupple-api/src/api/domain"
	"github.com/uma-co82/Shupple-api/src/api/domain/user"
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/db"
)

/************************************************************
 *                         ADMIN                            *
 ************************************************************/
type PASS struct {
	ADMIN string
}

/**
 * User一覧
 */
func (s UserService) GetAllUser(c *gin.Context) ([]user.User, error) {
	uid := c.Request.Header.Get("Uid")

	var env PASS
	_ = envconfig.Process("", &env)

	if uid != env.ADMIN {
		return nil, domain.RaiseError(403, "Forbidden", nil)
	}

	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		users []user.User
	)

	if err := tx.Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, domain.RaiseDBError()
	}

	return users, nil
}
