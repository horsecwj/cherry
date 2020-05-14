package v1

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"

	"cherry/initializers/locale"
	. "cherry/models"
	"cherry/utils"
)

type DataStruct struct {
	AppKey             string
	CallbackUrl        string
	Notice             string
	LoginWitihRfinex   string
	RfinexAccount      string
	InputPassword      string
	Authorize          string
	AuthorizeAndAgree  string
	ProtocolFromRfinex string
	RegisterNow        string
	AlreadyLogin       string
	OtherAccountLogin  string
	NeedGeetest        bool
}

func GetOauthAuthorization(context echo.Context) error {
	var language string
	var lqs []locale.LangQ
	if context.QueryParam("lang") != "" {
		lqs = locale.ParseAcceptLanguage(context.QueryParam("lang"))
	} else {
		lqs = locale.ParseAcceptLanguage(context.Request().Header.Get("Accept-Language"))
	}
	if lqs[0].Lang == "en" {
		language = "en"
	} else if lqs[0].Lang == "ja" {
		language = "ja"
	} else if lqs[0].Lang == "ko" {
		language = "ko"
	} else {
		language = "zh-CN"
	}
	callBackUrl := context.QueryParam("callback_url")
	appKey := context.QueryParam("app_key")
	db := utils.MainDbBegin()
	defer db.DbRollback()
	var service Service
	if db.Where("app_key = ?", appKey).Where("Aasm_State = ?", "verified").First(&service).RecordNotFound() {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid app_key."})
	}
	urlStruct, _ := url.Parse(callBackUrl)
	if !service.ValidataHost(urlStruct.Host) {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid callback_url."})
	}
	data := DataStruct{
		AppKey:             service.AppKey,
		CallbackUrl:        callBackUrl,
		LoginWitihRfinex:   fmt.Sprint(I18n.T(language, "oauth.notice.login_witih_rfinex")),
		RfinexAccount:      fmt.Sprint(I18n.T(language, "oauth.notice.rfinex_account")),
		InputPassword:      fmt.Sprint(I18n.T(language, "oauth.notice.input_password")),
		Authorize:          fmt.Sprint(I18n.T(language, "oauth.notice.authorize")),
		AuthorizeAndAgree:  fmt.Sprint(I18n.T(language, "oauth.notice.authorize_and_agree")),
		ProtocolFromRfinex: fmt.Sprint(I18n.T(language, "oauth.notice.protocol_from_rfinex")),
		RegisterNow:        fmt.Sprint(I18n.T(language, "oauth.notice.register_now")),
		AlreadyLogin:       fmt.Sprint(I18n.T(language, "oauth.notice.already_login")),
		OtherAccountLogin:  fmt.Sprint(I18n.T(language, "oauth.notice.other_account_login")),
	}
	return context.Render(http.StatusOK, "oauth/login", data)
}

func OauthLogin(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	callBackUrl := context.QueryParam("callback_url")
	appKey := context.QueryParam("app_key")
	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()
	db := utils.MainDbBegin()
	defer db.DbRollback()
	var service Service
	if db.Where("app_key = ?", appKey).Where("Aasm_State = ?", "verified").First(&service).RecordNotFound() {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid app_key."})
	}
	callbackUrlStruct, _ := url.Parse(callBackUrl)
	if !service.ValidataHost(callbackUrlStruct.Host) {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid callback_url."})
	}

	var identity Identity
	if db.Where("`source` = ?", params["source"]).Where("symbol").First(&identity).RecordNotFound() {
		loginError(context)
	}
	var user User
	if db.Where("id = ?", identity.UserId).First(&user).RecordNotFound() {
		loginError(context)
	}
	user.Password = context.QueryParam("password")
	if !user.CompareHashAndPassword() {
		loginError(context)
	}

	authorizationhCode := AuthorizationhCode{UserId: user.Id, ServiceId: service.Id}
	authorizationhCode.InitCode()
	dataRedis.Do("SETEX", authorizationhCode.Key(), 600, authorizationhCode.UserId)
	values := callbackUrlStruct.Query()
	values.Add("code", authorizationhCode.Code)
	callbackUrlStruct.RawQuery = values.Encode()
	return context.Redirect(http.StatusPermanentRedirect, callbackUrlStruct.String())
}

func GetOauthAccessToken(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	appKey := context.QueryParam("app_key")
	if appKey == "" {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "Need app_key."})
	}
	db := utils.MainDbBegin()
	defer db.DbRollback()
	var service Service
	if db.Where("aasm_state = ? AND app_key = ?", "verified", appKey).Find(&service).RecordNotFound() {
		return context.JSON(http.StatusForbidden, map[string]string{"message": "app 不存在."})
	}
	if params["grant_type"] == "authorization_code" {
		if params["code"] == "" {
			return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid code."})
		}
		authorizationhCode := AuthorizationhCode{Code: params["code"], ServiceId: service.Id}
		dataRedis := utils.GetRedisConn("data")
		defer dataRedis.Close()
		user_id_str, _ := redis.String(dataRedis.Do("GET", authorizationhCode.Key()))
		if user_id_str == "" {
			return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid code."})
		}
		user_id, _ := strconv.Atoi(user_id_str)
		var token Token
		token.InitializeAccessToken()
		token.UserId = user_id
		db.Save(&token)
		tokenAndApp := TokenAndApp{ServiceId: service.Id, TokenId: token.Id}
		db.Save(&tokenAndApp)
		db.DbCommit()
		response := utils.SuccessResponse
		response.Body = map[string]string{"uid": user_id_str, "token": token.Token, "expire_at": token.ExpireAt.Format("2006-01-02 15:04:05"), "type": "access_token"}
		return context.JSON(http.StatusOK, response)
	} else if params["grant_type"] == "login" {
		if service.Inside == false {
			return context.JSON(http.StatusForbidden, map[string]string{"message": "此app不支持grant_type为login的登录"})
		}
		var identity Identity
		var user User
		var token Token
		token.UserId = user.Id
		token.InitializeAccessToken()

		if db.Where("`source` = ?", params["source"]).Where("symbol").First(&identity).RecordNotFound() {
			return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid source or symbol."})
		}
		if db.Where("id = ?", identity.UserId).First(&user).RecordNotFound() {
			return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid source or symbol."})
		}
		user.Password = context.QueryParam("password")
		if !user.CompareHashAndPassword() {
			return context.JSON(http.StatusForbidden, map[string]string{"message": "Invalid source or symbol."})
		}

		token.UserId = user.Id
		db.Save(&token)
		tokenAndApp := TokenAndApp{ServiceId: service.Id, TokenId: token.Id}
		db.Save(&tokenAndApp)
		db.DbCommit()
		response := utils.SuccessResponse
		response.Body = map[string]string{"uid": strconv.Itoa(user.Id), "token": token.Token, "expire_at": token.ExpireAt.Format("2006-01-02 15:04:05"), "type": "login_token"}
		return context.JSON(http.StatusOK, response)
	}
	return context.JSON(http.StatusForbidden, map[string]string{"message": "Wrong grant_type."})
}

func GetOauthUserInfo(context echo.Context) error {
	user := context.Get("current_user").(User)
	db := utils.MainDbBegin()
	defer db.DbRollback()
	var app Service
	if db.Where("app_key = ?", context.QueryParam("app_key")).First(&app).RecordNotFound() {
		return utils.BuildError("1102")
	}
	response := utils.SuccessResponse
	var smsTwoFactor TwoFactor
	if !db.Where("type = ? AND activated = ?", "TwoFactor::Sms", true).Where("user_id = ?", user.Id).Find(&smsTwoFactor).RecordNotFound() {
		user.SmsValidated = true
	}
	response.Body = OauthUserInfo{
		Uid:          user.Id,
		Sn:           user.Sn,
		Name:         user.Nickname,
		Ancestry:     user.Ancestry,
		Activated:    user.Activated,
		SmsValidated: user.SmsValidated,
	}
	return context.JSON(http.StatusOK, response)
}

func loginError(context echo.Context) error {
	params := context.Get("params").(map[string]string)
	log.Println("{ error:", "登录密码错误,source:", params["source"], "账号: ", params["symbol"], "提交的密码:", "-->"+params["password"]+"<--", "}")
	return utils.BuildError("1101")
}
