package models

import (
	"cherry/utils"
	"github.com/jinzhu/gorm"
)

type Device struct {
	CommonModel
	UserId    int    `json:"user_id"`                       // 设备所属用户
	IsUsed    bool   `json:"is_used"`                       // 是否已经使用
	Token     string `json:"token" gorm:"type:varchar(64)"` // 授权token
	PublicKey string `json:"-"`                             // 公钥，用于公钥验签
}

func (device *Device) InitializeToken() {
	device.Token = utils.RandStringRunes(64)
}

func (device *Device) BeforeCreate(db *gorm.DB) {
	var count int
	db.Model(&Device{}).Where("token = ?", device.Token).Count(&count)
	for count > 0 {
		device.Token = utils.RandStringRunes(64)
		db.Model(&Device{}).Where("token = ?", device.Token).Count(&count)
	}
}
