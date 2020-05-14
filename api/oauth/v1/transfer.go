package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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

var transferNotifyRoutingKey string

func GetTransfer(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	mainDb := utils.MainDbBegin()
	defer mainDb.DbRollback()
	var service Service
	if mainDb.Where("app_key = ?", params["app_key"]).First(&service).RecordNotFound() {
		return utils.BuildError("1102")
	}
	var transfer Transfer
	if mainDb.Where("sn = ? AND service_id = ?", params["sn"], service.Id).First(&transfer).RecordNotFound() {
		return utils.BuildError("1113")
	}
	response := utils.SuccessResponse
	response.Body = transfer
	return context.JSON(http.StatusOK, response)
}

func PostTransfersCreate(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	user := context.Get("current_user").(User)
	service := context.Get("current_service").(Service)
	notifyUrlStruct, _ := url.Parse(params["notify_url"])
	if !service.ValidataHost(notifyUrlStruct.Host) {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid notify_url."})
	}
	if params["currency"] == "" {
		return utils.BuildError("1107")
	}

	amount, _ := decimal.NewFromString(params["amount"])
	amount = amount.Truncate(8)
	if amount.LessThanOrEqual(decimal.Zero) {
		return utils.BuildError("1105")
	}

	mainDb := utils.MainDbBegin()
	defer mainDb.DbRollback()

	var currency Currency
	if mainDb.Where("code = ?", params["currency"]).First(&currency).RecordNotFound() {
		return utils.BuildError("1107")
	}

	transfer := Transfer{
		From:       user.Id,
		To:         service.UserId,
		Amount:     amount,
		CurrencyId: currency.Id,
		ServiceId:  service.Id,
		Sn:         params["sn"],
		State:      "pending",
	}
	if !mainDb.Where("service_id = ?", service.Id).Where("sn = ?", params["sn"]).Find(&Transfer{}).RecordNotFound() {
		return utils.BuildError("1114")
	}
	mainDb.DbRollback()
	err := transferToChangeAccount(&transfer, 9)
	if err == nil {
		transfer.InitializeTimestamp()
		sendMessageToNotifyUrl(params["notify_url"], &map[string]string{
			"id":           strconv.Itoa(transfer.Id),
			"from":         strconv.Itoa(transfer.From),
			"to":           strconv.Itoa(transfer.To),
			"amount":       transfer.Amount.String(),
			"currency":     currency.Code,
			"service_id":   strconv.Itoa(service.Id),
			"sn":           params["sn"],
			"notify_times": "1",
			"timestamp":    strconv.Itoa(int(time.Now().Unix())),
		})
		response := utils.SuccessResponse
		return context.JSON(http.StatusOK, response)
	} else {
		return utils.BuildError("1020")
	}
}

func PostTransfersBack(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	user := context.Get("current_user").(User)
	service := context.Get("current_service").(Service)
	notifyUrlStruct, _ := url.Parse(params["notify_url"])
	if !service.ValidataHost(notifyUrlStruct.Host) {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid notify_url."})
	}
	amount, _ := decimal.NewFromString(params["amount"])
	amount = amount.Truncate(8)
	if amount.LessThanOrEqual(decimal.Zero) {
		return utils.BuildError("1105")
	}

	mainDb := utils.MainDbBegin()
	defer mainDb.DbRollback()
	var currency Currency
	if mainDb.Where("code = ?", params["currency"]).First(&currency).RecordNotFound() {
		return utils.BuildError("1107")
	}

	transfer := Transfer{
		From:       service.UserId,
		To:         user.Id,
		Amount:     amount,
		CurrencyId: currency.Id,
		ServiceId:  service.Id,
		Sn:         params["sn"],
		State:      "pending",
	}
	if !mainDb.Where("service_id = ?", service.Id).Where("sn = ?", params["sn"]).Find(&Transfer{}).RecordNotFound() {
		return utils.BuildError("1114")
	}
	mainDb.DbRollback()
	err := transferToChangeAccount(&transfer, 3)
	if err == nil {
		transfer.InitializeTimestamp()
		sendMessageToNotifyUrl(params["notify_url"], &map[string]string{
			"id":           strconv.Itoa(transfer.Id),
			"from":         strconv.Itoa(transfer.From),
			"to":           strconv.Itoa(transfer.To),
			"amount":       transfer.Amount.String(),
			"currency":     currency.Code,
			"service_id":   strconv.Itoa(service.Id),
			"sn":           params["sn"],
			"notify_times": "1",
			"timestamp":    strconv.Itoa(int(time.Now().Unix())),
		})
		response := utils.SuccessResponse
		return context.JSON(http.StatusOK, response)
	} else {
		return utils.BuildError("1020")
	}
}

func transferToChangeAccount(transfer *Transfer, times int) error {
	mainDb := utils.MainDbBegin()
	defer mainDb.DbRollback()
	var fromAccount, toAccount Account
	if mainDb.Where("user_id = ?", (*transfer).From).Where("currency = ?", (*transfer).CurrencyId).First(&fromAccount).RecordNotFound() {
		fromAccount.UserId = transfer.From
		fromAccount.CurrencyId = transfer.CurrencyId
		fromAccount.Balance = decimal.Zero
		fromAccount.Locked = decimal.Zero
		mainDb.Save(&fromAccount)
		mainDb.DbCommit()
		return utils.BuildError("1103")
	}
	if mainDb.Where("user_id = ?", (*transfer).To).Where("currency = ?", (*transfer).CurrencyId).First(&toAccount).RecordNotFound() {
		toAccount.UserId = transfer.To
		toAccount.CurrencyId = transfer.CurrencyId
		toAccount.Balance = decimal.Zero
		toAccount.Locked = decimal.Zero
		mainDb.Save(&toAccount)
	}

	if fromAccount.Balance.Sub(transfer.Amount).IsNegative() {
		return utils.BuildError("1103")
	}
	transfer.State = "done"
	mainDb.Create(transfer)
	if transfer.Id == 0 {
		return utils.BuildError("1112")
	}
	err1 := fromAccount.SubFunds(mainDb, (*transfer).Amount, decimal.Zero, TRANSFER_BACK, (*transfer).Id, "Transfer")
	err2 := toAccount.PlusFunds(mainDb, (*transfer).Amount, decimal.Zero, TRANSFER, (*transfer).Id, "Transfer")

	if err1 == nil && err2 == nil {
		mainDb.DbCommit()
		return nil
	}
	mainDb.DbRollback()
	if times > 0 {
		(*transfer).Id = 0
		return transferToChangeAccount(transfer, times-1)
	}
	return utils.BuildError("1104")
}

func sendMessageToNotifyUrl(notify_url string, message *map[string]string) {
	if transferNotifyRoutingKey == "" {
		for _, worker := range sneakerWorkers.AllWorkers {
			if worker.Name == "TransferNotifyWorker" {
				transferNotifyRoutingKey = worker.RoutingKey
			}
		}
	}
	params := map[string]string{"notify_url": notify_url}
	for k, v := range *message {
		params[k] = v
	}
	b, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = initializers.PublishMessageWithRouteKey(
		initializers.AmqpGlobalConfig.Exchange["default"]["key"],
		transferNotifyRoutingKey,
		"text/plain",
		&b,
		amqp.Table{},
		amqp.Persistent,
	)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
