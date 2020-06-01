package models

import (
	"github.com/shopspring/decimal"

	"cherry/utils"
)

type Transfer struct {
	CommonModel
	From       int             `json:"from"`
	To         int             `json:"to"`
	CurrencyId int             `json:"-"`
	ServiceId  int             `json:"service_id"`
	Grant      bool            `json:"grant" gorm:"default:false"`
	CenterSn   string          `json:"center_sn" gorm:"size:32"`
	Sn         string          `json:"sn" gorm:"size:32"`
	State      string          `json:"state" gorm:"default:'pending';size:12"`
	Amount     decimal.Decimal `json:"amount" gorm:"type:decimal(32,2);default:null;"`
}

func (transfer *Transfer) BeforeCreate() {
	if transfer.CenterSn == "" {
		transfer.CenterSn = utils.RandStringRunes(32)
	}
}

func (transfer *Transfer) IsDone() bool {
	return transfer.State == "done"
}

func (transfer *Transfer) IsCanceled() bool {
	return transfer.State == "canceled"
}

func (transfer *Transfer) IsGranting() bool {
	return transfer.Grant && transfer.State == "granting"
}
