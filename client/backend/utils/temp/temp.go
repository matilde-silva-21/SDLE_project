package main

import (
	"database/sql"
	"fmt"
	"log"
	"sdle/m/v2/database"
	"sdle/m/v2/utils/CRDT/shoppingList"
	"sdle/m/v2/utils/messageStruct"
	_ "github.com/mattn/go-sqlite3"
)

func joinReceivedList(serverMsg messageStruct.MessageStruct, db *sql.DB) {
	// Retrieve 


	// // Assuming you have a method to convert the message body into a ShoppingList

	remoteList := shoppingList.MessageFormatToCRDT(serverMsg.ToJSON())
	fmt.Println(remoteList.GetURL())

	// // Join the lists
	// localList.JoinShoppingList(remoteList)

	// // Assuming you have a method to update the local list in the database
	// updateLocalList(localList, db)
}

func retrieveLocalList(listURL string, db *sql.DB) (shoppingList.ShoppingList, error) {
	// Assuming you have a method to retrieve the list from the database based on the list URL
	// This is just a placeholder, and you should replace it with your actual implementation
	// The implementation might include querying the database, unmarshaling the JSON, etc.
	var list shoppingList.ShoppingList
	return list, nil
}

func updateLocalList(list shoppingList.ShoppingList, db *sql.DB) {
	// Assuming you have a method to update the list in the database
	// This is just a placeholder, and you should replace it with your actual implementation
	// The implementation might include marshaling the list to JSON, updating the database, etc.
}

func main() {
	const filename = "temp.db"
	db, err := sql.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}

	sqliteRepository := database.NewSQLiteRepository(db)
	createError := sqliteRepository.CreateTables()
	if createError != nil {
		fmt.Println(createError.Error())
		return
	}

	// Example: Joining the list from the local db to a list retrieved from the server
	fmt.Println(messageStruct.Update)
	// Create a MessageStruct using the provided details
	serverMsg := messageStruct.CreateMessage("123", "john.doe", messageStruct.Update,
		`{"Name":"My List 1", "List":{"Map":{"apple":{"First":1,"Second":3},"pear":{"First":2,"Second":2},"rice":{"First":3,"Second":2}}}, "State":{"Map":{"pear":{"First":0,"Second":0},"rice":{"First":2,"Second":0}}}}`)
	// Join the received list with the local list
	joinReceivedList(serverMsg, db)
}
