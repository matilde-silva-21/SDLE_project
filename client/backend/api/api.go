package api

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	//"sdle/m/v2/communication/communicator"
	"sdle/m/v2/database"
	"sdle/m/v2/utils/CRDT/shoppingList"
	"sdle/m/v2/utils/messageStruct"

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

func GetItem(c *gin.Context) {
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

	var shoppingListModel database.ShoppingListModel
	if err := c.ShouldBind(&shoppingListModel); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to shopping list"})
		return
	}

	shoppingListCRDT := shoppingList.Create(shoppingListModel.Name)

	listCRDT := shoppingListCRDT.ListFormatForDatabase()
	stateCRDT := shoppingListCRDT.StateFormatForDatabase()

	shoppingListModel.List = listCRDT
	shoppingListModel.State = stateCRDT

	newShoppingListModel, createErr := shoppingListModel.Create(db)

	if createErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating shopping list"})
		return
	}

	newShoppingList := newShoppingListModel.(*database.ShoppingListModel)

	var userList database.UserList
	userList.ListID = newShoppingList.Id
	userList.UserID = username
	_, userListErr := userList.Create(db)
	if userListErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating user list"})
		return
	}

	c.IndentedJSON(http.StatusOK, newShoppingList)
}

func RemoveShoppingList(c *gin.Context) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var shoppingList database.ShoppingListModel
	if err := c.ShouldBind(&shoppingList); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to shopping list"})
		return
	}

	shoppingListModel, _ := shoppingList.Read(db)
	shoppingListObj := shoppingListModel.(*database.ShoppingListModel)

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

	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var sl database.ShoppingListModel

	if err := c.ShouldBindUri(&sl); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list url not found"})
		return
	}

	shoppingListModel, err := sl.Read(db)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading shopping list"})
		return
	}

	shoppingListObj := shoppingListModel.(*database.ShoppingListModel)
	shoppingListCRDT := shoppingList.DatabaseShoppingListToCRDT(shoppingListObj)

	mapItems := shoppingListCRDT.GetItemsAndTheirQuantity()
	finalItems := []database.Item {}

	for name, quantity := range mapItems{
		item := database.Item {Id: 0, Name: name, Done: shoppingListCRDT.CheckIfItemBought(name), Quantity: int64(quantity), List: *shoppingListObj}
		finalItems = append(finalItems, item)
	}

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

	c.IndentedJSON(http.StatusOK, finalItems)
}

func AddItemToShoppingList(c *gin.Context) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var sl database.ShoppingListModel
	if err := c.ShouldBindUri(&sl); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list url not found"})
		return
	}

	shoppingListModel, err := sl.Read(db)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading shopping list"})
		return
	}

	shoppingListObj := shoppingListModel.(*database.ShoppingListModel)
	var item database.Item

	bindingErr := c.ShouldBind(&item)
	if bindingErr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error binding post request body to item"})
		return
	}
	
	shoppingListCRDT := shoppingList.DatabaseShoppingListToCRDT(shoppingListModel.(*database.ShoppingListModel))
	shoppingListCRDT.AddItem(item.Name, int(item.Quantity))
	newDB := shoppingListCRDT.ToDatabaseShoppingList(sl.Id)

	err = shoppingListModel.Update(db, newDB)
	if(err != nil){
		log.Print("Error writing to memory.")
		return
	}

	item.List = *shoppingListObj
	_, createErr := item.Create(db)
	if createErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating Item"})
		return
	}

	c.IndentedJSON(http.StatusOK, item)
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

	var sl database.Model

	if err := c.ShouldBindUri(&sl); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list not found"})
		return
	}

	shoppingListCRDT := shoppingList.DatabaseShoppingListToCRDT(sl.(*database.ShoppingListModel))
	shoppingListCRDT.DeleteItem(item.Name)
	newDB := shoppingListCRDT.ToDatabaseShoppingList(sl.(*database.ShoppingListModel).Id)

	err := sl.Update(db, newDB)
	if(err != nil){
		log.Print("Error deleting from memory.")
		return
	}

	itemModel.Delete(db)
	c.IndentedJSON(http.StatusOK, gin.H{"msg": "item deleted successfully"})
}

func Login(c *gin.Context) {
	if isLoggedIn(c) {
		c.IndentedJSON(http.StatusOK, "")
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
		_, createErr := user.Create(db)
		if createErr != nil {
			return "", createErr
		}
	}

	return usernameStr, nil
}

func isLoggedIn(c *gin.Context) bool {
	_, err := getUsernameFromCookie(c)
	fmt.Println(err)

	return (err == nil)
}

func SetDB(database *database.SQLiteRepository) {
	db = database
}

func GetDB() *database.SQLiteRepository {
	return db
}

func SetMessagesToSendChannel(ch chan messageStruct.MessageStruct) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("messagesToSend", ch)
        c.Next()
    }
}

func UpdateItemInShoppingList(c *gin.Context) {
    if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	var shoppingListt database.ShoppingListModel

	if err := c.ShouldBindUri(&shoppingListt); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list url not found"})
		return
	}

	shoppingListModel, err := shoppingListt.Read(db)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading shopping list"})
		return
	}

	shoppingListObj := shoppingListModel.(*database.ShoppingListModel)
	shoppingListCRDT := shoppingList.DatabaseShoppingListToCRDT(shoppingListObj)

	var item database.Item
	if err := c.ShouldBind(&item); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error binding item"})
		return
	}

	itemModel, err := item.Read(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "failed to read item"})
	}

	itemObj := itemModel.(*database.Item)
	updatedItemObj := database.Item{Name: itemObj.Name, Quantity: item.Quantity, Done: itemObj.Done, List: itemObj.List}
	itemObj.Update(db, &updatedItemObj)

	username, _ := getUsernameFromCookie(c)
	userList := database.UserList{ListID: shoppingListObj.Id, UserID: username}
	userListObj, _ := userList.Read(db)

	if userListObj != nil {
		updatedShoppingList := shoppingListCRDT.AlterItemQuantity(updatedItemObj.Name, int(updatedItemObj.Quantity))
		c.IndentedJSON(http.StatusOK, updatedShoppingList)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error..."})
}

func UploadList(c *gin.Context, connected bool) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	
	if !connected {
		return
	}


	var username, cookieErr = getUsernameFromCookie(c)
	if cookieErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error getting username from cookie"})
		return
	}

	var sl database.ShoppingListModel

	if err := c.ShouldBindUri(&sl); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list not found"})
		return
	}

	shoppingListModel, err := sl.Read(db)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "error reading shopping list"})
		return
	}

	shoppingListObj := shoppingListModel.(*database.ShoppingListModel)

	shoppingListCRDT := shoppingList.DatabaseShoppingListToCRDT(shoppingListObj)
	messageJSON := shoppingListCRDT.ConvertToMessageFormat(username, messageStruct.Write)
	message, convErr := messageStruct.JSONToMessage(messageJSON)

	if convErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error converting JSON to Message Struct"})
	}

	ch, ok := c.Get("messagesToSend")
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Channel not found"})
        return
    }

    messagesToSend := ch.(chan messageStruct.MessageStruct)

	messagesToSend <- message

	log.Printf("Sent write request to orchestrator for list %s.", message.ListURL)

	c.IndentedJSON(http.StatusOK, gin.H{"msg": "list uploaded successfully"})
}

func SetListsToAddChannel(ch chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("listsToAdd", ch)
		c.Next()
	}
}

func SetWriteListsToDatabaseChannel(ch chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("writeListsToDatabase", ch)
		c.Next()
	}
}


func FetchList(c *gin.Context, connected bool) {
	if !isLoggedIn(c) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"msg": "user must be logged in"})
		return
	}

	if !connected {
		return
	}

	var sl database.ShoppingListModel

	if err := c.ShouldBindUri(&sl); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "list not found"})
		return
	}

	shoppingListModel, _ := sl.Read(db)

	var username, cookieErr = getUsernameFromCookie(c)
	if cookieErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error getting username from cookie"})
		return
	}


	ch, ok := c.Get("listsToAdd")
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Channel not found"})
        return
    }

    listsToAdd := ch.(chan string)


	ch, ok = c.Get("messagesToSend")
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Channel not found"})
        return
    }

    messagesToSend := ch.(chan messageStruct.MessageStruct)


	ch, ok = c.Get("writeListsToDatabase")
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Channel not found"})
        return
    }

    writeListsToDatabase := ch.(chan string)


	shoppingListCRDT := shoppingList.DatabaseShoppingListToCRDT(shoppingListModel.(*database.ShoppingListModel))
	messageJSON := shoppingListCRDT.ConvertToMessageFormat(username, messageStruct.Read)
	message, _ := messageStruct.JSONToMessage(messageJSON)
	
	listsToAdd <- message.ListURL
	
	messagesToSend <- message
	
	log.Printf("Sent read request to orchestrator for list %s.", message.ListURL)

	writeListsToDatabase <- message.ListURL

	c.IndentedJSON(http.StatusOK, gin.H{"msg": "list fetched successfully"})
}
