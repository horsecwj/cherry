package initializers

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"

	. "cherry/models"
	"cherry/utils"
)

type ApiInterface struct {
	Method             string `yaml:"method"`
	Path               string `yaml:"path"`
	Auth               bool   `yaml:"auth"`
	Sign               bool   `yaml:"sign"`
	CheckFormat        bool   `yaml:"check_format"`
	CheckTimestamp     bool   `yaml:"check_timestamp"`
	LimitTrafficWithIp bool   `yaml:"limit_traffic_with_ip"`
}

var GlobalApiInterfaces []ApiInterface

func LoadInterfaces() {
	files, err := ioutil.ReadDir("config/interfaces/")
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, f := range files {
		if matched, err := regexp.MatchString(".yml$", f.Name()); matched && err == nil {
			path_str, _ := filepath.Abs("config/interfaces/" + f.Name())
			content, err := ioutil.ReadFile(path_str)
			if err != nil {
				log.Fatal(err)
				return
			}
			var interfaces []ApiInterface
			err = yaml.Unmarshal(content, &interfaces)
			if err != nil {
				log.Fatal(err)
			}
			GlobalApiInterfaces = append(GlobalApiInterfaces, interfaces...)
		}
	}
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		treatLanguage(context)
		params := make(map[string]string)
		for k, v := range context.QueryParams() {
			params[k] = v[0]
		}
		values, _ := context.FormParams()
		for k, v := range values {
			params[k] = v[0]
		}
		context.Set("params", params)
		var currentApiInterface ApiInterface
		for _, apiInterface := range GlobalApiInterfaces {
			if context.Path() == apiInterface.Path && context.Request().Method == apiInterface.Method {
				currentApiInterface = apiInterface
				if currentApiInterface.LimitTrafficWithIp && LimitTrafficWithIp(context) != true {
					return utils.BuildError("3043")
				}
				if apiInterface.Auth != true {
					return next(context)
				}
			}
		}
		if currentApiInterface.Path == "" {
			return utils.BuildError("1025")
		}
		if context.Request().Header.Get("Authorization") == "" && context.QueryParam("token") == "" {
			return utils.BuildError("2016")
		}
		if currentApiInterface.CheckTimestamp && checkTimestamp(context, &params) == false {
			return utils.BuildError("3050")
		}

		db := utils.MainDbBegin()
		defer db.DbRollback()

		var user User
		var err error
		var matchedOauth bool
		if matchedOauth, err = regexp.MatchString("^/oauth/", context.Path()); matchedOauth {
			var token Token
			user, token, err = oauthAuth(context, &params, db)
			context.Set("current_token", token.Token)
		} else {
			var token Token
			user, token, err = normalAuth(context, &params, db)
			context.Set("current_token", token.Token)
		}
		if err != nil {
			log.Println("err: ", err)
			return err
		}

		db.DbCommit()
		context.Set("current_user", user)
		return next(context)
	}
}

func oauthAuth(context echo.Context, params *map[string]string, db *utils.GormDB) (user User, token Token, err error) {
	tokenStr := (*params)["access_token"]
	if db.Where("`type` = ?", "AccessToken").Where("token = ? AND ? < expire_at", tokenStr, time.Now()).First(&token).RecordNotFound() {
		return user, token, utils.BuildError("4016")
	}
	if db.Where("id = ?", token.UserId).First(&user).RecordNotFound() {
		return user, token, utils.BuildError("4016")
	}
	return
}

func normalAuth(context echo.Context, params *map[string]string, db *utils.GormDB) (user User, token Token, err error) {
	tokenStr := (*params)["token"]
	if tokenStr == "" {
		tokenStr = context.Request().Header.Get("Authorization")
	}
	if db.Where("`type` = ?", "AccessToken").Where("token = ? AND ? < expire_at", tokenStr, time.Now()).First(&token).RecordNotFound() {
		return user, token, utils.BuildError("4016")
	}
	if db.Where("id = ?", token.UserId).First(&user).RecordNotFound() {
		return user, token, utils.BuildError("4016")
	}
	return
}
