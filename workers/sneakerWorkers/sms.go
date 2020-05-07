package sneakerWorkers

import (
	"encoding/json"

	envConfig "cherry/config"
	"cherry/utils"
)

func (worker Worker) SmsWorker(payloadJson *[]byte) (queueName string, message []byte) {
	var payload struct {
		Phone       string //必须参数
		Type        string //必须参数
		Content     string //必须参数
		TextOrVoice string `json:"text_or_voice"` // default `sendcloud`
		Platform    string // default `sendcloud`
	}
	json.Unmarshal([]byte(*payloadJson), &payload)

	if payload.TextOrVoice == "voice" && payload.Type == "validate_code" && payload.Platform != "253" { // only validation code, platform: sendcloud
		utils.SendSmsSendcloudVoice(payload.Phone, payloadJson)
	} else {
		// text message
		if payload.Type == "validate_code" {
			ValidateContent(payload.Phone, payloadJson, payload.Platform, payload.Content)
		}
	}
	return
}

func ValidateContent(phone string, payloadJson *[]byte, platform string, code string) {
	if platform == "253" {
		utils.SendSms253(phone, "【Rfinex】您的验证码为"+code)
	} else {
		utils.SendSmsSendcloud(phone, payloadJson, envConfig.CurrentEnv.Sms.ValidateCode)
	}
}
