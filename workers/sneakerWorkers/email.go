package sneakerWorkers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"time"

	. "cherry/models"
	"cherry/utils"
)

type Exception struct {
	Msg string
}

func (worker Worker) EmailWorker(payloadJson *[]byte) (err error) {
	start := time.Now().UnixNano()
	var payload MailerPayload
	json.Unmarshal([]byte(*payloadJson), &payload)
	exception := Exception{}
	excute(start, &payload, &worker, &exception)
	return
}

func excute(start int64, payload *MailerPayload, worker *Worker, exception *Exception) {
	defer func(e *Exception) {
		r := recover()
		if r != nil {
			e.Msg = fmt.Sprintf("%v", r)
		}
	}(exception)
	reflect.ValueOf(payload).MethodByName(payload.I18nFuncName()).Call([]reflect.Value{})
	t, err := template.New("content.html").Funcs(payload.FuncMap).ParseFiles(
		"public/workers/emailWorker/content.html",
		"public/workers/emailWorker/head.html",
		"public/workers/emailWorker/footer.html",
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
	for _, email := range payload.Emails {
		if err = utils.SendMail(email, payload.Subject, tpl.String()); err == nil {
			worker.LogInfo((time.Now().UnixNano()-start)/1000000, " ms, payload: ", payload, err)
		} else {
			worker.LogError((time.Now().UnixNano()-start)/1000000, " ms, payload: ", payload, err)
			return
		}
	}
	if err != nil {
		worker.LogError((time.Now().UnixNano()-start)/1000000, " ms, payload: ", payload, err)
	}
}
