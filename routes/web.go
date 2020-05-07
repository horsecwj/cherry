package routes

import (
	v1 "cherry/api/web/v1"
	"github.com/labstack/echo"
)

func SetWebInterfaces(e *echo.Echo) {

	e.GET("/api/web/currencies", v1.Currencies)

	e.POST("/api/web/grant/create", v1.PostGrantCreate)
	e.GET("/api/web/grant/confirm", v1.GrantConfirm)
	e.POST("/api/web/grant/confirm", v1.GrantConfirm)

	e.GET("/api/web/user/wechat/login", v1.GetUsersWechatLogin)
	e.GET("/api/web/user/wechat/callback", v1.UsersWechatCallback)

}
