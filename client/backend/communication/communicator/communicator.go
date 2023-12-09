package communicator

import (
	"log"
	"sdle/m/v2/communication/rabbitMQ"
	"sdle/m/v2/utils/messageStruct"

	amqp "github.com/rabbitmq/amqp091-go"
)


func addListToClient(ch *amqp.Channel, queue *amqp.Queue, exchangeName string, listsToAdd chan string){

	for {
		select {
			case listurl:= <-listsToAdd:
				rabbbitmq.BindRoutingKeys(ch, queue, exchangeName, "url."+listurl)
		}
	}
}

func sendMessageToServer(ch *amqp.Channel, exchangeName string, messagesToSend chan messageStruct.MessageStruct){
	
	for {
		select {
			case message:= <-messagesToSend:
				log.Println(ch.IsClosed())
				rabbbitmq.PublishMessage("text/json", string(message.ToJSON()), ch, exchangeName, "server/url."+message.ListURL)
		}
	}
	
}

/*
	listsToAdd := make(chan string, 100) // Se quiser ouvir uma lista, escrevo o URl da lista que quero ouvir no canal (listsToAdd <- url) 
	messagesToSend := make(chan messageStruct.MessageStruct, 100) // Se quiser enviar uma mensagem, escrevo o messageStruct da mensagem que quero enviar no canal (messagesToSend <- messageStruct) 
*/

func StartClientCommunication(listsToAdd chan string, messagesToSend chan messageStruct.MessageStruct) {

	// <------------ RabbitMQ Boiler plate ------------>
	conn, ch := rabbbitmq.CreateChannel()

	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	rabbbitmq.DeclareExchange(ch, exchangeName)
	
	q := rabbbitmq.DeclareQueue(ch, "")
	
	// <----------------------------------------------->

	go addListToClient(ch, q, exchangeName, listsToAdd)
	
	messages := rabbbitmq.CreateConsumerChannel(ch, q)
	go rabbbitmq.HandleIncomingMessages(messages)

	go sendMessageToServer(ch, exchangeName, messagesToSend)

}