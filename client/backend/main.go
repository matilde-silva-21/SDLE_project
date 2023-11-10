package main

import (
	"database/sql"
	"fmt"
	"log"
	"sdle/m/v2/api"
	"sdle/m/v2/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET"},
		ExposeHeaders: []string{"Access-Control-Allow-Headers"},
		AllowHeaders: []string{"Content-Type, Access-Control-Allow-Credentials, Access-Control-Allow-Headers, Access-Control-Allow-Methods, Access-Control-Allow-Origin"},
	}))
	
	router.POST("/login", api.Login)
	router.GET("/lists", api.GetShoppingLists)
	router.POST("/lists/create", api.CreateShoppingList)
	router.POST("/lists/remove", api.RemoveShoppingList)
	router.GET("/lists/:url", api.GetShoppingList)
	router.POST("/lists/:url/add", api.AddItemToShoppingList)
	router.POST("/lists/:url/remove", api.RemoveItemFromShoppingList)

	router.Run("localhost:8080")
}
