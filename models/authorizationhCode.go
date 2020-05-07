package models

import (
	"strconv"

	"cherry/utils"
)

type AuthorizationhCode struct {
	UserId int
	ServiceId  int
	Code   string
}

func (ac *AuthorizationhCode) Key() string {
	return "hex:redismodels::authorizationcode:" + strconv.Itoa(ac.ServiceId) + ":" + ac.Code
}

func (ac *AuthorizationhCode) InitCode() {
	ac.Code = utils.RandStringRunes(32)
}
