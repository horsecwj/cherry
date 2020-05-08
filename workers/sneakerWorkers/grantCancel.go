package sneakerWorkers

import (
	"encoding/json"

	. "cherry/models"
	"cherry/utils"
)

func (worker Worker) GrantCancelWorker(payloadJson *[]byte) (err error) {
	payload := make(map[string]string)
	json.Unmarshal([]byte(*payloadJson), &payload)

	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()
	var transfer Transfer
	mainDB.Where("id = ?", payload["transfer_id"]).First(&transfer)
	if transfer.IsDone() || transfer.IsCanceled() {
		return
	} else {
		transfer.State = "canceled"
	}
	mainDB.Save(&transfer)
	mainDB.DbCommit()
	return
}