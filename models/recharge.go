package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"

	"cherry/utils"
)

type Recharge struct {
	CommonModel
	UserId     int             `json:"user_id"`
	CurrencyId int             `json:"currency_id"`
	Title      string          `json:"title" gorm:"size:64"`
	Sn         string          `json:"sn" gorm:"size:64"`
	OutTradeNo string          `json:"out_trade_no" gorm:"size:64"`
	Source     string          `json:"source" gorm:"size:16"` // Alipay, Wechat, Upop, Upmp
	Amount     decimal.Decimal `json:"amount" gorm:"type:decimal(32,2);default:0;"`
	State      int             `json:"state" gorm:"default:1"`
}

func (recharge *Recharge) BeforeCreate(db *gorm.DB) {
	count := 5
	for count > 0 {
		recharge.Sn = utils.RandStringRunes(64)
		db.Model(&Recharge{}).Where("sn = ?", recharge.Sn).Count(&count)
	}
}
