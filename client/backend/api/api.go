package api

import (
	"net/http"
	"sdle/m/v2/database"

	"github.com/gin-gonic/gin"
)

var db *database.SQLiteRepository

func GetAllItems(c *gin.Context) {
	item := database.Item{}
	items, err := item.ReadAll(db)

	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, items)
}

func SetDB (database *database.SQLiteRepository) {
	db = database
}

func GetDB() *database.SQLiteRepository {
	return db
}