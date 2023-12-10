package main

import (
	"database/sql"
	"fmt"
	"log"
	"sdle/m/v2/database"
	"sdle/m/v2/utils/messageStruct"
	"sdle/m/v2/communication/communicator"
	_ "github.com/mattn/go-sqlite3"
)

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

	seedError := sqliteRepository.Seed()
	if seedError != nil {
		fmt.Println(seedError.Error())
	}

	// Example: Joining the list from the local db to a list retrieved from the server
	serverMsg := messageStruct.CreateMessage("123", "john.doe", "Write",
		`{"Name":"My List 1", "List":{"Map":{"apple":{"First":1,"Second":3},"pear":{"First":2,"Second":2},"rice":{"First":3,"Second":2}}}, "State":{"Map":{"pear":{"First":0,"Second":0},"rice":{"First":2,"Second":0}}}}`)
	// Join the received list with the local list
	_, err = communicator.ReadAndMergeCRDT(serverMsg, sqliteRepository)
	if(err != nil){
		log.Print("Error joining the lists.")
	} else{
		log.Print("Join successful!")	
	}
}
