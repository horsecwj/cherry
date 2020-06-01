package v1

import (
	"net/http"

	"github.com/labstack/echo"

	. "cherry/models"
	"cherry/utils"
)

func ServiceList(context echo.Context) error {
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var services []Service
	mainDB.Where("inside = true").Find(&services)
	response := utils.SuccessResponse
	response.Body = services
	return context.JSON(http.StatusOK, response)
}
