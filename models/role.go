package models

import (
	"cherry/utils"
)

type Role struct {
	CommonModel
	Name string `json:"name"`
}

func InitAllRoles(db *utils.GormDB) {
	db.Find(&AllRoles)
}
