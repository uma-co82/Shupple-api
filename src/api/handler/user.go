package handler

import (
	"net/http"

	"../model"
	"github.com/jinzhu/gorm"
)

func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	defer db.Close()
}
