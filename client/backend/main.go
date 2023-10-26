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
	sqliteRepository.CreateTables()
	sqliteRepository.Seed()

	modelToUpdate := &database.Item{Id: 1}
	modelToUpdate.Update(sqliteRepository, &database.Item{Id: modelToUpdate.Id, Name: "Updated", Done: modelToUpdate.Done})
	
	modelToDelete := database.Item{Id: 2}
	modelToDelete.Delete(sqliteRepository)

	router := gin.Default()
	api.SetDB(sqliteRepository)

	if apiDB := api.GetDB(); apiDB == nil {
		fmt.Print("Database not initialized in API")
		return
	}
	
	router.GET("/items", api.GetAllItems)

	router.Run("localhost:8080")
}
