package models

import (
	envConfig "cherry/config"
	"cherry/utils"
)

type Config struct {
	CommonModel
	Description string `json:"description"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

func InitConfigInDB(db *utils.GormDB) {
	var configs []Config
	db.Find(&configs)
	envConfig.CurrentEnv.ConfigInDB = map[string]string{}
	for _, config := range configs {
		envConfig.CurrentEnv.ConfigInDB[config.Key] = config.Value
	}
}
