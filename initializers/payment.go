package initializers

import (
	"github.com/objcoding/wxpay"

	"cherry/config"
)

var (
	WxpayAccount *wxpay.Account
	WxpayClient  *wxpay.Client
)

func InitAllPayments() {
	InitWcpay()
}

func InitWcpay() {
	WxpayAccount = wxpay.NewAccount(
		config.CurrentEnv.Wechat["appid"],
		config.CurrentEnv.Wechat["mch_id"],
		config.CurrentEnv.Wechat["key"],
		false,
	)
	WxpayClient = wxpay.NewClient(WxpayAccount)
	WxpayAccount.SetCertData("config/pems/apiclient_cert.p12")
	// WxpayClient.setAccount(WxpayAccount)
	WxpayClient.SetHttpConnectTimeoutMs(2000)
	WxpayClient.SetHttpReadTimeoutMs(1000)
}
