package util

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"rozprochy_rabbit/model"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func CreateQueue(queueName string, channel *amqp.Channel) amqp.Queue {
	q, err := channel.QueueDeclare(
		queueName, 				// name
		false,     		// durable
		false,     	// delete when unused
		false,     	// exclusive
		false,     		// no-wait
		nil,       		// arguments
	)
	FailOnError(err, "Failed to declare a queue")

	return q
}

func PublishMessageToQueue(channel *amqp.Channel, message model.Message, routingKey string) () {
	body, _ := json.Marshal(message)

	PublishToQueue(channel, body, routingKey)
}

func PublishToQueue(channel *amqp.Channel, body []byte, routingKey string) () {
	err := channel.Publish(
		"exchange1",         	// exchange
		routingKey, 					// routing key
		false,      			// mandatory
		false,      			// immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        body,
		})

	FailOnError(err, "Failed to send message")
}

func CreateExchange(exchangeName string, channel *amqp.Channel)(){
	err := channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare an exchange")
}

func BindQueue(queue amqp.Queue, channel *amqp.Channel, routingKey string, exchangeName string) () {
	err := channel.QueueBind(
		queue.Name, // queue name
		routingKey,     // routing key
		exchangeName, // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind a queue " + queue.Name)
}