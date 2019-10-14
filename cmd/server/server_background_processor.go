package main

import (
	cfg "calendar/internal/config"
	"calendar/internal/interfaces/postgres"
	"calendar/internal/interfaces/rabbitmq"
	lg "calendar/internal/logger"
	"calendar/internal/services"
)

func main() {
	logger := lg.GetLogger(cfg.GetConfig())
	logger.Info("Service loading!")

	rabbit, err := rabbitmq.NewRabbitMQ(logger, cfg.GetConfig())
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer rabbit.Close()

	psql, err := postgres.NewPSQL(logger, cfg.GetConfig())
	if err != nil {
		logger.Error(err.Error())
	}
	defer psql.Close()

	bgProcessor := services.BackgroundProcessor{
		Logger:   logger,
		RabbitMQ: rabbit,
		PSQL:     psql,
	}

	forever := make(chan bool)

	err = bgProcessor.Run()
	if err != nil {
		logger.Error(err.Error())
	}
	<-forever
}
