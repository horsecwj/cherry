package initializers

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"time"

	"github.com/GeeTeam/GtGoSdk"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"

	"cherry/initializers/locale"
	. "cherry/models"
	"cherry/utils"
)

func LimitTrafficWithIp(context echo.Context) bool {
	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()
	key := "limit-traffic-with-ip:" + context.Path() + ":" + context.RealIP()
	timesStr, _ := redis.String(dataRedis.Do("GET", key))
	if timesStr == "" {
		dataRedis.Do("SETEX", key, 1, 60)
	} else {
		times, _ := strconv.Atoi(timesStr)
		if times > 10 {
			return false
		} else {
			dataRedis.Do("INCR", key)
		}
	}
	return true
}

func disabledValidate(context echo.Context, user *User) error {
	if (*user).Disabled {
		return utils.BuildError("3047")
	}
	return nil
}

func lowIntensityUserTrust(context echo.Context, user *User) error {
	if user.Activated == false {
		return utils.BuildError("4017")
	}
	return nil
}

func treatLanguage(context echo.Context) {
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
	context.Set("language", language)
}

func validateGeetest(context echo.Context, params *map[string]string) bool {
	if !validateGeetestUserId(context, params) {
		return false
	}
	if (*params)[GtGoSdk.FN_CHALLENGE] == "" || (*params)[GtGoSdk.FN_VALIDATE] == "" || (*params)[GtGoSdk.FN_SECCODE] == "" || (*params)[GtGoSdk.GT_STATUS_SESSION_KEY] == "" || (*params)["user_id"] == "" {
		return false
	}
	var result bool
	gt := GtGoSdk.GeetestLib(GeetestConfig.PrivateKey, GeetestConfig.CaptchaID)
	challenge := (*params)[GtGoSdk.FN_CHALLENGE]
	validate := (*params)[GtGoSdk.FN_VALIDATE]
	seccode := (*params)[GtGoSdk.FN_SECCODE]
	status, err := strconv.Atoi((*params)[GtGoSdk.GT_STATUS_SESSION_KEY])

	if err != nil {
		return false
	}
	userId := (*params)["user_id"]
	if status == 0 {
		result = gt.FailbackValidate(challenge, validate, seccode)
	} else {
		result = gt.SuccessValidate(challenge, validate, seccode, userId)
	}
	return result
}

func trustGeetest(context echo.Context, params *map[string]string) bool {
	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()
	key := "limit-traffic-with-ip-" + context.RealIP()
	timesStr, _ := redis.String(dataRedis.Do("GET", key))
	if timesStr == "" {
		return true
	} else {
		times, _ := strconv.Atoi(timesStr)
		if times < 3 {
			return true
		}
	}

	if (*params)[GtGoSdk.FN_CHALLENGE] == "" || (*params)[GtGoSdk.FN_VALIDATE] == "" || (*params)[GtGoSdk.FN_SECCODE] == "" || (*params)[GtGoSdk.GT_STATUS_SESSION_KEY] == "" || (*params)["user_id"] == "" {
		return false
	}
	var result bool
	gt := GtGoSdk.GeetestLib(GeetestConfig.PrivateKey, GeetestConfig.CaptchaID)
	challenge := (*params)[GtGoSdk.FN_CHALLENGE]
	validate := (*params)[GtGoSdk.FN_VALIDATE]
	seccode := (*params)[GtGoSdk.FN_SECCODE]
	status, err := strconv.Atoi((*params)[GtGoSdk.GT_STATUS_SESSION_KEY])

	if err != nil {
		return false
	}
	userId := (*params)["user_id"]
	if status == 0 {
		result = gt.FailbackValidate(challenge, validate, seccode)
	} else {
		result = gt.SuccessValidate(challenge, validate, seccode, userId)
	}
	return result
}

func checkTimestamp(context echo.Context, params *map[string]string) bool {
	timestamp, _ := strconv.Atoi((*params)["timestamp"])
	now := time.Now()
	if int(now.Add(-time.Second*60*5).Unix()) <= timestamp && timestamp <= int(now.Add(time.Second*60*5).Unix()) {
		return true
	}
	return false
}

func validateGeetestUserId(context echo.Context, params *map[string]string) bool {
	decodeBytes, err := base64.StdEncoding.DecodeString((*params)["user_id"])
	if err != nil {
		return false
	}
	matched, _ := regexp.MatchString("-"+context.RealIP(), string(decodeBytes))
	return matched
}

func verifyEmailOrPhoneBeforeRegist(context echo.Context) error {
	dataRedis := utils.GetRedisConn("data")
	defer dataRedis.Close()
	email := context.FormValue("email")
	phone := context.FormValue("phone_number")
	token := context.FormValue("token")
	if token == "" {
		return utils.BuildError("5003")
	}
	if email == "" && phone == "" {
		return utils.BuildError("5004")
	}
	existedToken := ""
	if email == "" {
		existedToken, _ = redis.String(dataRedis.Do("GET", "user_register_verify:phone:"+phone))
	} else {
		existedToken, _ = redis.String(dataRedis.Do("GET", "user_register_verify:email:"+email))
	}
	if existedToken != token {
		return utils.BuildError("5002")
	}
	return nil
}
