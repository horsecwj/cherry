package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"

	envConfig "cherry/config"
)

func SendSms253(phone string, msg string) {
	params := make(map[string]interface{})
	params["account"] = envConfig.CurrentEnv.Sms.Sms253Account
	params["password"] = envConfig.CurrentEnv.Sms.Sms253Password
	params["mobile"] = phone

	params["msg"] = msg
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://intapi.253.com/send/json" // 国际短信接口，如果使用国内接口，需变更API_URL
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
}
