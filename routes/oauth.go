package routes

import (
	v1 "cherry/api/oauth/v1"
	"github.com/labstack/echo"
)

func SetOauthInterfaces(e *echo.Echo) {

	e.GET("/oauth/authorize", v1.GetOauthAuthorization)
	e.GET("/oauth/login", v1.OauthLogin)
	e.GET("/oauth/access_token", v1.GetOauthAccessToken)
	e.GET("/oauth/user_info", v1.GetOauthUserInfo)

	e.GET("/oauth/transfer", v1.GetTransfer)
	e.POST("/oauth/transfers/create", v1.PostTransfersCreate)
	e.POST("/oauth/transfers/back", v1.PostTransfersBack)

	e.POST("/oauth/email", v1.PostEmail)
	e.POST("/oauth/sms", v1.PostSms)
	e.GET("/oauth/geetest", v1.GetGeetest)

	e.GET("/oauth/services", v1.ServiceList)
}
