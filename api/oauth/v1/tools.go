package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/GeeTeam/GtGoSdk"
	"github.com/labstack/echo"
	"github.com/streadway/amqp"

	"cherry/initializers"
	. "cherry/models"
	"cherry/utils"
	"cherry/workers/sneakerWorkers"
)

func GetGeetest(context echo.Context) error {
	geetest_token := initializers.GeetestConfig.GenerateToken(context.RealIP())
	gt := GtGoSdk.GeetestLib(initializers.GeetestConfig.PrivateKey, initializers.GeetestConfig.CaptchaID)
	gt.PreProcess(geetest_token)
	responseStr := gt.GetResponseStr()
	type GeetestAttrs struct {
		Challenge string `json:"challenge"`
		Gt        string `json:"gt"`
		Success   int    `json:"success"`
		UserId    string `json:"user_id"`
	}
	var geetestAttrs GeetestAttrs
	json.Unmarshal([]byte(strings.Replace(responseStr, `\`, "", -1)), &geetestAttrs)

	response := utils.SuccessResponse
	response.Body = geetestAttrs
	return context.JSON(http.StatusOK, response)
}

func PostEmail(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	db := utils.MainDbBegin()
	defer db.DbRollback()
	var identity Identity
	if db.Where("source = ?", "Email").Where("user_id = ?", params["uid"]).First(&identity).RecordNotFound() {
		return utils.BuildError("1021")
	}
	pushMessageToEmail(&map[string]string{
		"email":   identity.Symbol,
		"title":   params["subject"],
		"content": params["content"],
	})
	response := utils.SuccessResponse
	return context.JSON(http.StatusOK, response)
}

func PostSms(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	db := utils.MainDbBegin()
	defer db.DbRollback()
	var identity Identity
	if db.Where("source = ?", "Phone").Where("user_id = ?", params["uid"]).First(&identity).RecordNotFound() {
		return utils.BuildError("1021")
	}
	pushMessageToSms(&map[string]string{
		"phone":   identity.Symbol,
		"type":    "enterprise",
		"content": params["content"],
	})
	response := utils.SuccessResponse
	return context.JSON(http.StatusOK, response)
}

var emailRoutingKey, smsRoutingKey string

func pushMessageToEmail(payload *map[string]string) {
	if emailRoutingKey == "" {
		for _, worker := range sneakerWorkers.AllWorkers {
			if worker.Name == "EmailWorker" {
				emailRoutingKey = worker.RoutingKey
			}
		}
	}
	b, err := json.Marshal(*payload)
	if err != nil {
		fmt.Println("{ error:", err, "}")
		panic(err)
	}
	err = initializers.PublishMessageWithRouteKey(
		initializers.AmqpGlobalConfig.Exchange["default"]["key"],
		emailRoutingKey,
		"text/plain",
		&b,
		amqp.Table{},
		amqp.Persistent,
	)
	if err != nil {
		fmt.Println("{ error:", err, "}")
		panic(err)
	}
}

func pushMessageToSms(payload *map[string]string) {
	if smsRoutingKey == "" {
		for _, worker := range sneakerWorkers.AllWorkers {
			if worker.Name == "SmsWorker" {
				smsRoutingKey = worker.RoutingKey
			}
		}
	}
	b, err := json.Marshal(*payload)
	if err != nil {
		fmt.Println("{ error:", err, "}")
		panic(err)
	}
	err = initializers.PublishMessageWithRouteKey(
		initializers.AmqpGlobalConfig.Exchange["default"]["key"],
		smsRoutingKey,
		"text/plain",
		&b,
		amqp.Table{},
		amqp.Persistent,
	)
	if err != nil {
		fmt.Println("{ error:", err, "}")
		panic(err)
	}
}
