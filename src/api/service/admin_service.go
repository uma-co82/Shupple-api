package service

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/uma-co82/Shupple-api/src/api/db"
)

/************************************************************
 *                         ADMIN                            *
 ************************************************************/
type PASS struct {
	Admin string
}

/**
 * User一覧
 */
func (s UserService) GetAllUser(c *gin.Context) ([]User, error) {
	uid := c.Request.Header.Get("Uid")

	var env PASS
	_ = envconfig.Process("", &env)

	if uid != env.Admin {
		return nil, RaiseError(403, "Forbidden", nil)
	}

	db := db.Init()
	tx := db.Begin()
	defer db.Close()
	var (
		users []User
	)

	if err := tx.Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, RaiseDBError()
	}

	return users, nil
}
