package api

import (
	"encoding/base64"
	"net/http"
	"sdle/m/v2/database"
	"fmt"
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

func GetItem(c* gin.Context) {
	var i database.Item

	if err := c.ShouldBindUri(&i); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}

	item, err := i.Read(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}

	c.IndentedJSON(http.StatusOK, item)
}

func CreateShoppingList(c *gin.Context) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	username, cookieErr := getUsernameFromCookie(c)
	if cookieErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error reading username from cookie"})
		return
	}

	var shoppingList database.ShoppingList
	if err := c.ShouldBind(&shoppingList); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to shopping list"})
		return
	}

	newShoppingListModel, createErr := shoppingList.Create(db)
	if createErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating shopping list"})
		return
	}

	newShoppingList := newShoppingListModel.(*database.ShoppingList)

	var userList database.UserList
	userList.ListID = newShoppingList.Id
	userList.UserID = username
	userList.Create(db)
	
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "shopping list created successfully"})
}

func RemoveShoppingList(c *gin.Context) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var shoppingList database.ShoppingList
	if err := c.ShouldBind(&shoppingList); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to shopping list"})
		return
	}

	shoppingListModel, _ := shoppingList.Read(db)
	shoppingListObj := shoppingListModel.(*database.ShoppingList)

	deleteErr := shoppingListObj.Delete(db)
	if deleteErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error deleting shopping list"})
		return
	}

	username, cookieErr := getUsernameFromCookie(c)
	if cookieErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error getting username from cookie"})
		return
	}

	userList := database.UserList{ListID: shoppingListObj.Id, UserID: username}
	err := userList.Delete(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error deleting userList"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"msg": "shopping list deleted successfully"})
}

func GetShoppingLists(c *gin.Context) {

	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	username, err := getUsernameFromCookie(c)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	u := database.User{Username: username}
	userModel, err := u.Read(db)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user unauthorized"})
		return
	}

	user := userModel.(*database.User)

	userLists, readErr := user.ReadUserLists(db)

	if readErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
	}
	
	c.IndentedJSON(http.StatusOK, userLists)
}

func GetShoppingList(c *gin.Context) {

	fmt.Println(c)

	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var shoppingList database.ShoppingList

	fmt.Println(shoppingList)
	
	if err := c.ShouldBindUri(&shoppingList); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list url not found"})
		return
	}

	shoppingListModel, err := shoppingList.Read(db)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading shopping list"})
		return
	}

	shoppingListObj := shoppingListModel.(*database.ShoppingList)
	
	items, err := shoppingListObj.GetShoppingListItems(db)
	
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error retrieving shopping list items"})
		return
	}

	username, _ := getUsernameFromCookie(c)
	userList := database.UserList{ListID: shoppingListObj.Id, UserID: username}
	userListObj, _ := userList.Read(db)

	if userListObj == nil {
		_, createUserListErr := userList.Create(db)
		if createUserListErr != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating UserList entry"})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, items)
}

func AddItemToShoppingList(c *gin.Context) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var shoppingList database.ShoppingList
	if err := c.ShouldBindUri(&shoppingList); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list url not found"})
		return
	}

	shoppingListModel, err := shoppingList.Read(db)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading shopping list"})
		return
	}

	shoppingListObj := shoppingListModel.(*database.ShoppingList)
	var item database.Item

	bindingErr := c.ShouldBind(&item) 
	if bindingErr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to item"})
		return
	}

	item.List = *shoppingListObj
	_, createErr := item.Create(db)
	if createErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating Item"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"msg": "Item added to list successfully"})
}

func RemoveItemFromShoppingList(c *gin.Context) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var item database.Item
	if err := c.ShouldBind(&item); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to item"})
		return
	}

	itemModel, readErr := item.Read(db)
	if readErr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading item from db"})
		return
	}

	itemModel.Delete(db)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "item deleted successfully"})
}

func Login(c *gin.Context) {
	if isLoggedIn(c) {
		c.Redirect(http.StatusFound, "/lists")
		return
	}

	var user database.User
	err := c.ShouldBind(&user)

	_, readErr := user.Read(db)
	if readErr != nil {
		_, createErr := user.Create(db)
		if createErr != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating or reading user"})
			return
		}
	}

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "username not found"})
		return
	}

	cookie := base64.StdEncoding.EncodeToString([]byte(user.Username))
	c.SetCookie("session", cookie, 0, "/", "localhost", false, false)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "user logged in successfully"})
}

func getUsernameFromCookie(c *gin.Context) (string, error) {
	cookie, cookieErr := c.Cookie("session")
	if cookieErr != nil {
		return "", cookieErr
	}

	username, decodeErr := base64.StdEncoding.DecodeString(cookie)
	if decodeErr != nil {
		return "", decodeErr
	}

	usernameStr := string(username)
	user := database.User{Username: usernameStr}
	_, readErr := user.Read(db)
	if readErr != nil {
		return "", readErr
	}

	return usernameStr, nil
}

func isLoggedIn(c *gin.Context) bool {
	_, err := getUsernameFromCookie(c)

	return (err == nil)
}

func SetDB (database *database.SQLiteRepository) {
	db = database
}

func GetDB() *database.SQLiteRepository {
	return db
}
