package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"os"
	"rozprochy_rabbit/constants"
	"rozprochy_rabbit/model"
	"rozprochy_rabbit/util"
	"time"
)

func consume(msg []byte, channel *amqp.Channel)(){
	time.Sleep(2 * time.Second)

	msgUnmarshalled := model.Message{}
	err := json.Unmarshal(msg, &msgUnmarshalled)

	util.FailOnError(err, "failed to unmarshall message: " + string(msg))

	println("Finished handling message : " + string(msg))
	body := msgUnmarshalled.Name + ", " + msgUnmarshalled.Injury + ", DONE"
	util.PublishToQueue(channel, []byte(body), "app.doctor."+msgUnmarshalled.DocId)
}


func main() {
	println("Started with parameters: " + os.Args[1] + ", " + os.Args[2])
	conn, err := amqp.Dial(constants.MqUrl)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	util.CreateExchange(constants.ExchangeName, ch)

	elbowQueue := util.CreateQueue("elbow", ch)
	kneeQueue := util.CreateQueue("knee", ch)
	hipQueue := util.CreateQueue("hip", ch)

	util.BindQueue(elbowQueue, ch, constants.ElbowRoutingKey, constants.ExchangeName)
	util.BindQueue(kneeQueue, ch, constants.KneeRoutingKey, constants.ExchangeName)
	util.BindQueue(hipQueue, ch, constants.HipRoutingKey, constants.ExchangeName)

	msgs1, err := ch.Consume(
		os.Args[1], 		// queue
		"",     	// consumer
		false,  	// auto-ack
		false, 	// exclusive
		false,  	// no-local
		false,  	// no-wait
		nil,    		// args
	)
	util.FailOnError(err, "Failed to register a consumer on " + os.Args[1])

	msgs2, err := ch.Consume(
		os.Args[2], 		// queue
		"",     	// consumer
		false,  	// auto-ack
		false, 	// exclusive
		false,  	// no-local
		false,  	// no-wait
		nil,    		// args
	)

	util.FailOnError(err, "Failed to register a consumer on " + os.Args[2])

	forever := make(chan bool)
	go func() {
		for {
			select {
			case msg := <-msgs1:
				go consume(msg.Body, ch)
			case msg := <-msgs2:
				go consume(msg.Body, ch)
			}
		}
	}()
	<- forever
}
