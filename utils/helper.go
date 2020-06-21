package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func GetRealIp(context echo.Context) (ip string) {
	ips := context.RealIP()
	rips := strings.Split(ips, ",")
	ip = rips[0]
	return
}

func SetLogAndPid(name string) {
	err := os.Mkdir("logs", 0755)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	file, err := os.OpenFile("logs/"+name+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	log.SetOutput(file)
	err = os.Mkdir("pids", 0755)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	err = ioutil.WriteFile("pids/"+name+".pid", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
