package models

import (
	"github.com/shopspring/decimal"
)

type Transfer struct {
	CommonModel
	NotifyTimes int             `json:"notify_times"`
	From        int             `json:"from"`
	To          int             `json:"to"`
	CurrencyId  int             `json:"-"`
	ServiceId       int             `json:"service_id"`
	Sn          string          `json:"sn" gorm:"size:32"`
	State       string          `json:"state" gorm:"default:'pending';size:12"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(32,16);default:null;"`
}

func (transfer *Transfer) IsDone() (stat bool) {
	if transfer.State == "done" {
		stat = true
	}
	return
}

func (transfer *Transfer) IsCanceled() (stat bool) {
	if transfer.State == "canceled" {
		stat = true
	}
	return
}
