package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"

	"cherry/initializers"
	. "cherry/models"
	"cherry/utils"
	"cherry/workers/sneakerWorkers"
)

// 拨款至某个应用
// 参数: app_key, transfer_id, currency, amount
func PostGrantCreate(context echo.Context) error {
	user := context.Get("current_user").(User)
	params := context.Get("params").(map[string]string)
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	service := context.Get("current_service").(Service)
	if service.Id == 0 {
		if mainDB.Where("aasm_state = ? AND app_key = ?", "verified", params["app_key"]).First(&service).RecordNotFound() {
			return utils.BuildError("1102")
		}
	}
	if service.CanNotGrant() {
		return utils.BuildError("1106")
	}
	amount, _ := decimal.NewFromString(params["amount"])
	amount = amount.Truncate(8)
	if amount.LessThanOrEqual(decimal.Zero) {
		return utils.BuildError("1105")
	}
	currency, e := FindCurrencyBySymbol(params["currency"])
	if e != nil {
		return utils.BuildError("1107")
	}
	transfer := Transfer{
		From:       user.Id,
		To:         service.UserId,
		Amount:     amount,
		CurrencyId: currency.Id,
		ServiceId:  service.Id,
		State:      "granting",
		Grant:      true,
	}
	err := tryLockFunds(&transfer, 9)
	if err != nil {
		return err
	}
	sendMessageToGrantNotifyUrl(&map[string]string{
		"id":         strconv.Itoa(transfer.Id),
		"from":       strconv.Itoa(transfer.From),
		"to":         strconv.Itoa(transfer.To),
		"amount":     transfer.Amount.String(),
		"currency":   currency.Code,
		"service_id": strconv.Itoa(service.Id),
		"notify_url": service.GrantUrl,
		"timestamp":  strconv.Itoa(int(time.Now().Unix())),
	})

	response := utils.SuccessResponse
	response.Body = transfer
	return context.JSON(http.StatusOK, response)
}

// 应用确认收到拨款，并写入SN
// 参数: app_key, transfer_id, currency, amount, sn
func GrantConfirm(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var transfer Transfer
	service := context.Get("current_service").(Service)
	if service.Id == 0 {
		if mainDB.Where("aasm_state = ? AND app_key = ?", "verified", params["app_key"]).First(&service).RecordNotFound() {
			return utils.BuildError("1102")
		}
	}
	if mainDB.Where("id = ?", params["transfer_id"]).First(&transfer).RecordNotFound() {
		return utils.BuildError("1108")
	}
	if transfer.ServiceId != service.Id {
		return utils.BuildError("1109")
	}
	if !transfer.IsDone() {
		var currency Currency
		if mainDB.Where("code = ?", params["currency"]).First(&currency).RecordNotFound() {
			return utils.BuildError("1107")
		}
		amount, _ := decimal.NewFromString(params["amount"])
		if transfer.Amount.Equal(amount) {
			return utils.BuildError("1110")
		}
		if transfer.CurrencyId == currency.Id {
			return utils.BuildError("1111")
		}
		transfer.Sn = params["sn"]
		err := grantToTargetAccount(&transfer, 50)
		if err != nil {
			return err
		}

	}
	return context.String(http.StatusOK, "success")
}

// 尝试锁币
func tryLockFunds(transfer *Transfer, times int) error {
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var account Account
	if mainDB.Where("user_id = ?", (*transfer).From).Where("currency = ?", (*transfer).CurrencyId).First(&account).RecordNotFound() {
		account.UserId = (*transfer).From
		account.CurrencyId = (*transfer).CurrencyId
		account.Balance = decimal.Zero
		account.Locked = decimal.Zero
		mainDB.Save(&account)
		mainDB.DbCommit()
		return utils.BuildError("1103")
	}

	if account.Balance.Sub((*transfer).Amount).IsNegative() {
		return utils.BuildError("1103")
	}
	mainDB.Create(transfer)
	if (*transfer).Id == 0 {
		return utils.BuildError("1112")
	}
	ferr := account.LockFunds(mainDB, (*transfer).Amount, GRAINT_LOCK, (*transfer).Id, "Transfer")
	if ferr == nil {
		mainDB.DbCommit()
		return nil
	}
	mainDB.DbRollback()
	if times > 0 {
		(*transfer).Id = 0
		return tryLockFunds(transfer, times-1)
	}
	return utils.BuildError("1104")
}

func grantToTargetAccount(transfer *Transfer, times int) error {
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var fromAccount, toAccount Account
	if mainDB.Where("user_id = ?", (*transfer).From).Where("currency = ?", (*transfer).CurrencyId).First(&fromAccount).RecordNotFound() {
		fromAccount.UserId = (*transfer).From
		fromAccount.CurrencyId = (*transfer).CurrencyId
		fromAccount.Balance = decimal.Zero
		fromAccount.Locked = decimal.Zero
		mainDB.Save(&fromAccount)
		mainDB.DbCommit()
		return utils.BuildError("1103")
	}
	if mainDB.Where("user_id = ?", (*transfer).To).Where("currency = ?", (*transfer).CurrencyId).First(&toAccount).RecordNotFound() {
		toAccount.UserId = (*transfer).To
		toAccount.CurrencyId = (*transfer).CurrencyId
		toAccount.Balance = decimal.Zero
		toAccount.Locked = decimal.Zero
		mainDB.Save(&toAccount)
	}

	if fromAccount.Locked.Sub((*transfer).Amount).IsNegative() {
		return utils.BuildError("1103")
	}
	(*transfer).State = "done"
	ferr := fromAccount.UnlockedAndSubFunds(mainDB, (*transfer).Amount, (*transfer).Amount, decimal.Zero, GRAINT_SUB_LOCK, (*transfer).Id, "Transfer")
	terr := toAccount.PlusFunds(mainDB, (*transfer).Amount, decimal.Zero, TRANSFER, (*transfer).Id, "Transfer")
	mainDB.Save(transfer)
	if ferr == nil && terr == nil {
		mainDB.DbCommit()
		return nil
	}
	mainDB.DbRollback()
	if times > 0 {
		(*transfer).Id = 0
		return grantToTargetAccount(transfer, times-1)
	}
	return utils.BuildError("1104")
}

var grantRoutingKey, grantCancelQueue string

func sendMessageToGrantNotifyUrl(message *map[string]string) {
	if grantRoutingKey == "" {
		for _, worker := range sneakerWorkers.AllWorkers {
			if worker.Name == "GrantNotifyWorker" {
				grantRoutingKey = worker.RoutingKey
			}
		}
	}
	if grantCancelQueue == "" {
		for _, worker := range sneakerWorkers.AllWorkers {
			if worker.Name == "GrantCancelWorker" {
				grantCancelQueue = worker.GetQueue() + ".delay.1"
			}
		}
	}
	b, err := json.Marshal(*message)
	if err != nil {
		log.Println(err)
	}
	err = initializers.PublishMessageWithRouteKey(initializers.AmqpGlobalConfig.Exchange["default"]["name"], grantRoutingKey, "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		log.Println(err)
	}
	err = initializers.PublishMessageToQueue(grantCancelQueue, "text/plain", &b, amqp.Table{}, amqp.Persistent)
	if err != nil {
		log.Println(err)
	}
}
