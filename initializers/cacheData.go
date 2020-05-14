package initializers

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	. "cherry/models"
	"cherry/utils"
)

type Payload struct {
	Update string `json:"update"`
}

func InitCacheData() {
	db := utils.MainDbBegin()
	defer db.DbRollback()
	InitAllCurrencies(db)
	InitConfigInDB(db)
}

func LoadCacheData() {
	InitCacheData()
	go func() {
		channel, err := RabbitMqConnect.Channel()
		if err != nil {
			log.Println(err)
			return
		}
		channel.ExchangeDeclare(AmqpGlobalConfig.Exchange["fanout"]["default"], "fanout", true, false, false, false, nil)
		queue, err := channel.QueueDeclare("", true, true, false, false, nil)
		if err != nil {
			log.Println(err)
			return
		}
		channel.QueueBind(queue.Name, queue.Name, AmqpGlobalConfig.Exchange["fanout"]["default"], false, nil)
		msgs, _ := channel.Consume(queue.Name, "", true, false, false, false, nil)
		for d := range msgs {
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err == nil {
				reflect.ValueOf(&payload).MethodByName(strings.Title(strings.ToLower(payload.Update))).Call([]reflect.Value{})
			} else {
				log.Println(err)
			}
		}
		return
	}()
}

func (payload *Payload) Currencies() {
	db := utils.MainDbBegin()
	defer db.DbRollback()
	InitAllCurrencies(db)
}

func (payload *Payload) Configs() {
	db := utils.MainDbBegin()
	defer db.DbRollback()
	InitConfigInDB(db)
}
