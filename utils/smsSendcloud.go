package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	envConfig "cherry/config"
)

// text
func SendSmsSendcloud(phone string, payloadJson *[]byte, templateId string) {
	payload := string([]byte(*payloadJson))
	msgType := "0"
	// purge phone number
	if !CheckIsChinesePhone(phone) {
		msgType = "2"
		phone = "00" + phone
	} else {
		phone = strings.TrimLeft(phone, "86")
	}
	params := map[string]string{"smsUser": envConfig.CurrentEnv.Sms.SmsUser, "templateId": templateId, "msgType": msgType, "phone": phone, "vars": payload}
	apiUrl := envConfig.CurrentEnv.Sms.ApiUrl
	HandleRequest(params, apiUrl)
}

// voice code
func SendSmsSendcloudVoice(phone string, payloadJson *[]byte) {
	var payload struct {
		Code string
	}
	json.Unmarshal([]byte(*payloadJson), &payload)
	// purge phone number
	if !CheckIsChinesePhone(phone) {
		fmt.Printf("voice code only for chinese phone.")
	} else {
		phone = strings.TrimLeft(phone, "86")
	}
	params := map[string]string{"smsUser": envConfig.CurrentEnv.Sms.SmsUser, "phone": phone, "code": payload.Code}
	apiUrl := envConfig.CurrentEnv.Sms.VoiceApiUrl
	HandleRequest(params, apiUrl)
}

// Signature
func SignatureSmsSendcloud(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	paramStr := ""
	for _, k := range keys {
		str := k + "=" + params[k] + "&"
		paramStr += str
	}
	paramStr = envConfig.CurrentEnv.Sms.SmsKey + "&" + paramStr + envConfig.CurrentEnv.Sms.SmsKey
	signature := Md5Sendcloud(paramStr)
	return strings.ToUpper(signature)
}

func Md5Sendcloud(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func CheckIsChinesePhone(phone string) bool {
	return strings.HasPrefix(phone, "86")
}

func HandleRequest(params map[string]string, apiUrl string) {
	signature := SignatureSmsSendcloud(params)
	params["signature"] = signature
	postValues := url.Values{}
	for postKey, PostValue := range params {
		postValues.Set(postKey, PostValue)
	}
	paramsStr := postValues.Encode()
	paramsBytes := []byte(paramsStr)
	postBytesReader := bytes.NewReader(paramsBytes)
	httpReq, _ := http.NewRequest("POST", apiUrl, postBytesReader)
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Printf("http get strUrl=%s response error=%s\n", apiUrl, err.Error())
	}
	defer httpResp.Body.Close()
	body, errReadAll := ioutil.ReadAll(httpResp.Body)
	if errReadAll != nil {
		fmt.Printf("get response for strUrl=%s got error=%s\n", apiUrl, errReadAll.Error())
	}
	fmt.Println(string(body))
}
