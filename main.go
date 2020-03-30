package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Paintxd/compassitoMail/model"
	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

func main() {
	fmt.Println("Start messaging...")

	startMessaging()
}

func startMessaging() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("compassito.queue", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		for d := range messageChannel {
			// Received message

			Info := &model.Info{}

			err := json.Unmarshal(d.Body, Info)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Sending email - %s, %s, %g ", Info.Nome, Info.Email, Info.Valor)
				Info.Send()
			}

		}
	}()

	<-stopChan

}
