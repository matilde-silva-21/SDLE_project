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

// Loops infinitely, waiting for messages. Will keep running indefinitely until the channel is closed.
func main() {

	// <------------ Boiler plate ------------>
	conn, ch := CreateChannel()
	
	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	DeclareExchange(ch, exchangeName)
	
	q := DeclareQueue(ch, "")
	
	// <-------------------------------------->

	// Tópicos a ser ouvidos pelo cliente (poderão ser os URLs da shoppingList)
	BindRoutingKeys(ch, q, exchangeName, "info.server", "topic.second", "topic.third")

	messages := CreateConsumerChannel(ch, q)
 
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	for msg := range messages {
	   log.Printf(" [x] %s", msg.Body)
	}
}
