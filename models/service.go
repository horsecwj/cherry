package models

import (
	"encoding/json"
)

type Service struct {
	CommonModel
	UserId    int    `json:"user_id"`
	Inside      bool   `json:"inside"`
	Name        string `json:"name"`
	AppKey      string `json:"-"`
	AppSecret   string `json:"-"`
	AasmState   string `json:"-" gorm:"default:pending"`
	CompanyName string `json:"company_name"`
	CustomKey   string `json:"-" gorm:"type:text"`
	PrivateKey  string `json:"-" gorm:"type:text"`
	PublicKey   string `json:"-" gorm:"type:text"`
	Hosts       string `json:"-" gorm:"type:text"`
	GrantUrl    string `json:"-"`

	CallBackHosts []string `sql:"-"`
}

func (as *Service) AfterFind() {
	json.Unmarshal([]byte(as.Hosts), &as.CallBackHosts)
}

func (as *Service) ValidataHost(host string) (re bool) {
	for _, h := range as.CallBackHosts {
		if h == host {
			re = true
			return
		}
	}
	return
}

func (as *Service) CanNotGrant() (canNot bool) {
	if as.GrantUrl != "" {
		canNot = true
	}
	return
}
