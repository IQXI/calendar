package services

import (
	"calendar/internal/interfaces/postgres"
	"calendar/internal/interfaces/rabbitmq"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type BackgroundProcessor struct {
	RabbitMQ rabbitmq.RabbitMQ
	PSQL     postgres.PSQL
	Logger   *zap.Logger
}

func (bp *BackgroundProcessor) Run() error {

	go func() {

		start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Now().Location())
		stop := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Now().Location())

		for {

			bp.Logger.Info(fmt.Sprintf("Checking %v  --  %v", start, stop))
			//чекаем базу на наличие сообщений для рассылки
			events, err := bp.PSQL.GetPublishEvents(start, stop)
			if err != nil {
				bp.Logger.Error(err.Error())
			} else {
				if events != nil {
					for _, event := range events {

						bp.Logger.Info(fmt.Sprintf("Publish to RMQ %v", event))

						//постим в RabbitMQ
						err = bp.RabbitMQ.Publish(event)
						if err != nil {
							bp.Logger.Error(err.Error())
						}
					}
				}
			}

			//смещаем интервал времени
			start = stop
			stop = stop.Add(time.Second * 10)
			bp.Logger.Info(fmt.Sprintf("Sleep %v", time.Second*10))
			//повторяем раз в 10 сек
			time.Sleep(10 * time.Second)
		}

	}()

	return nil
}
