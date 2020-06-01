package main

import (
	"log"
	"os"
	"os/signal"

	sneaker "github.com/oldfritter/sneaker-go"
	"github.com/streadway/amqp"

	"cherry/initializers"
	"cherry/utils"
	"cherry/workers/sneakerWorkers"
)

func main() {
	initialize()
	sneakerWorkers.InitWorkers()

	StartAllWorkers()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	closeResource()
}

func initialize() {
	initializers.InitializeAmqpConfig()
	utils.SetLogAndPid("workers")
}

func closeResource() {
	initializers.CloseAmqpConnection()
}

func StartAllWorkers() {
	for _, w := range sneakerWorkers.AllWorkers {
		if !w.Banned {
			for i := 0; i < w.GetThreads(); i++ {
				go func(w sneakerWorkers.Worker) {
					sneaker.SubscribeMessageByQueue(initializers.RabbitMqConnect.Connection, w, amqp.Table{})
					log.Println("stated ", w.GetName())
				}(w)
			}
		}
	}
}
