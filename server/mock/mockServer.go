package main

import (
	"log"
	"net"
	messageStruct "sdle/server/utils/messageStruct"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func tcp() {
	address := "localhost:8080"
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return
	}
	
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}
	
	defer conn.Close()
	
	_, err = conn.Write([]byte("first message"))
	if err != nil {
		log.Print("Error sending message:", err)
	}

	time.Sleep(10 * time.Second)

	_, err = conn.Write([]byte("second message"))
	if err != nil {
		log.Print("Error sending message:", err)
	}

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {

			log.Print("Error reading message: ", err)
			return
		}

		log.Print(string(buffer[:n]))

		rabbit()
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
	tcp()
}