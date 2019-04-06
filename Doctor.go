package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"math/rand"
	"os"
	"rozprochy_rabbit/constants"
	"rozprochy_rabbit/model"
	"rozprochy_rabbit/util"
	"time"
)

/*
1. Declare queue with doctor name app.doctor.doctor_name
2. Check if app.doctor.requests is alive, if not create it
3. Randomly send requests to exchange
4. Process responses from app.doctor.doctor_name
 */

func sendMessagesToTechnicians(channel *amqp.Channel, docId string)() {
	injuries := [...]string{"knee","elbow","hip"}
	names := [...]string{"john", "mark", "ed"}

	rand.Seed(time.Now().UTC().UnixNano())

	for {
		index := rand.Intn(3)
		message := model.Message{
			Injury: injuries[index],
			Name: names[index],
			DocId: docId}

		util.PublishMessageToQueue(channel, message, routingKeyFromInjury(injuries[index]))

		marshaled,_ := json.Marshal(message)

		println("Message send: " + string(marshaled))
		time.Sleep(9 * time.Second)
	}
}

func routingKeyFromInjury(injury string) string {
	switch injury {
	case "knee": return constants.KneeRoutingKey
	case "hip": return constants.HipRoutingKey
	case "elbow": return constants.ElbowRoutingKey
	default:
		return constants.KneeRoutingKey
	}
}

func main() {
	conn, err := amqp.Dial(constants.MqUrl)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	util.CreateExchange(constants.ExchangeName, ch)

	queueName := os.Args[1]

	queue := util.CreateQueue(queueName, ch)

	routingKey := "app.doctor." + os.Args[1]
	util.BindQueue(queue, ch, routingKey, constants.ExchangeName)

	msgs, err := ch.Consume(
		queue.Name, 		// queue
		"",     	// consumer
		false,  	// auto-ack
		false, 	// exclusive
		false,  	// no-local
		false,  	// no-wait
		nil,    		// args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			println( "["+time.Now().String() + "]Received message: " + string(d.Body))
		}
	}()

	go sendMessagesToTechnicians(ch, queueName)
	<- forever


}
