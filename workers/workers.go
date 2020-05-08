package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"

	sneaker "github.com/oldfritter/sneaker-go"
	"github.com/streadway/amqp"

	"cherry/initializers"
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

	setLog()
	err := ioutil.WriteFile("pids/workers.pid", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Println(err)
	}
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
				}(w)
			}
		}
	}
}

func setLog() {
	err := os.Mkdir("logs", 0755)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}

	file, err := os.OpenFile("logs/workers.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	log.SetOutput(file)
}