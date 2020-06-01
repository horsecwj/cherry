package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"

	"cherry/utils"
)

type Account struct {
	CommonModel

	UserId     int             `json:"user_id"`                                      // 所属用户
	CurrencyId int             `json:"currency_id"`                                  // 币种
	Balance    decimal.Decimal `gorm:"type:decimal(32,2);default:0;" json:"balance"` // 可用余额
	Locked     decimal.Decimal `gorm:"type:decimal(32,2);default:0;" json:"locked"`  // 锁定余额

	// 以下字段不存入数据库
	Options      []int    `sql:"-" json:"options"`
	CurrencyName string   `sql:"-" json:"currency_name"`
	Currency     Currency `sql:"-" json:"-"`
	Fixed        int      `sql:"-" json:"fixed"`
}

func DefaultAccount(db *gorm.DB) *gorm.DB {
	return db.Where("accounts.modifiable_id = ?", 0).Where("accounts.modifiable_type = ?", "Normal")
}

func (account *Account) AfterUpdate() {
	account.notifyAccount()
}
func (account *Account) AfterCreate() {
	account.notifyAccount()
}

func (account *Account) notifyAccount() {
	b, err := json.Marshal(*account)
	if err != nil {
		log.Println(err)
	}
	utils.PublishToPubSubChannels(NotifyAccountWithRedis, &b)
}

func (account *Account) SetOptions() {
	if account.Currency.Id == 0 {
		account.Currency, _ = FindCurrencyById(account.CurrencyId)
	}
	account.Options = account.Currency.Options
	account.Fixed = account.Currency.Fixed
}

func (account *Account) Amount() (amount decimal.Decimal) {
	amount = account.Balance.Add(account.Locked)
	return
}

// 1
func (account *Account) UnlockFunds(db *utils.GormDB, amount decimal.Decimal, reason, modifiableId int, modifiableType string) (err error) {
	if amount.LessThanOrEqual(decimal.Zero) || amount.GreaterThan(account.Locked) {
		err = fmt.Errorf("cannot unlock funds (amount: %v)", amount)
		return
	}

	err = account.changeBalanceAndLocked(db, amount, amount.Neg())
	if err != nil {
		return
	}
	opts := map[string]string{
		"reason":          strconv.Itoa(reason),
		"modifiable_id":   strconv.Itoa(modifiableId),
		"modifiable_type": modifiableType,
	}
	err = account.after(db, FUNS["UnlockFunds"], amount, opts)
	return
}

// 2
func (account *Account) LockFunds(db *utils.GormDB, amount decimal.Decimal, reason, modifiableId int, modifiableType string) (err error) {
	if amount.LessThanOrEqual(decimal.Zero) || amount.GreaterThan(account.Balance) {
		err = fmt.Errorf("cannot lock funds (amount: %v)", amount)
		return
	}

	err = account.changeBalanceAndLocked(db, amount.Neg(), amount)
	if err != nil {
		return
	}
	opts := map[string]string{
		"reason":          strconv.Itoa(reason),
		"modifiable_id":   strconv.Itoa(modifiableId),
		"modifiable_type": modifiableType,
	}
	err = account.after(db, FUNS["LockFunds"], amount, opts)
	return
}

// 3
func (account *Account) PlusFunds(db *utils.GormDB, amount, fee decimal.Decimal, reason, modifiableId int, modifiableType string) (err error) {
	if amount.LessThan(decimal.Zero) || fee.GreaterThan(amount) {
		err = fmt.Errorf("cannot add funds (amount: %v)", amount)
		return
	}
	err = account.changeBalanceAndLocked(db, amount, decimal.Zero)
	if err != nil {
		return
	}
	opts := map[string]string{
		"fee":             fee.String(),
		"reason":          strconv.Itoa(reason),
		"modifiable_id":   strconv.Itoa(modifiableId),
		"modifiable_type": modifiableType,
	}
	err = account.after(db, FUNS["PlusFunds"], amount, opts)
	return
}

// 4
func (account *Account) SubFunds(db *utils.GormDB, amount, fee decimal.Decimal, reason, modifiableId int, modifiableType string) (err error) {
	if amount.LessThan(decimal.Zero) || fee.GreaterThan(amount) {
		err = fmt.Errorf("cannot add funds (amount: %v)", amount)
		return
	}
	err = account.changeBalanceAndLocked(db, amount.Neg(), decimal.Zero)
	if err != nil {
		return
	}
	opts := map[string]string{
		"fee":             fee.String(),
		"reason":          strconv.Itoa(reason),
		"modifiable_id":   strconv.Itoa(modifiableId),
		"modifiable_type": modifiableType,
	}
	err = account.after(db, FUNS["SubFunds"], amount, opts)
	return
}

// 5
func (account *Account) UnlockedAndSubFunds(db *utils.GormDB, amount, locked, fee decimal.Decimal, reason, modifiableId int, modifiableType string) (err error) {
	if amount.LessThan(decimal.Zero) || amount.GreaterThan(locked) {
		err = fmt.Errorf("cannot unlock and subtract funds (amount: %v)", amount)
		return
	}
	if locked.LessThanOrEqual(decimal.Zero) {
		err = fmt.Errorf("invalid lock amount")
		return
	}
	if locked.GreaterThan(account.Locked) {
		err = fmt.Errorf("Account# %v invalid lock amount (amount: %v, locked: %v, self.locked: %v)", account.Id, amount, locked, account.Locked)
		return
	}
	err = account.changeBalanceAndLocked(db, locked.Sub(amount), locked.Neg())
	if err != nil {
		return
	}
	opts := map[string]string{
		"fee":             fee.String(),
		"locked":          locked.String(),
		"reason":          strconv.Itoa(reason),
		"modifiable_id":   strconv.Itoa(modifiableId),
		"modifiable_type": modifiableType,
	}
	err = account.after(db, FUNS["UnlockedAndSubFunds"], amount, opts)
	return
}

// 6
func (account *Account) PlusAndLockFunds(db *utils.GormDB, amount, locked, fee decimal.Decimal, reason, modifiableId int, modifiableType string) (err error) {
	if amount.LessThan(decimal.Zero) || fee.GreaterThan(amount) {
		err = fmt.Errorf("cannot add funds (amount: %v)", amount)
		return
	}
	if locked.LessThanOrEqual(decimal.Zero) {
		err = fmt.Errorf("invalid lock amount")
		return
	}
	err = account.changeBalanceAndLocked(db, locked.Sub(amount).Neg(), locked)
	if err != nil {
		return
	}
	opts := map[string]string{
		"fee":             fee.String(),
		"locked":          locked.String(),
		"reason":          strconv.Itoa(reason),
		"modifiable_id":   strconv.Itoa(modifiableId),
		"modifiable_type": modifiableType,
	}
	err = account.after(db, FUNS["PlusAndLockFunds"], amount, opts)
	return
}

func (account *Account) changeBalanceAndLocked(db *utils.GormDB, deltaB, deltaL decimal.Decimal) (err error) {
	db.Set("gorm:query_option", "FOR UPDATE").First(&account, account.Id)
	balance := account.Balance
	account.Balance = account.Balance.Add(deltaB)
	locked := account.Locked
	account.Locked = account.Locked.Add(deltaL)
	updateSql := fmt.Sprintf("UPDATE accounts SET balance = balance + %v, locked = locked + %v WHERE accounts.id = %v AND balance = ? AND locked = ?", deltaB, deltaL, account.Id, balance, locked)
	accountresult := db.Exec(updateSql)
	if accountresult.RowsAffected != 1 {
		err = fmt.Errorf("Update row failed.")
	}
	return
}

// global
func (account *Account) after(db *utils.GormDB, fun int, amount decimal.Decimal, opts map[string]string) (err error) {
	var fee decimal.Decimal
	if opts["fee"] != "" {
		fee, _ = decimal.NewFromString(opts["fee"])
	}
	if opts["reason"] == "" {
		opts["reason"] = strconv.Itoa(UNKNOWN)
	}
	attributes := map[string]string{
		"fun":             strconv.Itoa(fun),
		"fee":             fee.String(),
		"reason":          opts["reason"],
		"amount":          account.Amount().String(),
		"currency_id":     strconv.Itoa(account.CurrencyId),
		"user_id":         strconv.Itoa(account.UserId),
		"account_id":      strconv.Itoa(account.Id),
		"modifiable_id":   opts["modifiable_id"],
		"modifiable_type": opts["modifiable_type"],
	}
	attributes["locked"], attributes["balance"], err = account.computeLockedAndBalance(fun, amount, opts)
	if err != nil {
		log.Println(err)
		return
	}
	err = account.optimisticallyLockAccountAndCreate(db, account.Balance, account.Locked, attributes)
	return
}

func (account *Account) computeLockedAndBalance(fun int, amount decimal.Decimal, opts map[string]string) (locked, balance string, err error) {
	switch fun {
	case 1:
		locked = amount.Neg().String()
		balance = amount.String()
	case 2:
		locked = amount.String()
		balance = amount.Neg().String()
	case 3:
		locked = "0"
		balance = amount.String()
	case 4:
		locked = "0"
		balance = amount.Neg().String()
	case 5:
		l, _ := decimal.NewFromString(opts["locked"])
		locked = l.Neg().String()
		balance = l.Sub(amount).String()
	case 6:
		l, _ := decimal.NewFromString(opts["locked"])
		locked = l.String()
		balance = amount.Sub(l).String()
	default:
		err = fmt.Errorf("forbidden account operation")
	}
	return
}

func (account *Account) optimisticallyLockAccountAndCreate(db *utils.GormDB, balance, locked decimal.Decimal, attrs map[string]string) (err error) {
	if attrs["account_id"] == "" {
		err = fmt.Errorf("account must be specified")
	}
	attrs["created_at"] = time.Now().Format("2006-01-02 15:04:05")
	attrs["updated_at"] = attrs["created_at"]

	sql := `INSERT INTO account_versions
  (user_id, account_id, reason, balance, locked, amount, modifiable_id, modifiable_type, currency_id, fun, created_at, updated_at)
  SELECT ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?  FROM accounts WHERE accounts.balance = ? AND accounts.locked = ? AND accounts.id = ?`
	result := db.Exec(sql,
		attrs["user_id"],
		attrs["account_id"],
		attrs["reason"],
		attrs["balance"],
		attrs["locked"],
		attrs["amount"],
		attrs["modifiable_id"],
		attrs["modifiable_type"],
		attrs["currency_id"], attrs["fun"],
		attrs["created_at"],
		attrs["updated_at"],
		balance,
		locked,
		attrs["account_id"],
	)
	if result.RowsAffected != 1 {
		err = fmt.Errorf("Insert row failed.")
	}
	return
}
