package models

import (
	"time"

	"cherry/utils"
	"github.com/jinzhu/gorm"
)

type Token struct {
	CommonModel
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

func (token *Token) InitializeLoginToken() {
	token.Token = utils.RandStringRunes(64)
	token.ExpireAt = time.Now().Add(time.Hour * 24 * 7)
}

func (token *Token) BeforeCreate(db *gorm.DB) {
	var count int
	db.Model(&Token{}).Where("token = ?", token.Token).Count(&count)
	for count > 0 {
		token.Token = utils.RandStringRunes(64)
		db.Model(&Token{}).Where("token = ?", token.Token).Count(&count)
	}
}

func (token *Token) InitializeAccessToken() {
	token.Token = utils.RandStringRunes(64)
	token.ExpireAt = time.Now().Add(time.Hour * 24 * 7)
}
