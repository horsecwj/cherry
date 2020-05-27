package models

import (
	"encoding/json"

	"cherry/utils"
)

type Service struct {
	CommonModel
	UserId      int    `json:"-"`
	Inside      bool   `json:"-"`
	Name        string `json:"name"`
	AppKey      string `json:"-" gorm:"type:varchar(64)"`
	AasmState   string `json:"-" gorm:"type:varchar(32);default:'pending'"`
	CompanyName string `json:"company_name"`
	CustomKey   string `json:"custom_key" gorm:"type:text"`
	PrivateKey  string `json:"-" gorm:"type:text"`
	PublicKey   string `json:"public_key" gorm:"type:text"`
	Hosts       string `json:"-" gorm:"type:text"`
	GrantUrl    string `json:"-"`

	CallbackHosts []string `json:"callback_hosts" sql:"-"`
}

func (as *Service) BeforeCreate() {
	if as.AppKey == "" {
		as.AppKey = utils.RandStringRunes(64)
	}
}

func (as *Service) AfterFind() {
	json.Unmarshal([]byte(as.Hosts), &as.CallbackHosts)
}

func (as *Service) ValidataHost(host string) (re bool) {
	for _, h := range as.CallbackHosts {
		if h == host {
			re = true
			return
		}
	}
	return
}

func (as *Service) CanNotGrant() (canNot bool) {
	if as.GrantUrl != "" || as.Inside == false {
		canNot = true
	}
	return
}
