package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"rozprochy_rabbit/constants"
	"rozprochy_rabbit/util"
	"time"
)

/*
1. Check if app.logging exists
2. Randomly send topic requests to exchange
*/

func main() {
	conn, err := amqp.Dial(constants.MqUrl)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	util.CreateExchange(constants.ExchangeName, ch)
	util.CreateLoggingExchange(constants.LoggingExchangeName, ch)

	q := util.CreateQueue(constants.AdminLoggingQueue, ch)
	util.BindQueue(q, ch, constants.AdminRoutingKey, constants.ExchangeName)
	util.BindQueue(q, ch, constants.AdminRoutingKey, constants.LoggingExchangeName)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("[%s] Received a message: %s", time.Now().String(), d.Body)
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			println("[" + time.Now().String() + "]Publishing to all queues")

			body, _ := json.Marshal("info message")
			util.PublishToExchange(ch, body, "", constants.LoggingExchangeName)
		}
	}()
	<-forever

}
