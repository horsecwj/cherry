package v1

import (
	"net/http"

	"github.com/labstack/echo"

	. "cherry/models"
	"cherry/utils"
)

func Currencies(context echo.Context) error {
	response := utils.SuccessResponse
	response.Body = AllCurrencies
	return context.JSON(http.StatusOK, response)
}
