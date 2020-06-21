package models

import (
	"time"

	"cherry/utils"
	"github.com/jinzhu/gorm"
)

type Token struct {
	CommonModel
	Type         string    `gorm:"type:varchar(32)" json:"type"`  // 令牌类型 Tokens::Login, Tokens::AccessToken
	Token        string    `gorm:"type:varchar(64)" json:"token"` // 授权令牌
	UserId       int       `json:"user_id"`                       // 所属用户
	IsUsed       bool      `json:"is_used"`                       // 是否已使用
	ExpireAt     time.Time `gorm:"default:null" json:"expire_at"` // 过期时间
	LastVerifyAt time.Time `gorm:"default:null" json:"-"`         // 最后验证时间
}

type TokenAndApp struct {
	CommonModel
	TokenId   int
	ServiceId int
}

func (token *Token) BeforeCreate(db *gorm.DB) {
	count := 9
	for count > 0 {
		token.Token = utils.RandStringRunes(64)
		db.Model(&Token{}).Where("token = ?", token.Token).Count(&count)
	}
}

func (token *Token) InitializeAccessToken() {
	token.Token = utils.RandStringRunes(64)
	token.Type = "AccessToken"
	token.ExpireAt = time.Now().Add(time.Hour * 24 * 7)
}
