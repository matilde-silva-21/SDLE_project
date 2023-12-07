package main

import (
	messageStruct "sdle/server/utils/messageStruct"
	"sdle/server/utils/communication/rabbitMQ"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func rabbit() {

	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a Topic Exchange named "logs"
	err = ch.ExchangeDeclare(
		"logs", // exchange name
		"topic", // exchange type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Create String to match  
	body := messageStruct.CreateMessage("123", "jonh.doe", messageStruct.Write, "CRDT").ToJSON()


	err = ch.Publish(
		"logs",                   // exchange
		"url.123",            // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}


func main(){

	// <------------ Boiler plate ------------>
	conn, ch := rabbitmq.CreateChannel()
	
	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	rabbitmq.DeclareExchange(ch, exchangeName)
	
	q := rabbitmq.DeclareQueue(ch, "")
	
	// <-------------------------------------->
	
	// Tópicos a ser ouvidos pelo orchestrator (todos os URLs)
	rabbitmq.BindRoutingKeys(ch, q, exchangeName, "url.*")

	messages := rabbitmq.CreateConsumerChannel(ch, q)

	go rabbitmq.PrintIncomingMessages(messages)

	body := messageStruct.CreateMessage("123", "jonh.doe", messageStruct.Write, "CRDT").ToJSON()
	

	for {

		err := ch.Publish(
			"logs",                   // exchange
			"server/url.123",            // routing key
			false,                    // mandatory
			false,                    // immediate
			amqp.Publishing{
				ContentType: "text/json",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", body)

		time.Sleep(10*time.Second)
	}

}