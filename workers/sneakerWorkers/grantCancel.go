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
	if transfer.IsGranting() {
		transfer.State = "canceled"
	} else if transfer.IsDone() {
		return
	} else {
		worker.LogError("transfer type is error, transfer id: ", payload["transfer_id"])
	}
	mainDB.Save(&transfer)
	mainDB.DbCommit()
	return
}
