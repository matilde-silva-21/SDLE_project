package main

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func CreateChannel() (*amqp.Connection, *amqp.Channel){
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return conn, ch
}

func DeclareExchange(ch *amqp.Channel, exchangeName string) {

	err := ch.ExchangeDeclare(
		exchangeName, // exchange name
		"topic", // exchange type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

}

func DeclareQueue(ch *amqp.Channel, queueName string) *amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &q
}

func BindRoutingKeys(ch *amqp.Channel, queue *amqp.Queue, exchangeName string, topics ...string) {
	for _, topic := range topics {
		err := ch.QueueBind(
			queue.Name,   // queue name
			topic,        // routing key (topic)
			exchangeName, // exchange
			false,
			nil,
		)
	   failOnError(err, "Failed to bind queue to exchange")
	}
 }


func CreateConsumerChannel(ch *amqp.Channel, queue *amqp.Queue)  <-chan amqp.Delivery {

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	return msgs
}

// For example, messages with JSON payload should use application/json
func PublishMessage(contentType string, body string, ch *amqp.Channel, exchangeName string, topics ...string) {
	for _, topic := range topics {
		err := ch.Publish(
			exchangeName,                   // exchange
			topic,            // routing key
			false,                    // mandatory
			false,                    // immediate
			amqp.Publishing{
				ContentType: contentType,
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}
}

// Loops infinitely, waiting for messages. Will keep running indefinitely until the channel is closed.
func HandleIncomingMessages(messages <-chan amqp.Delivery) {

	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	for msg := range messages {
	   log.Printf("[x] %s", msg.Body)
	}

}

// Joining the pieces
func RabbitMQExample() {

	// <------------ Boiler plate ------------>
	conn, ch := CreateChannel()
	
	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	DeclareExchange(ch, exchangeName)
	
	q := DeclareQueue(ch, "")
	
	// <-------------------------------------->
	
	// Tópicos a ser ouvidos pelo orchestrator (todos os URLs)
	BindRoutingKeys(ch, q, exchangeName, "url.*")
	
	messages := CreateConsumerChannel(ch, q)

	HandleIncomingMessages(messages)
}
