package v1

import (
	"net/http"

	"github.com/labstack/echo"

	. "cherry/orm/db/models"
	"cherry/utils"
)

func UserAccounts(context echo.Context) error {
	user := context.Get("current_user").(User)
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var accounts []Account
	mainDB.Where("user_id = ?", user.Id).Find(&accounts)

	response := utils.SuccessResponse
	response.Body = accounts
	return context.JSON(http.StatusOK, response)
}
