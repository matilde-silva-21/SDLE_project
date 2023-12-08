package database

import (
	"fmt"	
);

type Model interface {
	CreateTable(r *SQLiteRepository) error
	Create(r *SQLiteRepository) (Model, error)
	Delete(r *SQLiteRepository) error
	Update(r *SQLiteRepository, updated Model) error
	Read(r *SQLiteRepository) (Model, error)
	ReadAll(r *SQLiteRepository) ([]Model, error)
}

type ShoppingList struct {
	Id   int64  `json:"id"`
	Name string `json:"name" form:"listName"`
	Url  string `json:"url" uri:"url"`
	List string `json:"list"`
	State string `json:"state"`
}

type User struct {
	Username string `json:"username" form:"username"`
}

type Item struct {
	Id   int64        `json:"id" uri:"id"`
	Name string       `json:"name" form:"itemName"`
	Done bool         `json:"done" form:"itemDone"`
	Quantity int64	  `json:"quantity"`
	List ShoppingList `json:"list"`
}

type UserList struct {
	ListID int64  `json:"listID"`
	UserID string `json:"userID"`
}

// ITEM MODELS METHODS

func (item *Item) CreateTable(r *SQLiteRepository) error {
	r.db.Exec("DROP TABLE IF EXISTS Item")
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS Item (Id INTEGER PRIMARY KEY, Name TEXT, Done INTEGER, Quantity INTEGER, List TEXT REFERENCES ShoppingList)")

	if err != nil {
		return err
	}
	return nil
}

func (item *Item) Create(r *SQLiteRepository) (Model, error) {
	res, err := r.db.Exec("INSERT INTO Item(Name, Done, Quantity, List) VALUES (?, ?, ?, ?)", &item.Name, &item.Done, &item.Quantity, &item.List.Id)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	item.Id = id

	return item, nil
}

func (item *Item) Delete(r *SQLiteRepository) error {
	_, err := r.db.Exec("DELETE FROM Item WHERE Id = (?)", &item.Id)

	if err != nil {
		return err
	}

	return nil
}

func (item *Item) Update(r *SQLiteRepository, updated Model) error {
	updatedItem := updated.(*Item)
	res, err := r.db.Exec("UPDATE Item SET Name = (?), Done = (?), Quantity = (?), List = (?) WHERE Id = (?)", 
		updatedItem.Name, updatedItem.Done, updatedItem.Quantity, updatedItem.List.Id, &item.Id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil || rows == 0 {
		return err
	}

	return nil
}

func (item *Item) Read(r *SQLiteRepository) (Model, error) {
	res := r.db.QueryRow("SELECT * FROM Item WHERE Id = (?) OR Name = (?)", &item.Id, &item.Name)

	var updated Item
	if err := res.Scan(&updated.Id, &updated.Name, &updated.Done, &updated.Quantity, &updated.List.Id); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (item *Item) ReadAll(r *SQLiteRepository) ([]Model, error) {
	rows, err := r.db.Query("SELECT * FROM Item")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Model

	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.Id, &i.Name, &i.Done); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}

	return items, nil
}

// USER MODEL METHODS

func (user *User) CreateTable(r *SQLiteRepository) error {
	r.db.Exec("DROP TABLE IF EXISTS User")
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS User (Username TEXT PRIMARY KEY)")

	if err != nil {
		return err
	}
	return nil
}

func (user *User) Create(r *SQLiteRepository) (Model, error) {
	_, err := r.db.Exec("INSERT INTO User(username) VALUES (?)", &user.Username)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) Delete(r *SQLiteRepository) error {
	_, err := r.db.Exec("DELETE FROM User WHERE Username = (?)", user.Username)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) Update(r *SQLiteRepository, updated Model) error {
	return nil
}

func (user *User) Read(r *SQLiteRepository) (Model, error) {
	res := r.db.QueryRow("SELECT * FROM User WHERE Username = (?)", user.Username)

	var updated User
	if err := res.Scan(&updated.Username); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (user *User) ReadAll(r *SQLiteRepository) ([]Model, error) {
	rows, err := r.db.Query("SELECT * FROM User")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []Model

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (user *User) ReadUserLists(r *SQLiteRepository) ([]ShoppingList, error) {
	fmt.Println(user.Username)
	rows, err := r.db.Query("SELECT ShoppingList.Id, ShoppingList.Name, ShoppingList.Url FROM ShoppingList JOIN UserList ON UserList.ListId = ShoppingList.Id WHERE UserList.UserId = ?", user.Username)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userLists []ShoppingList

	for rows.Next() {
		var ul ShoppingList

		if scanErr := rows.Scan(&ul.Id, &ul.Name, &ul.Url); scanErr != nil {
			return nil, scanErr
		}

		userLists = append(userLists, ul)
	}

	return userLists, nil
}

// SHOPPING LIST MODEL METHODS

func (list *ShoppingList) CreateTable(r *SQLiteRepository) error {
	r.db.Exec("DROP TABLE IF EXISTS ShoppingList")
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS ShoppingList (Id INTEGER PRIMARY KEY, Name TEXT, Url TEXT UNIQUE, List TEXT, State TEXT)")

	if err != nil {
		return err
	}
	return nil
}

func (list *ShoppingList) Create(r *SQLiteRepository) (Model, error) {
	_, err := r.db.Exec("INSERT INTO ShoppingList(Name, Url, List, State) VALUES (?, ?, ? ,?)", &list.Name, &list.Url, &list.List, &list.State)

	if err != nil {
		return nil, err
	}
	newList, _ := list.Read(r)

	return newList, nil
}

func (list *ShoppingList) Delete(r *SQLiteRepository) error {
	_, err := r.db.Exec("DELETE FROM ShoppingList WHERE Id = (?)", list.Id)

	if err != nil {
		return err
	}

	return nil
}

func (list *ShoppingList) Update(r *SQLiteRepository, updated Model) error {
	return nil
}

func (list *ShoppingList) Read(r *SQLiteRepository) (Model, error) {
	res := r.db.QueryRow("SELECT * FROM ShoppingList WHERE Url = ?", &list.Url)

	var updated ShoppingList
	if err := res.Scan(&updated.Id, &updated.Name, &updated.Url, &updated.List, &updated.State); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (list *ShoppingList) ReadAll(r *SQLiteRepository) ([]Model, error) {
	rows, err := r.db.Query("SELECT * FROM ShoppingList")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lists []Model

	for rows.Next() {
		var l ShoppingList
		if err := rows.Scan(&l.Id); err != nil {
			return nil, err
		}
		lists = append(lists, &l)
	}

	return lists, nil
}

func (list *ShoppingList) GetShoppingListItems(r *SQLiteRepository) ([]Item, error) {
	rows, err := r.db.Query("SELECT * FROM Item WHERE List = (?)", &list.Id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Item

	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.Id, &i.Name, &i.Done, &i.Quantity, &i.List.Id); err != nil {
			return nil, err
		}
		i.List.Name = list.Name
		i.List.Url = list.Url
		items = append(items, i)
	}

	return items, nil
}

func GetIDByURL(r *SQLiteRepository, url string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT Id FROM ShoppingList WHERE Url = ?", url).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return a custom error indicating that no matching URL was found
			return 0, fmt.Errorf("ShoppingList with URL '%s' not found", url)
		}
		// Return other errors as is
		return 0, err
	}

	return id, nil
}

// USER LIST MODEL METHODS

func (userList *UserList) CreateTable(r *SQLiteRepository) error {
	r.db.Exec("DROP TABLE IF EXISTS UserList")
	_, err := r.db.Exec(`
	CREATE TABLE IF NOT EXISTS UserList (
		ListId INTEGER REFERENCES ShoppingList ON DELETE CASCADE ON UPDATE CASCADE,
		UserId INTEGER REFERENCES User ON DELETE CASCADE ON UPDATE CASCADE, 
		PRIMARY KEY(ListId, UserId)
	)`)

	if err != nil {
		return err
	}
	return nil
}

func (userList *UserList) Create(r *SQLiteRepository) (Model, error) {
	fmt.Println(userList)
	_, err := r.db.Exec("INSERT INTO UserList(ListId, UserId) VALUES (?, ?)", &userList.ListID, &userList.UserID)

	if err != nil {
		return nil, err
	}

	return userList, nil
}

func (userList *UserList) Delete(r *SQLiteRepository) error {
	_, err := r.db.Exec("DELETE FROM UserList WHERE ListId = (?) AND UserId = (?)", &userList.ListID, &userList.UserID)

	if err != nil {
		return err
	}

	return nil
}

func (userList *UserList) Update(r *SQLiteRepository, updated Model) error {
	return nil
}

func (userList *UserList) Read(r *SQLiteRepository) (Model, error) {
	res := r.db.QueryRow("SELECT * FROM UserList WHERE ListId = (?) AND UserId = (?)", &userList.ListID, &userList.UserID)

	var updated UserList
	if err := res.Scan(&updated.ListID, &updated.UserID); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (userList *UserList) ReadAll(r *SQLiteRepository) ([]Model, error) {
	rows, err := r.db.Query("SELECT * FROM UserList")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lists []Model

	for rows.Next() {
		var ul UserList
		if err := rows.Scan(&ul.ListID, &ul.UserID); err != nil {
			return nil, err
		}
		lists = append(lists, &ul)
	}

	return lists, nil
}
