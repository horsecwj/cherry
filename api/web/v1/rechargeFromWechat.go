package v1

// import (
//   "encoding/json"
//   "io/ioutil"
//   "net/http"
//   "net/url"
//   "time"
//
//   "github.com/labstack/echo"
//
//   envConfig "cherry/config"
//   . "cherry/models"
//   "cherry/utils"
// )
//
// func GetUsersWechatLogin(context echo.Context) error {
//   urlStruct, _ := url.Parse(envConfig.CurrentEnv.Wechat["authorize"])
//   values := urlStruct.Query()
//   values.Add("appid", envConfig.CurrentEnv.Wechat["app_id"])
//   values.Add("scope", envConfig.CurrentEnv.Wechat["snsapi_login"])
//   values.Add("redirect_uri", "https://"+context.Request().Host+envConfig.CurrentEnv.Wechat["callback_url"]+"?source=wechat")
//   urlStruct.RawQuery = values.Encode()
//   return context.Redirect(http.StatusPermanentRedirect, urlStruct.String())
// }
