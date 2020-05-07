package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"cherry/utils"
)

type User struct {
	CommonModel
	Sn               string          `gorm:"type:varchar(16)" json:"sn"`                      // 唯一编号
	PasswordDigest   string          `gorm:"type:varchar(64)" json:"-"`                       // 经加密的密码
	Nickname         string          `gorm:"type:varchar(32)" json:"nickname"`                // 昵称
	Ancestry         string          `gorm:"type:varchar(32)" json:"ancestry"`                // 邀请关系
	State            int             `gorm:"default:null" json:"state"`                       // 状态
	Activated        bool            `gorm:"default:null" json:"activated"`                   // 是否已激活
	Disabled         bool            `json:"disabled"`                                        // 是否禁用
	ApiDisabled      bool            `json:"api_disabled"`                                    // 是否禁用API
	ActivityPriority bool            `gorm:"default:true" json:"activity_priority"`           // 优先使用空投币
	ProfitRatio      decimal.Decimal `gorm:"type:decimal(4,3);default:1" json:"profit_ratio"` // 赔率

	SmsValidated  bool            `sql:"-" json:"sms_validated"`
	OrderFeeRatio decimal.Decimal `sql:"-" json:"order_fee_ratio"`
	Password      string          `sql:"-" json:"-"`
	Tokens        []Token         `sql:"-" json:"tokens"`
}

type OauthUserInfo struct {
	Uid       int    `json:"uid"`
	Sn        string `json:"sn"`
	Name      string `json:"name"`
	Ancestry  string `json:"ancestry"`
	InviteUrl string `json:"invite_url"`

	Activated    bool `json:"activated"`
	SmsValidated bool `json:"sms_validated"`
}

func (user *User) GenerateSn() {
	user.Sn = "CHE" + utils.RandStringRunes(10) + "RRY"
}

func (user *User) BeforeCreate(db *gorm.DB) {
	var count int
	db.Model(&User{}).Where("sn = ?", user.Sn).Count(&count)
	for count > 0 {
		user.GenerateSn()
		db.Model(&User{}).Where("sn = ?", user.Sn).Count(&count)
	}
}

func (user *User) AfterFind(db *gorm.DB) {
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
