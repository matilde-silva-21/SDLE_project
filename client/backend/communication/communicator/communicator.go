package communicator

import (
	"fmt"
	"log"
	"sdle/m/v2/communication/rabbitMQ"
	"sdle/m/v2/utils/messageStruct"
	"sdle/m/v2/utils/CRDT/shoppingList"
	"sdle/m/v2/database"
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
				rabbbitmq.PublishMessage("text/json", string(message.ToJSON()), ch, exchangeName, "server/url."+message.ListURL)
		}
	}
	
}

func ReadAndMergeCRDT(serverMsg messageStruct.MessageStruct, repo *database.SQLiteRepository) (error){
	// Get the list sent in the message
	remoteList := shoppingList.MessageStructToCRDT(serverMsg)

	// Get the list in the local db 
	id, _ := database.GetIDByURL(repo, serverMsg.ListURL)
	dbList := remoteList.ToDatabaseShoppingList(id)

	localList, err := dbList.Read(repo)
	if(err != nil){
		log.Print("Error reading from memory.")
		return err
	}
	localCRDT := shoppingList.DatabaseShoppingListToCRDT(localList.(*database.ShoppingListModel))
	fmt.Println(localCRDT)
	fmt.Println()

	// Join the lists
	localCRDT.JoinShoppingList(remoteList)
	localList = localCRDT.ToDatabaseShoppingList(id)

	//Save the new list in the local db
	err = localList.Update(repo, localList)
	if(err != nil){
		log.Print("Error writing to memory.")
		return err
	}
	return nil
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
	
	go sendMessageToServer(ch, exchangeName, messagesToSend)
	
	rabbbitmq.HandleIncomingMessages(messages)

}