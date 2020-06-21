package models

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"cherry/utils"
)

type Currency struct {
	CommonModel
	Key         string `json:"key"`                                                                 // 币的唯一标示
	Symbol      string `json:"symbol"`                                                              // 币的简称
	Logo        string `json:"logo"`                                                                // 币的图标
	Visible     bool   `json:"visible"`                                                             // 是否可用
	Fixed       int    `json:"fixed" gorm:"default:2"`                                              // 小数位精度
	OptionsJson string `json:"-" gorm:"column:options;type:varchar(64);default:\"[10,20,50,100]\""` // 转账快捷选项JSON
	Options     []int  `sql:"-" json:"options"`                                                     // 转账快捷选项
}

func InitAllCurrencies(db *utils.GormDB) {
	db.FirstOrCreate(&Currency{}, map[string]interface{}{"symbol": "cny", "key": "CNY", "visible": true})
	db.Where("visible = ?", true).Find(&AllCurrencies)
}

func FindCurrencyById(id int) (Currency, error) {
	for _, currency := range AllCurrencies {
		if currency.Id == id {
			return currency, nil
		}
	}
	var currency Currency
	return currency, fmt.Errorf("No currency can be found.")
}

func FindCurrencyBySymbol(symbol string) (Currency, error) {
	for _, currency := range AllCurrencies {
		if currency.Symbol == symbol {
			return currency, nil
		}
	}
	var currency Currency
	return currency, fmt.Errorf("No currency can be found.")
}

func (c *Currency) AfterFind(db *gorm.DB) {
	json.Unmarshal([]byte(c.OptionsJson), &c.Options)
}
