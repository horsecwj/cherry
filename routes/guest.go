package routes

import (
	v1 "cherry/api/guest/v1"
	"github.com/labstack/echo"
)

func SetGuestInterfaces(e *echo.Echo) {

	e.GET("/api/guest/currencies", v1.Currencies)

	e.GET("/api/guest/user/accounts", v1.UserAccounts)

	e.GET("/api/guest/user/wechat/login", v1.GetUsersWechatLogin)
	e.GET("/api/guest/user/wechat/callback", v1.UsersWechatCallback)

	e.POST("/api/guest/recharge/wechat", v1.PostRechargeFromWechat)
	e.GET("/api/guest/recharge/wechat/verify", v1.RechargeFromWechatVerify)
	e.POST("/api/guest/recharge/wechat/verify", v1.RechargeFromWechatVerify)
}
