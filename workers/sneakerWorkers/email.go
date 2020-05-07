package sneakerWorkers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"reflect"
	"time"
	"unicode"

	. "cherry/models"
	"cherry/utils"
)

func (worker Worker) EmailWorker(payloadJson *[]byte) (err error) {
	start := time.Now().UnixNano()
	var payload MailerPayload
	err = json.Unmarshal([]byte(*payloadJson), &payload)
	if err != nil {
		worker.LogError("parse JSON err:", err)
		return
	}
	worker.LogInfo("..........payload", reflect.ValueOf(&payload))
	reflect.ValueOf(&payload).MethodByName(payload.I18nFuncName()).Call([]reflect.Value{})
	t, err := template.New(payload.Method+".html").Funcs(payload.FuncMap).ParseFiles(
		"public/workers/emailWorker/"+ToCamel(payload.MailerClass)+"/"+payload.Method+".html",
		"public/workers/emailWorker/head.html",
		"public/workers/emailWorker/footer.html",
		"public/workers/emailWorker/content.html",
	)
	if err != nil {
		worker.LogError("parse file err:", err)
		return
	}
	var tpl bytes.Buffer
	if err = t.Execute(&tpl, payload); err != nil {
		worker.LogError((time.Now().UnixNano()-start)/1000000, " ms, payload: ", payload, err)
		return
	}
	if err = utils.SendMail(payload.Email, payload.Subject, tpl.String()); err == nil {
		worker.LogError((time.Now().UnixNano()-start)/1000000, " ms, payload: ", payload, err)
		return
	}
	worker.LogInfo((time.Now().UnixNano()-start)/1000000, " ms, payload: ", payload, err)
	return
}

func ToCamel(in string) string {
	for i, v := range in {
		return string(unicode.ToLower(v)) + in[i+1:]
	}
	return ""
}
