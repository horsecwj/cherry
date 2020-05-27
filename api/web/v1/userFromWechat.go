package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo"

	envConfig "cherry/config"
	. "cherry/models"
	"cherry/utils"
)

func GetUsersWechatLogin(context echo.Context) error {
	urlStruct, _ := url.Parse(envConfig.CurrentEnv.Wechat["authorize"])
	values := urlStruct.Query()
	values.Add("appid", envConfig.CurrentEnv.Wechat["app_id"])
	values.Add("scope", envConfig.CurrentEnv.Wechat["snsapi_login"])
	values.Add("redirect_uri", "https://"+context.Request().Host+envConfig.CurrentEnv.Wechat["redirect_url_path"]+"?source=wechat")
	urlStruct.RawQuery = values.Encode()
	return context.Redirect(http.StatusPermanentRedirect, urlStruct.String())
}

func UsersWechatCallback(context echo.Context) error {
	var user User
	if context.Get("current_user") != nil {
		user = context.Get("current_user").(User)
	}
	urlStruct, _ := url.Parse(envConfig.CurrentEnv.Wechat["access_token"])
	values := urlStruct.Query()
	values.Add("appid", envConfig.CurrentEnv.Wechat["app_id"])
	values.Add("secret", envConfig.CurrentEnv.Wechat["secret"])
	values.Add("code", context.QueryParam("code"))
	values.Add("grant_type", "authorization_code")
	urlStruct.RawQuery = values.Encode()
	resp, err := http.Get(urlStruct.String())
	if err != nil {
		return utils.BuildError("1028")
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return utils.BuildError("1028")
	}
	var result struct {
		Head map[string]string
		Body map[string]string
	}
	json.Unmarshal(responseBody, &result)
	if result.Head["code"] != "1000" {
		return utils.BuildError("1028")
	}
	accessToken, openid := result.Body["access_token"], result.Body["openid"]

	urlStruct, _ = url.Parse(envConfig.CurrentEnv.Wechat["user_info"])
	values = urlStruct.Query()
	values.Add("openid", openid)
	values.Add("access_token", accessToken)
	urlStruct.RawQuery = values.Encode()
	resp, err = http.Get(urlStruct.String())
	if err != nil {
		return utils.BuildError("1028")
	}
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return utils.BuildError("1028")
	}
	json.Unmarshal(responseBody, &result)
	if result.Head["code"] != "1000" {
		return utils.BuildError("1028")
	}
	infoBody := result.Body
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var identity Identity
	if mainDB.Where("source = ?", "Wechat").Where("symbol = ?", openid).First(&identity).RecordNotFound() {
		identity = Identity{
			UserId:      user.Id,
			Source:      "Wechat",
			Symbol:      openid,
			AccessToken: accessToken,
			ExpiredAt:   time.Now().Add(time.Hour * 24 * 7),
		}
	} else {
		identity.AccessToken = accessToken
		if user.Id > 0 {
			identity.UserId = user.Id
		} else {
			mainDB.Where("id = ?", identity.UserId).First(&user)
		}
		identity.ExpiredAt = time.Now().Add(time.Hour * 24 * 7)
	}
	user.Nickname = infoBody["display_name"]
	if user.Nickname == "" {
		user.Nickname = "Wechat" + openid
	}
	if user.Id == 0 {
		mainDB.Save(&user)
		identity.UserId = user.Id
	}
	mainDB.Save(&identity)
	var token Token
	token.UserId = user.Id
	token.InitializeAccessToken()
	mainDB.Create(&token)
	mainDB.Where("user_id = ?", user.Id).Where("expire_at > ?", time.Now()).First(&user.Tokens)
	mainDB.DbCommit()
	response := utils.SuccessResponse
	response.Body = user
	return context.JSON(http.StatusOK, response)
}
