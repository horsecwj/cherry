package models

import (
	"cherry/utils"
)

type Role struct {
	CommonModel
	Name string `json:"name"`
}

func InitAllRoles(db *utils.GormDB) {
	db.FirstOrCreate(&Role{}, map[string]interface{}{"name": "Guest"})
	db.FirstOrCreate(&Role{}, map[string]interface{}{"name": "Admin"})
	db.FirstOrCreate(&Role{}, map[string]interface{}{"name": "Merchant"})
	db.Find(&AllRoles)
}
