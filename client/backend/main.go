package main

import (
	"database/sql"
	"fmt"
	"log"
	"sdle/m/v2/api"
	"sdle/m/v2/database"
	"sdle/m/v2/communication/communicator"
	"sdle/m/v2/utils/messageStruct"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Se quiser ouvir uma lista, escrevo o URl da lista que quero ouvir no canal (listsToAdd <- url) 
	listsToAdd := make(chan string, 100)
	
	// Se quiser enviar uma mensagem, escrevo o messageStruct da mensagem que quero enviar no canal (messagesToSend <- messageStruct) 
	messagesToSend := make(chan messageStruct.MessageStruct, 100)

	writeListsToDatabase := make(chan string, 100)

	const filename = "local.db"
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

	go communicator.StartClientCommunication(listsToAdd, messagesToSend, writeListsToDatabase, sqliteRepository)

	router := gin.Default()
	api.SetDB(sqliteRepository)

	if apiDB := api.GetDB(); apiDB == nil {
		fmt.Print("Database not initialized in API")
		return
	}

	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"POST", "GET"},
		ExposeHeaders: []string{"Access-Control-Allow-Headers"},
		AllowHeaders: []string{"Content-Type, Access-Control-Allow-Credentials, Access-Control-Allow-Headers, Access-Control-Allow-Methods, Access-Control-Allow-Origin"},
	}))

	router.Use(api.SetMessagesToSendChannel(messagesToSend))
	router.Use(api.SetListsToAddChannel(listsToAdd))
	router.Use(api.SetWriteListsToDatabaseChannel(writeListsToDatabase))
	
	router.POST("/login", api.Login)
	router.GET("/lists", api.GetShoppingLists)
	router.POST("/lists/create", api.CreateShoppingList)
	router.POST("/lists/remove", api.RemoveShoppingList)
	router.GET("/lists/:url", api.GetShoppingList)
	router.POST("/lists/:url/add", api.AddItemToShoppingList)
	router.POST("/lists/:url/remove", api.RemoveItemFromShoppingList)
	router.POST("lists/:url/upload", api.SetMessagesToSendChannel(messagesToSend),func(c *gin.Context) {
		connected := GetConnectedStatus()
		api.UploadList(c, connected)
	})
	router.POST("lists/:url/fetch", api.SetListsToAddChannel(listsToAdd), func(c *gin.Context) {
		connected := GetConnectedStatus()
		api.FetchList(c, connected)
	})

	router.Run("localhost:8082")
}

func GetConnectedStatus() bool{
	return true
}