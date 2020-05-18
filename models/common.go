package models

import (
	"time"

	"cherry/utils"
)

const (
	NotifyAccountWithRedis = "notify:account"
)

var (
	UNKNOWN              = 0   // 未知类型
	FIX                  = 1   // 修复
	TRANSFER             = 700 // 转账入金
	TRANSFER_BACK_LOCK   = 705 // 转账出金锁币
	TRANSFER_BACK_UNLOCK = 706 // 转账出金锁定的币解锁
	TRANSFER_BACK        = 710 // 转账出金，扣款
	GRAINT_LOCK          = 750 // 拨款锁币
	GRAINT_UNLOCK        = 751 // 拨款解锁
	GRAINT_SUB_LOCK      = 752 // 拨款扣除锁币

	FUNS = map[string]int{
		"UnlockFunds":         1,
		"LockFunds":           2,
		"PlusFunds":           3,
		"SubFunds":            4,
		"UnlockedAndSubFunds": 5,
		"PlusAndLockFunds":    6,
	}

	// transfer state
	Pending  = 1
	Canceled = 10
	Done     = 20

	AllCurrencies []Currency
)

type CommonModel struct {
	Id             int       `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	CreatedAtStamp int64     `sql:"-" json:"created_at"`
	UpdatedAtStamp int64     `sql:"-" json:"updated_at"`
}

func (cm *CommonModel) InitializeTimestamp() {
	cm.CreatedAtStamp = cm.CreatedAt.UnixNano() / 1000000
	cm.UpdatedAtStamp = cm.UpdatedAt.UnixNano() / 1000000
}

func AutoMigrations() {
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()

	// account
	mainDB.AutoMigrate(&Account{})
	mainDB.Model(&Account{}).AddUniqueIndex("index_accounts_on_idx0", "user_id", "currency_id", "modifiable_type", "modifiable_id")

	// account_version
	mainDB.AutoMigrate(&AccountVersion{})
	mainDB.Model(&AccountVersion{}).AddIndex("index_account_versions_on_user_id_and_reason", "user_id", "reason")
	mainDB.Model(&AccountVersion{}).AddIndex("index_account_versions_on_account_id_and_reason", "account_id", "reason")
	mainDB.Model(&AccountVersion{}).AddIndex("index_account_versions_on_currency_id_and_created_at", "currency_id", "created_at")
	mainDB.Model(&AccountVersion{}).AddUniqueIndex("account_versions_idx0", "account_id", "modifiable_id", "modifiable_type", "reason")

	// account_version_check_point
	mainDB.AutoMigrate(&AccountVersionCheckPoint{})
	mainDB.Model(&AccountVersionCheckPoint{}).AddIndex("index_account_version_check_points_on_account_id", "account_id")

	// config
	mainDB.AutoMigrate(&Config{})

	// currency
	mainDB.AutoMigrate(&Currency{})
	mainDB.Model(&Currency{}).AddIndex("index_currencies_on_visible", "visible")

	// device
	mainDB.AutoMigrate(&Device{})

	// identity
	mainDB.AutoMigrate(&Identity{})
	mainDB.Model(&Identity{}).AddUniqueIndex("index_identity_on_source_and_symbol", "source", "symbol")

	// recharge
	mainDB.AutoMigrate(&Recharge{})
	mainDB.Model(&Recharge{}).AddUniqueIndex("index_recharges_on_sn", "sn")

	// service
	mainDB.AutoMigrate(&Service{})

	// token
	mainDB.AutoMigrate(&Token{})

	// transfer
	mainDB.AutoMigrate(&Transfer{})
	mainDB.Model(&Transfer{}).AddUniqueIndex("index_transfers_on_sn_and_service_id", "sn", "service_id")
	mainDB.AutoMigrate(&TransferNotifyLog{})

	mainDB.AutoMigrate(&TwoFactor{})

	// user
	mainDB.AutoMigrate(&User{})

}
