package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	newrelic "github.com/oldfritter/echo-middleware"

	envConfig "cherry/config"
	"cherry/initializers"
	"cherry/models"
	"cherry/routes"
	"cherry/utils"
	"cherry/workers/sneakerWorkers"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	initialize()
	sneakerWorkers.InitWorkers()
	e := echo.New()
	e.File("/web", "public/assets/index.html")
	e.File("/web/*", "public/assets/index.html")
	e.Static("/assets", "public/assets")
	if envConfig.CurrentEnv.Newrelic.AppName != "" && envConfig.CurrentEnv.Newrelic.LicenseKey != "" {
		e.Use(newrelic.NewRelic(envConfig.CurrentEnv.Newrelic.AppName, envConfig.CurrentEnv.Newrelic.LicenseKey))
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(initializers.Auth)
	routes.SetWebInterfaces(e)
	routes.SetOauthInterfaces(e)
	routes.SetWsInterfaces(e)

	e.HTTPErrorHandler = customHTTPErrorHandler
	t := &Template{
		templates: template.Must(template.ParseGlob("public/api/*/*.html")),
	}
	e.Renderer = t
	e.HideBanner = true
	go func() {
		if err := e.Start(":9700"); err != nil {
			log.Println("start close echo")
			time.Sleep(500 * time.Millisecond)
			closeResource()
			log.Println("shutting down the server")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("accepted signal")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Println("shutting down failed, err:" + err.Error())
		e.Logger.Fatal(err)
	} else {
		log.Println("shutting down complete.")
	}
	closeResource()
}

func customHTTPErrorHandler(err error, context echo.Context) {
	language := context.Get("language").(string)
	if response, ok := err.(utils.Response); ok {
		response.Head["msg"] = fmt.Sprint(models.I18n.T(language, "error_code."+response.Head["code"]))
		context.JSON(http.StatusBadRequest, response)
	} else {
		panic(err)
	}
}

func initialize() {
	envConfig.InitEnv()
	utils.InitMainDB()
	utils.InitRedisPools()
	models.AutoMigrations()
	models.InitI18n()
	initializers.InitializeAmqpConfig()
	initializers.LoadInterfaces()
	initializers.LoadCacheData()
	initializers.InitAllPayments()
	utils.SetLogAndPid("api")
}

func closeResource() {
	initializers.CloseAmqpConnection()
	utils.CloseRedisPools()
	utils.CloseMainDB()
}
