package api

import (
	"encoding/base64"
	"fmt"
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

func GetShoppingLists(c *gin.Context) {
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
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
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
	
	items, err := shoppingListObj.GetShoppingListItems(db)
	
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error retrieving shopping list items"})
		return
	}

	username, _ := getUsernameFromCookie(c)
	userList := database.UserList{ListID: shoppingListObj.Id, UserID: username}
	_, createUserListErr := userList.Create(db)

	if createUserListErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "error creating UserList entry"})
		return
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
		fmt.Println(err.Error())
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
	c.SetCookie("session", cookie, 0, "/", "localhost:8080", false, false)
}

func LoginPage(c *gin.Context) {
	cookie, err := c.Cookie("session")

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"msg": "cookie not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"msg": cookie})
}

func getUsernameFromCookie(c *gin.Context) (string, error) {
	cookie, cookieErr := c.Cookie("session")
	if cookieErr != nil {
		fmt.Println(cookieErr.Error())
		return "", cookieErr
	}

	username, decodeErr := base64.StdEncoding.DecodeString(cookie)
	if decodeErr != nil {
		fmt.Println(decodeErr.Error())
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
