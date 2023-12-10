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

func ReadWriteAndMergeCRDT(serverMsg messageStruct.MessageStruct, repo *database.SQLiteRepository) (error){
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

func ReadAndMergeCRDT(serverMsg messageStruct.MessageStruct, repo *database.SQLiteRepository) (shoppingList.ShoppingList, error) {
	remoteList := shoppingList.MessageStructToCRDT(serverMsg)

	// Get the list in the local db 
	id, _ := database.GetIDByURL(repo, serverMsg.ListURL)
	dbList := remoteList.ToDatabaseShoppingList(id)

	localList, err := dbList.Read(repo)
	if(err != nil){
		log.Print("Error reading from memory.")
		var dummy shoppingList.ShoppingList
		return dummy, err
	}
	localCRDT := shoppingList.DatabaseShoppingListToCRDT(localList.(*database.ShoppingListModel))
	fmt.Println(remoteList)
	
	// Join the lists
	localCRDT.JoinShoppingList(remoteList)
	fmt.Println(localCRDT)

	return localCRDT, nil
}

func WriteListsToDatabase(updatedMap *map[string](shoppingList.ShoppingList), writeListsToDatabase chan string, repo *database.SQLiteRepository) error{


	for {

		select{

			case urlToWrite := <- writeListsToDatabase:

				id, _ := database.GetIDByURL(repo, urlToWrite)

				localCRDT, ok := (*updatedMap)[urlToWrite]
				localList := localCRDT.ToDatabaseShoppingList(id)

				if(!ok){
					writeListsToDatabase <- urlToWrite // Wait until message arrives 
					continue
				}

				if (id != 0){
					err := localList.Update(repo, localList)
					if(err != nil){
						log.Print("Error writing to memory.")
						return err
					}
					log.Printf("Updated memory value for list %s.", urlToWrite)
					return nil
				} else {
					_, err := localList.Create(repo)
					if(err != nil){
						log.Print("Error writing to memory.")
						return err
					}
					log.Printf("Wrote new list to memory %s.", urlToWrite)
					return nil
				}
		}
	}
	
}

/*
	listsToAdd := make(chan string, 100) // Se quiser ouvir uma lista, escrevo o URl da lista que quero ouvir no canal (listsToAdd <- url) 
	messagesToSend := make(chan messageStruct.MessageStruct, 100) // Se quiser enviar uma mensagem, escrevo o messageStruct da mensagem que quero enviar no canal (messagesToSend <- messageStruct) 
*/

func StartClientCommunication(connected chan bool, listsToAdd chan string, messagesToSend chan messageStruct.MessageStruct, writeListsToDatabase chan string, repo *database.SQLiteRepository) error {

	// <------------ RabbitMQ Boiler plate ------------>
	
	conn, ch, err := rabbbitmq.CreateChannel()
	if(err != nil) {
		connected <- false
	}

	for (err != nil){
		conn, ch, err = rabbbitmq.CreateChannel()
	}

	connected <- true
	
	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	rabbbitmq.DeclareExchange(ch, exchangeName)
	
	q := rabbbitmq.DeclareQueue(ch, "")
	
	// <----------------------------------------------->

	go addListToClient(ch, q, exchangeName, listsToAdd)
	
	messages := rabbbitmq.CreateConsumerChannel(ch, q)
	
	go sendMessageToServer(ch, exchangeName, messagesToSend)

	updatedMap := make(map[string] shoppingList.ShoppingList)

	go WriteListsToDatabase(&updatedMap, writeListsToDatabase, repo)
	
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	for msg := range messages {
		messageObject, _ := messageStruct.JSONToMessage(msg.Body)
	   	log.Printf("Received a message for URL %s: %s", messageObject.ListURL, msg.Body)

		update, err := ReadAndMergeCRDT(messageObject, repo)
		if(err == nil){
			updatedMap[messageObject.ListURL] = update
		}
	}

	return nil

}