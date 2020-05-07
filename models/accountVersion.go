package models

import (
	"cherry/utils"
	"github.com/shopspring/decimal"
)

type AccountVersion struct {
	CommonModel
	UserId         int             `json:"user_id"`                                       // 所属用户
	AccountId      int             `json:"account_id"`                                    // 所属账号
	Reason         int             `json:"-"`                                             // 原因
	Balance        decimal.Decimal `gorm:"type:decimal(32,16);default:0;" json:"balance"` // 可用余额改变量
	Locked         decimal.Decimal `gorm:"type:decimal(32,16);default:0;" json:"locked"`  // 锁定余额改变量
	Amount         decimal.Decimal `gorm:"type:decimal(32,16);default:0;" json:"amount"`  // 修改完成后的量
	ModifiableId   int             `json:"modifiable_id"`                                 // 本次修改所属的order_id或者transfer_id，以及其它id
	ModifiableType string          `json:"modifiable_type"`                               // Order或Transfer,以及其它
	CurrencyId     int             `json:"currency_id"`                                   // 币种
	Fun            int             `json:"fun"`                                           // 流转函数

	// 以下字段不存入数据库
	CurrencySymbol string `sql:"-" json:"currency_symbol"`
	CurrencyLogo   string `sql:"-" json:"logo"`
	Sn             string `sql:"-" json:"sn"`
	ReasonStr      string `sql:"-" json:"reason"`
}

func (av *AccountVersion) InitForAPI(db *utils.GormDB) {
	av.setCurrency()
	av.SetReasonStr()
	av.InitializeTimestamp()
}

func (av *AccountVersion) setCurrency() {
	currency, _ := FindCurrencyById(av.CurrencyId)
	av.CurrencyLogo = currency.Logo
	av.CurrencySymbol = currency.Symbol
}

func (av *AccountVersion) setModifiable(db *utils.GormDB) {
	if av.ModifiableType == "Transfer" {
		db.Model(&Transfer{}).Where("id = ?", av.ModifiableId).Select("sn").Scan(&av)
	}
}

func (av *AccountVersion) SetReasonStr() {
	if av.Reason == TRANSFER_BACK_LOCK {
		av.ReasonStr = "TRANSFER_BACK_LOCK"
	} else if av.Reason == TRANSFER {
		av.ReasonStr = "TRANSFER"
	} else if av.Reason == TRANSFER_BACK_UNLOCK {
		av.ReasonStr = "TRANSFER_BACK_UNLOCK"
	} else if av.Reason == UNKNOWN {
		av.ReasonStr = "UNKNOWN"
	}
}
