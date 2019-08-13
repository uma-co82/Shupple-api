package handler

import (
	"net/http"

	"github.com/holefillingco-ltd/Shupple-api/src/api/model"
	"github.com/jinzhu/gorm"
)

func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	defer db.Close()
}
