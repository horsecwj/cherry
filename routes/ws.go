package routes

import (
	v1 "cherry/api/ws/v1"
	"github.com/labstack/echo"
)

func SetWsInterfaces(e *echo.Echo) {

	e.GET("/ws/accounts", v1.Accounts)

}
