package initializers

// import (
//   "encoding/base64"
//   "io/ioutil"
//   "log"
//   "math/rand"
//   "path/filepath"
//
//   "github.com/objcoding/wxpay"
// )
//
// func InitWcpay() {
//   account := wxpay.NewAccount(CurrentEnv.Wechat["appid"], CurrentEnv.Wechat["mch_id"], CurrentEnv.Wechat["key"], false)
//   client := wxpay.NewClient(account)
//   account.SetCertData("config/pems/apiclient_cert.p12")
//   client.setAccount(account)
//   client.SetHttpConnectTimeoutMs(2000)
//   client.SetHttpReadTimeoutMs(1000)
// }
