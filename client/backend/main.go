package main

import (
	"database/sql"
	"fmt"
	"log"
	"sdle/m/v2/api"
	"sdle/m/v2/database"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
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
	
	seedError := sqliteRepository.Seed()
	if seedError != nil {
		fmt.Println(seedError.Error())
	}

	router := gin.Default()
	api.SetDB(sqliteRepository)

	if apiDB := api.GetDB(); apiDB == nil {
		fmt.Print("Database not initialized in API")
		return
	}
	
	router.POST("/login", api.Login)
	router.GET("/lists", api.GetShoppingLists)
	router.POST("/lists/create", api.CreateShoppingList)
	router.POST("/lists/remove", api.RemoveShoppingList)
	router.GET("/lists/:url", api.GetShoppingList)
	router.POST("/lists/:url/add", api.AddItemToShoppingList)
	router.POST("/lists/:url/remove", api.RemoveItemFromShoppingList)

	router.Run("localhost:8080")
}
