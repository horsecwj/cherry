package sneakerWorkers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cherry/initializers"
	. "cherry/models"
	"cherry/utils"
)

func (worker Worker) GrantNotifyWorker(payloadJson *[]byte) (err error) {
	payload := make(map[string]string)
	json.Unmarshal([]byte(*payloadJson), &payload)
	payload["timestamp"] = strconv.Itoa(int(time.Now().Unix()))

	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var service Service
	mainDB.Where("id = ?", payload["service_id"]).First(&service)
	if service.CanNotGrant() {
		return
	}
	var transfer Transfer
	mainDB.Where("id = ?", payload["transfer_id"]).First(&transfer)
	if transfer.IsDone() || transfer.IsCanceled() {
		return
	}
	mainDB.Save(&transfer)
	mainDB.DbCommit()

	urlStruct, _ := url.Parse(service.GrantUrl)
	values := urlStruct.Query()
	for k, _ := range values {
		payload[k] = values.Get(k)
	}
	for k, v := range payload {
		values.Add(k, v)
	}
	var signature string
	signature, _ = initializers.PrivateKeySign("POST", urlStruct.Path, &service, &payload)
	values.Add("signature", signature)
	ctx, cancelFun := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFun()
	b, err := json.Marshal(payload)
	if err != nil {
		worker.LogInfo(err.Error())
		return
	}
	req, _ := http.NewRequest(http.MethodPost, service.GrantUrl, bytes.NewReader(b))
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		worker.LogInfo(err.Error())
		return
	}
	db := utils.MainDbBegin()
	defer db.DbRollback()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || strings.ToLower(string(body)) != "success" {
		worker.LogInfo(err.Error())
		db.Save(&TransferNotifyLog{
			TransferId:   transfer.Id,
			RequestUrl:   service.GrantUrl,
			RequestBody:  values.Encode(),
			ErrorInfo:    err.Error(),
			ResponseInfo: string(body),
		})
	} else {
		transfer.State = "done"
		db.Save(&transfer)
	}
	db.DbCommit()
	return
}
