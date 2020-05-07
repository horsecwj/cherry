package models

type TransferNotifyLog struct {
	CommonModel
	TransferId   int    `json:"transfer_id"`
	RequestUrl   string `json:"request_url"`
	RequestBody  string `gorm:"type:text" json:"request_body"`
	ErrorInfo    string `gorm:"type:text" json:"error_info"`
	ResponseInfo string `gorm:"type:text" json:"response_info"`
}
