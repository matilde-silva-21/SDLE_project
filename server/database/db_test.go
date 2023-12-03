package database

import (
	"log"
	"testing"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func TestShoppingListModel(t *testing.T) {
	const filename = "test.db"
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteRepository := NewSQLiteRepository(db)

	resetDBErr := resetDatabase(sqliteRepository)
	if resetDBErr != nil {
		t.Errorf("ResetDatabase error: %v", resetDBErr)
	}

	// Create ShoppingList instance
	newList := &ShoppingList{
		Name: "TestList",
		Url:  "test-url",
	}

	// Test CreateTable
	createTableErr := newList.CreateTable(sqliteRepository)
	if createTableErr != nil {
		t.Errorf("CreateTable error: %v", createTableErr)
	}

	// Test Create
	createdList, createErr := newList.Create(sqliteRepository)
	if createErr != nil {
		t.Errorf("Create error: %v", createErr)
	}

	// Test Read
	readList, readErr := createdList.Read(sqliteRepository)
	if readErr != nil {
		t.Errorf("Read error: %v", readErr)
	}

	// Compare the created list with the read list
	if readList.(*ShoppingList).Name != newList.Name || readList.(*ShoppingList).Url != newList.Url {
		t.Errorf("Read result does not match the created ShoppingList")
	}

	// Test Update
	updatedList := &ShoppingList{
		Id:   readList.(*ShoppingList).Id,
		Name: "UpdatedTestList",
		Url:  "updated-url",
	}

	updateErr := updatedList.Update(sqliteRepository, updatedList)
	if updateErr != nil {
		t.Errorf("Update error: %v", updateErr)
	}

	// Test ReadAll
	allLists, readAllErr := newList.ReadAll(sqliteRepository)
	if readAllErr != nil {
		t.Errorf("ReadAll error: %v", readAllErr)
	}

	// Check if the updated list is in the list of all lists
	found := false
	for _, l := range allLists {
		if l.(*ShoppingList).Id == updatedList.Id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Updated list not found in ReadAll result")
	}

	// Test Delete
	deleteErr := updatedList.Delete(sqliteRepository)
	if deleteErr != nil {
		t.Errorf("Delete error: %v", deleteErr)
	}

	// Test Read after deletion
	deletedList, readDeletedErr := updatedList.Read(sqliteRepository)
	if readDeletedErr == nil {
		t.Errorf("Read after deletion should return an error, but got nil")
	}
	
	// Ensure that the returned ShoppingList is nil after deletion
	if deletedList != nil {
		t.Errorf("Read after deletion should return nil, but got: %v", deletedList)
	}
}

// TestGetShoppingListItems tests the GetShoppingListItems method
func TestGetShoppingListItems(t *testing.T) {
	const filename = "test.db"
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteRepository := NewSQLiteRepository(db)

	resetDBErr := resetDatabase(sqliteRepository)
	if resetDBErr != nil {
		t.Errorf("ResetDatabase error: %v", resetDBErr)
	}

	// Create a ShoppingList instance
	newList := &ShoppingList{
		Name: "TestList",
		Url:  "test-url",
	}

	// Test CreateTable
	createTableListErr := newList.CreateTable(sqliteRepository)
		if createTableListErr != nil {
			t.Errorf("CreateTable error: %v", createTableListErr)
		}

	// Create the ShoppingList in the database
	createdList, createErr := newList.Create(sqliteRepository)
	if createErr != nil {
		t.Errorf("Create error: %v", createErr)
	}

	// Create an Item associated with the ShoppingList
	newItem := &Item{
		Name:     "TestItem",
		Done:     false,
		Quantity: 1,
		List:     *createdList.(*ShoppingList),
	}

	// Test CreateTable
	createTableItemErr := newItem.CreateTable(sqliteRepository)
		if createTableItemErr != nil {
			t.Errorf("CreateTable error: %v", createTableItemErr)
		}
	

	// Create the Item in the database
	createdItem, createItemErr := newItem.Create(sqliteRepository)
	if createItemErr != nil {
		t.Errorf("CreateItem error: %v", createItemErr)
	}

	// Call GetShoppingListItems
	items, getItemsErr := createdList.(*ShoppingList).GetShoppingListItems(sqliteRepository)
	if getItemsErr != nil {
		t.Errorf("GetShoppingListItems error: %v", getItemsErr)
	}

	// Check if the created item is in the list of items
	found := false
	for _, item := range items {
		if item.Id == createdItem.(*Item).Id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Created item not found in GetShoppingListItems result")
	}

	// Clean up: Delete the ShoppingList and associated items
	deleteListErr := createdList.(*ShoppingList).Delete(sqliteRepository)
	if deleteListErr != nil {
		t.Errorf("DeleteList error: %v", deleteListErr)
	}
}

// User model test function
func TestUserModel(t *testing.T) {
	const filename = "test.db"
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteRepository := NewSQLiteRepository(db)

	// Create User instance
	newUser := User{
		Username: "TestUser",
	}

	// Test CreateTable
	createTableErr := newUser.CreateTable(sqliteRepository)
	if createTableErr != nil {
		t.Errorf("CreateTable error: %v", createTableErr)
	}

	// Test Create
	createdUser, createErr := newUser.Create(sqliteRepository)
	if createErr != nil {
		t.Errorf("Create error: %v", createErr)
	}

	// Test Read
	readUser, readErr := createdUser.Read(sqliteRepository)
	if readErr != nil {
		t.Errorf("Read error: %v", readErr)
	}

	// Compare the created user with the read user
	if readUser.(*User).Username != newUser.Username {
		t.Errorf("Read result does not match the created User")
	}

	// Call ReadAllUsers
	users, readAllErr := newUser.ReadAll(sqliteRepository)
	if readAllErr != nil {
		t.Errorf("ReadAllUsers error: %v", readAllErr)
	}

	// Check if the created user is in the list of users
	found := false
	for _, user := range users {
		if user.(*User).Username == createdUser.(*User).Username {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Created user not found in ReadAllUsers result")
	}

	// Create a ShoppingList associated with the User
	newList := &ShoppingList{
		Name: "TestList",
		Url:  "test-url",
	}

	// Test CreateTable
	createTableListErr := newList.CreateTable(sqliteRepository)
		if createTableListErr != nil {
			t.Errorf("CreateTable error: %v", createTableListErr)
		}

	// Create the ShoppingList in the database
	createdList, createListErr := newList.Create(sqliteRepository)
	if createListErr != nil {
		t.Errorf("CreateList error: %v", createListErr)
	}

	// Create a UserList association
	userList := &UserList{
		ListID: createdList.(*ShoppingList).Id,
		UserID: createdUser.(*User).Username,
	}

	// Test CreateTable
	createTableUserListErr := userList.CreateTable(sqliteRepository)
		if createTableUserListErr != nil {
			t.Errorf("CreateTable error: %v", createTableUserListErr)
		}

	// Create the UserList in the database
	_, createUserListErr := userList.Create(sqliteRepository)
	if createUserListErr != nil {
		t.Errorf("CreateUserList error: %v", createUserListErr)
	}

	// Call ReadUserLists
	userLists, readUserListsErr := createdUser.(*User).ReadUserLists(sqliteRepository)
	if readUserListsErr != nil {
		t.Errorf("ReadUserLists error: %v", readUserListsErr)
	}

	// Check if the created ShoppingList is in the list of user lists
	found = false
	for _, ul := range userLists {
		if ul.Id == createdList.(*ShoppingList).Id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Created ShoppingList not found in ReadUserLists result")
	}

	// Test Delete
	deleteErr := createdUser.Delete(sqliteRepository)
	if deleteErr != nil {
		t.Errorf("Delete error: %v", deleteErr)
	}

	// Test Read after deletion
	deletedUser, readDeletedErr := createdUser.Read(sqliteRepository)
	if readDeletedErr == nil {
		t.Errorf("Read after deletion should return an error, but got nil")
	}

	// Ensure that the returned User is nil after deletion
	if deletedUser != nil {
		t.Errorf("Read after deletion should return nil, but got: %v", deletedUser)
	}

}

// TestItemMethods tests the methods of the Item model
func TestItemMethods(t *testing.T) {
	const filename = "test.db"
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteRepository := NewSQLiteRepository(db)

	// Cleanup and reset the database
	resetDBErr := resetDatabase(sqliteRepository)
	if resetDBErr != nil {
		t.Errorf("ResetDatabase error: %v", resetDBErr)
	}

	// Create an Item instance
	newItem := &Item{
		Name:     "TestItem",
		Done:     false,
		Quantity: 10,
		List: ShoppingList{
			Name: "TestList",
			Url:  "test-url",
		},
	}

	// Call CreateTable to create the Item table
	createTableErr := newItem.CreateTable(sqliteRepository)
	if createTableErr != nil {
		t.Errorf("CreateTable error: %v", createTableErr)
	}

	// Call Create to insert the Item into the database
	createdItem, createErr := newItem.Create(sqliteRepository)
	if createErr != nil {
		t.Errorf("CreateItem error: %v", createErr)
	}

	// Call Read to retrieve the Item from the database
	readItem, readErr := createdItem.(*Item).Read(sqliteRepository)
	if readErr != nil {
		t.Errorf("ReadItem error: %v", readErr)
	}

	// Verify that the created Item matches the retrieved Item
	if createdItem.(*Item).Id != readItem.(*Item).Id {
		t.Errorf("Created Item ID does not match Read Item ID")
	}

	// Call Update to modify the Item in the database
	createdItem.(*Item).Name = "UpdatedItem"
	updateErr := createdItem.(*Item).Update(sqliteRepository, createdItem)
	if updateErr != nil {
		t.Errorf("UpdateItem error: %v", updateErr)
	}

	// Call Read to retrieve the updated Item from the database
	updatedItem, readUpdatedErr := createdItem.(*Item).Read(sqliteRepository)
	if readUpdatedErr != nil {
		t.Errorf("ReadUpdatedItem error: %v", readUpdatedErr)
	}

	// Verify that the Name of the updated Item matches the modification
	if updatedItem.(*Item).Name != "UpdatedItem" {
		t.Errorf("Updated Item Name does not match expected value")
	}

	// Call ReadAll to retrieve all Items from the database
	items, readAllErr := newItem.ReadAll(sqliteRepository)
	if readAllErr != nil {
		t.Errorf("ReadAllItems error: %v", readAllErr)
	}

	// Verify that the created Item is in the list of Items
	found := false
	for _, item := range items {
		if item.(*Item).Id == createdItem.(*Item).Id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Created Item not found in ReadAllItems result")
	}

	// Call Delete to remove the Item from the database
	deleteErr := createdItem.(*Item).Delete(sqliteRepository)
	if deleteErr != nil {
		t.Errorf("DeleteItem error: %v", deleteErr)
	}

	// Verify that the deleted Item is not in the list of Items
	deletedItem, readDeletedErr := createdItem.(*Item).Read(sqliteRepository)
	if readDeletedErr == nil || deletedItem != nil {
		t.Errorf("Deleted Item found in ReadDeletedItem result")
	}

}

// TestUserListMethods tests the methods of the UserList model
func TestUserListMethods(t *testing.T) {
	const filename = "test.db"
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteRepository := NewSQLiteRepository(db)

	// Cleanup and reset the database
	resetDBErr := resetDatabase(sqliteRepository)
	if resetDBErr != nil {
		t.Errorf("ResetDatabase error: %v", resetDBErr)
	}

	// Create a User and a ShoppingList for the UserList
	user := &User{
		Username: "TestUser",
	}

	// Test CreateTable
	createTableUserErr := user.CreateTable(sqliteRepository)
		if createTableUserErr != nil {
			t.Errorf("CreateTable error: %v", createTableUserErr)
		}


	createdUser, createUserErr := user.Create(sqliteRepository)
	if createUserErr != nil {
		t.Errorf("CreateUser error: %v", createUserErr)
	}
	defer createdUser.Delete(sqliteRepository)

	shoppingList := &ShoppingList{
		Name: "TestList",
		Url:  "test-url",
	}

	// Test CreateTable
	createTableListErr := shoppingList.CreateTable(sqliteRepository)
		if createTableListErr != nil {
			t.Errorf("CreateTable error: %v", createTableListErr)
		}

	createdList, createListErr := shoppingList.Create(sqliteRepository)
	if createListErr != nil {
		t.Errorf("CreateList error: %v", createListErr)
	}
	defer createdList.Delete(sqliteRepository)

	// Create a UserList instance
	newUserList := &UserList{
		ListID: createdList.(*ShoppingList).Id,
		UserID: createdUser.(*User).Username,
	}

	// Test CreateTable
	createTableUserListErr := newUserList.CreateTable(sqliteRepository)
	if createTableUserListErr != nil {
		t.Errorf("CreateTable error: %v", createTableUserListErr)
	}


	// Call Create to insert the UserList into the database
	createdUserList, createErr := newUserList.Create(sqliteRepository)
	if createErr != nil {
		t.Errorf("CreateUserList error: %v", createErr)
	}

	// Call Read to retrieve the UserList from the database
	readUserList, readErr := createdUserList.(*UserList).Read(sqliteRepository)
	if readErr != nil {
		t.Errorf("ReadUserList error: %v", readErr)
	}

	// Verify that the created UserList matches the retrieved UserList
	if createdUserList.(*UserList).ListID != readUserList.(*UserList).ListID ||
		createdUserList.(*UserList).UserID != readUserList.(*UserList).UserID {
		t.Errorf("Created UserList does not match Read UserList")
	}

	// Call ReadAll to retrieve all UserLists from the database
	userLists, readAllErr := newUserList.ReadAll(sqliteRepository)
	if readAllErr != nil {
		t.Errorf("ReadAllUserLists error: %v", readAllErr)
	}

	// Verify that the created UserList is in the list of UserLists
	found := false
	for _, userList := range userLists {
		if userList.(*UserList).ListID == createdUserList.(*UserList).ListID &&
			userList.(*UserList).UserID == createdUserList.(*UserList).UserID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Created UserList not found in ReadAllUserLists result")
	}

	// Call Delete to remove the UserList from the database
	deleteErr := createdUserList.(*UserList).Delete(sqliteRepository)
	if deleteErr != nil {
		t.Errorf("DeleteUserList error: %v", deleteErr)
	}

	// Verify that the deleted UserList is not in the list of UserLists
	deletedUserList, readDeletedErr := createdUserList.(*UserList).Read(sqliteRepository)
	if readDeletedErr == nil || deletedUserList != nil {
		t.Errorf("Deleted UserList found in ReadDeletedUserList result")
	}
}


func resetDatabase(repository *SQLiteRepository) error {
	_, err := repository.db.Exec("DROP TABLE IF EXISTS User")
	if err != nil {
		return err
	}

	_, err = repository.db.Exec("DROP TABLE IF EXISTS ShoppingList")
	if err != nil {
		return err
	}

	_, err = repository.db.Exec("DROP TABLE IF EXISTS Item")
	if err != nil {
		return err
	}

	_, err = repository.db.Exec("DROP TABLE IF EXISTS UserList")
	return err
}