package v1

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	. "cherry/models"
	"cherry/utils"
)

func Accounts(echoContext echo.Context) (err error) {
	upgrader := websocket.Upgrader{}
	user := echoContext.Get("current_user").(User)
	c, err := upgrader.Upgrade(echoContext.Response(), echoContext.Request(), nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	c.SetWriteDeadline(time.Now().Add(time.Minute * 10))
	defer c.Close()
	var params struct {
		CurrencyIds []int `json:"currency_ids"`
	}
	_, m, err := c.ReadMessage()
	json.Unmarshal(m, &params)
	if len(params.CurrencyIds) == 0 {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	err = utils.ListenPubSubChannels(
		ctx,
		func() error {
			return nil
		},
		func(channel string, m *[]byte) error {
			var account Account
			if channel == NotifyAccountWithRedis {
				json.Unmarshal(*m, &account)
				if account.UserId != user.Id {
					return nil
				}
				var in bool
				for _, currencyId := range params.CurrencyIds {
					if currencyId == account.CurrencyId {
						in = true
					}
				}
				if !in {
					return nil
				}

				err := c.WriteMessage(websocket.TextMessage, *m)
				if err != nil {
					log.Println("write:", err)
					cancel()
				}
			}
			return nil
		},
		NotifyAccountWithRedis,
	)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
