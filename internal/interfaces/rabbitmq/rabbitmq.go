package rabbitmq

import (
	"calendar/internal/structs"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	queue      amqp.Queue
	channel    amqp.Channel
	connection amqp.Connection
	Storage    chan structs.Event
	logger     *zap.Logger
	config     map[string]interface{}
}

func NewRabbitMQ(logger *zap.Logger, config map[string]interface{}) (RabbitMQ, error) {

	user := config["rabbitmq.user"]
	password := config["rabbitmq.password"]
	host := config["rabbitmq.host"]
	port := config["rabbitmq.port"]
	vhost := config["rabbitmq.vhost"]

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/%v", user, password, host, port, vhost))
	if err != nil {
		return RabbitMQ{}, err
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return RabbitMQ{}, err
	}
	//defer ch.Close()

	q, err := ch.QueueDeclare(
		"calendar", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return RabbitMQ{}, err
	}
	return RabbitMQ{
		queue:      q,
		channel:    *ch,
		connection: *conn,
		Storage:    make(chan structs.Event),
		logger:     logger,
		config:     config,
	}, nil
}

func (r *RabbitMQ) Close() error {
	err := r.channel.Close()
	if err != nil {
		return err
	}

	err = r.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Publish(body structs.Event) error {

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	r.logger.Info(fmt.Sprintf(" [x] Sent %s", body))
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Receive() error {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var jsEvent structs.Event
			err = json.Unmarshal(d.Body, &jsEvent)
			if err != nil {
				r.logger.Error(fmt.Sprintf("Unmarshal error %v", err))
			} else {
				r.logger.Info(fmt.Sprintf("Publish to chan: %v", jsEvent))
				r.Storage <- jsEvent
			}
		}
	}()

	r.logger.Info("Receiver Waiting for messages.")
	<-forever
	return nil
}
