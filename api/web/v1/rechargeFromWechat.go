package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/objcoding/wxpay"
	"github.com/shopspring/decimal"

	envConfig "cherry/config"
	"cherry/initializers"
	. "cherry/orm/db/models"
	"cherry/utils"
)

// params: symbol,amount,id(可选)
func PostRechargeFromWechat(context echo.Context) (err error) {
	params := context.Get("params").(map[string]string)
	user := context.Get("current_user").(User)
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var currency Currency
	if mainDB.Where("symbol = ?", strings.ToLower(params["symbol"])).First(&currency).RecordNotFound() {
		return utils.BuildError("1107")
	}
	var recharge Recharge
	if params["id"] != "" {
		mainDB.FirstOrInit(&recharge, map[string]interface{}{"id": params["id"]})
	}
	recharge.Title = "Wechat:Recharge:" + recharge.Sn + "; Amount:" + params["amount"]
	recharge.Source = "Wechat"
	recharge.UserId = user.Id
	recharge.Amount, err = decimal.NewFromString(params["amount"])
	if err != nil || recharge.Amount.LessThanOrEqual(decimal.Zero) {
		return utils.BuildError("1105")
	}
	mainDB.Save(&recharge)
	mainDB.DbCommit()
	wxpayParams := make(wxpay.Params)
	wxpayParams.SetString("body", recharge.Title).
		SetString("out_trade_no", recharge.Sn).
		SetString("total_fee", params["amount"]).
		SetString("spbill_create_ip", context.RealIP()).
		SetString("notify_url", "https://"+context.Request().Host+envConfig.CurrentEnv.Wechat["notify_url_path"]).
		SetString("trade_type", params["trade_type"])
	p, err := initializers.WxpayClient.UnifiedOrder(wxpayParams)
	if err != nil {
		log.Println(err)
	}
	response := utils.SuccessResponse
	response.Body = p
	return context.JSON(http.StatusOK, response)
}

func RechargeFromWechatVerify(context echo.Context) (err error) {
	response := make(wxpay.Params)
	var b []byte
	context.Request().Body.Read(b)
	params := wxpay.XmlToMap(string(b))
	if initializers.WxpayClient.ValidSign(params) {
		mainDB := utils.MainDbBegin()
		defer mainDB.DbRollback()
		var recharge Recharge
		if mainDB.Where("`source` = ?", "Wechat").Where("out_trade_no = ?", params.GetString("out_trade_no")).First(&recharge).RecordNotFound() {
			response.SetString("return_code", "FAIL")
			response.SetString("return_msg", "未找到此数据。")
			return context.XML(http.StatusForbidden, wxpay.MapToXml(response))
		}
		recharge.OutTradeNo = params.GetString("out_trade_no")
		totalFee, _ := decimal.NewFromString(params["total_fee"])
		if !totalFee.Equal(recharge.Amount) {
			response.SetString("return_code", "FAIL")
			response.SetString("return_msg", "金额不一致。")
			return context.XML(http.StatusForbidden, wxpay.MapToXml(response))
		}
		recharge.State = Done
		var currency Currency
		if mainDB.Where("symbol = ?", strings.ToLower(params.GetString("fee_type"))).First(&currency).RecordNotFound() {
			response.SetString("return_code", "FAIL")
			response.SetString("return_msg", "不支持此币种。")
			return context.XML(http.StatusForbidden, wxpay.MapToXml(response))
		}
		var account Account
		mainDB.FirstOrCreate(&account, map[string]interface{}{"user_id": recharge.UserId, "currency_id": currency.Id})
		err = account.PlusFunds(mainDB, recharge.Amount, decimal.Zero, RECHARGE, recharge.Id, "Recharge")
		if err != nil {
			response.SetString("return_code", "FAIL")
			response.SetString("return_msg", "增加余额出错。")
			return context.XML(http.StatusForbidden, wxpay.MapToXml(response))
		}
		mainDB.DbCommit()
		response.SetString("return_code", "SUCCESS")
		return context.XML(http.StatusOK, wxpay.MapToXml(response))
	}
	response.SetString("return_code", "FAIL")
	response.SetString("return_msg", "签名验证失败。")
	return context.XML(http.StatusForbidden, wxpay.MapToXml(response))
}
