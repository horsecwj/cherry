package models

import (
	"time"
)

type Identity struct {
	CommonModel
	UserId      int       `json:"user_id"`                              // 所属用户
	Source      string    `json:"source" gorm:"type:varchar(32)"`       // Email or Phone, Wechat, Alipay
	Symbol      string    `json:"symbol" gorm:"type:varchar(64)"`       // Email address or Phone number, openid, uid
	AccessToken string    `json:"access_token" gorm:"type:varchar(64)"` // 授权token
	ExpiredAt   time.Time `json:"expired_at" gorm:"default:null"`       // 过期时间
}
