package main

import (
	messageStruct "sdle/server/utils/messageStruct"
	"sdle/server/utils/communication/rabbitMQ"
	shoppingList "sdle/server/utils/CRDT/shoppingList"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
	"fmt"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func ShopListExample() shoppingList.ShoppingList{
	
	shopList1 := shoppingList.Create("My List 1")
	shopList2 := shoppingList.Create("My List 2")

	//fmt.Println(shopList1.GetURL())

	shopList1.AddItem("apple", 3)
	shopList1.AddItem("rice", 5)
	
	shopList2.AddItem("pear", 2)
	shopList2.AddItem("rice", 3)
	shopList2.BuyItem("rice")

	/*fmt.Println("\nShop List 1")
	fmt.Println(shopList1.JSON())

	fmt.Println("\nShop List 2")
	fmt.Println(shopList2.JSON())*/


	shopList1.JoinShoppingList(shopList2)

	/*fmt.Println("\nShop List 1 after merging with Shop List 2")
	fmt.Println(shopList1.JSON())*/

	shopList1.JoinShoppingList(shopList2)

	/*fmt.Println("\nShop List 1 after merging with Shop List 2 (again)")
	fmt.Println(shopList1.JSON())

	messageFormat := shopList1.ConvertToMessageFormat("john.doe", messageStruct.Write)

	fmt.Println("\n", string(messageFormat))

	fmt.Println("\n", shoppingList.MessageByteToCRDT(messageFormat))
	shoplistCopy := shoppingList.CreateFromStrings(shopList1.GetURL(), shopList1.GetListName(), shopList1.ListFormatForDatabase(), shopList1.StateFormatForDatabase())

	fmt.Println("\n", shoplistCopy)*/

	return shopList1
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
	body := ShopListExample().ConvertToMessageFormat("john.doe", messageStruct.Write)


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

	crdt := ShopListExample()
	crdt.Url = "2990f9e7-cbbc-4852-9906-471436639829"
	fmt.Println(crdt)
	body := crdt.ConvertToMessageFormat("john.doe", messageStruct.Write)
	
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