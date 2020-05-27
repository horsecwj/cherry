package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"cherry/utils"
)

type User struct {
	CommonModel
	RoleId         int    `json:"role_id"`                          // 角色
	Sn             string `gorm:"type:varchar(16)" json:"sn"`       // 唯一编号
	PasswordDigest string `gorm:"type:varchar(64)" json:"-"`        // 经加密的密码
	Nickname       string `gorm:"type:varchar(32)" json:"nickname"` // 昵称
	Ancestry       string `gorm:"type:varchar(32)" json:"ancestry"` // 邀请关系
	State          int    `gorm:"default:null" json:"state"`        // 状态
	Activated      bool   `gorm:"default:null" json:"activated"`    // 是否已激活
	Disabled       bool   `json:"disabled"`                         // 是否禁用
	ApiDisabled    bool   `json:"api_disabled"`                     // 是否禁用API

	SmsValidated bool    `sql:"-" json:"sms_validated"`
	Password     string  `sql:"-" json:"-"`
	Tokens       []Token `sql:"-" json:"tokens"`
	Role         Role    `sql:"-" json:"role"`
}

type OauthUserInfo struct {
	Uid          int    `json:"uid"`
	Sn           string `json:"sn"`
	Name         string `json:"name"`
	Ancestry     string `json:"ancestry"`
	InviteUrl    string `json:"invite_url"`
	Role         string `json:"role"`
	Activated    bool   `json:"activated"`
	SmsValidated bool   `json:"sms_validated"`
}

func (user *User) BeforeCreate(db *gorm.DB) {
	count := 4
	for count > 0 {
		user.Sn = "CHE" + utils.RandStringRunes(10) + "RRY"
		db.Model(&User{}).Where("sn = ?", user.Sn).Count(&count)
	}
}

func (user *User) AfterFind(db *gorm.DB) {
	user.setRole()
}

func (user *User) CompareHashAndPassword() bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(user.Password))
	if err == nil {
		return true
	}
	return false
}

func (user *User) SetPasswordDigest() {
	b, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.PasswordDigest = string(b[:])
}

func (user *User) setRole() {
	for _, role := range AllRoles {
		if user.RoleId == 0 && role.Name == "Guest" {
			user.Role = role
			return
		} else if user.RoleId == role.Id {
			user.Role = role
			return
		}
	}
}
