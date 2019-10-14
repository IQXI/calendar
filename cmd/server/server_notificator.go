package main

import (
	cfg "calendar/internal/config"
	"calendar/internal/interfaces/rabbitmq"
	lg "calendar/internal/logger"
	"fmt"
)

func main() {
	logger := lg.GetLogger(cfg.GetConfig())
	logger.Info("Service loading!")

	rabbit, err := rabbitmq.NewRabbitMQ(logger, cfg.GetConfig())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer rabbit.Close()

	forever := make(chan bool)
	go func() {
		err = rabbit.Receive()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}()

	go func() {
		for d := range rabbit.Storage {
			logger.Info(fmt.Sprintf("Новая встреча у %v в %v \nТема: %v \nОписание: %v", d.Owner, d.DateTime, d.Header, d.Description))
		}
	}()

	logger.Info("Notificator started!")
	<-forever
}
