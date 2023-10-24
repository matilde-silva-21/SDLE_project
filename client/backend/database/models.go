package database

import (
	"fmt"
)

type Model interface {
	CreateTable(r *SQLiteRepository) error
	Create(r *SQLiteRepository) (Model, error)
	Delete(r *SQLiteRepository) error
	Update(r *SQLiteRepository, updated Model) error
	Read(r *SQLiteRepository) (Model, error)
	ReadAll(r *SQLiteRepository) ([]Model, error)
}

type ShoppingList struct {
	Id int64 `json:"id"`
}

type User struct {
	Username string `json:"username"`
}

type Item struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type UserList struct {
	ListID int64 `json:"listID"`
	UserID int64 `json:"userID"`
}

type ListItem struct {
	ListID int64 `json:"listID"`
	ItemID int64 `json:"itemID"`
}

// ITEM MODELS METHODS

func (item *Item) CreateTable(r *SQLiteRepository) error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS Item (Id INTEGER PRIMARY KEY, Name TEXT, Done INTEGER)")

	if err != nil {
		return err
	}
	return nil
}

func (item *Item) Create(r *SQLiteRepository) (Model, error) {
	res, err := r.db.Exec("INSERT INTO Item(Name, Done) VALUES (?, ?)", &item.Name, &item.Done)

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
	res, err := r.db.Exec("UPDATE Item SET Name = (?), Done = (?) WHERE Id = (?)", updatedItem.Name, updatedItem.Done, &item.Id)

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
	res := r.db.QueryRow("SELECT * FROM Item WHERE Id = (?)", &item.Id)

	var updated Item
	if err := res.Scan(&updated.Id, &updated.Name, &updated.Done); err != nil {
		return nil, err
	}

	fmt.Println(updated)

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

// SHOPPING LIST MODEL METHODS

func (list *ShoppingList) CreateTable(r *SQLiteRepository) error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS ShoppingList (Id INTEGER PRIMARY KEY)")

	if err != nil {
		return err
	}
	return nil
}

func (list *ShoppingList) Create(r *SQLiteRepository) (Model, error) {
	_, err := r.db.Exec("INSERT INTO ShoppingList(Id) VALUES (?)", &list.Id)

	if err != nil {
		return nil, err
	}

	return list, nil
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
	res := r.db.QueryRow("SELECT * FROM ShoppingList WHERE Id = (?)", list.Id)

	var updated ShoppingList
	if err := res.Scan(&updated.Id); err != nil {
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

// USER LIST MODEL METHODS

func (userList *UserList) CreateTable(r *SQLiteRepository) error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS UserList (ListId INTEGER REFERENCES ShoppingList, UserId INTEGER REFERENCES User, PRIMARY KEY(ListId, UserId))")

	if err != nil {
		return err
	}
	return nil
}

func (userList *UserList) Create(r *SQLiteRepository) (Model, error) {
	_, err := r.db.Exec("INSERT INTO UserList(ListId, UserId) VALUES (?)", &userList.ListID, &userList.UserID)

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
	res := r.db.QueryRow("SELECT * FROM UserList WHERE ListId = (?) AND UserId = (?)",  userList.ListID, userList.ListID)

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

// LIST ITEM MODEL METHODS

func (listItem *ListItem) CreateTable(r *SQLiteRepository) error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS ListItem (ListId INTEGER REFERENCES ShoppingList, ItemId INTEGER REFERENCES Item, PRIMARY KEY(ListId, ItemId))")

	if err != nil {
		return err
	}
	return nil
}

func (listItem *ListItem) Create(r *SQLiteRepository) (Model, error) {
	_, err := r.db.Exec("INSERT INTO ListItem(ListId, ItemId) VALUES (?)", &listItem.ListID, &listItem.ItemID)

	if err != nil {
		return nil, err
	}

	return listItem, nil
}

func (listItem *ListItem) Delete(r *SQLiteRepository) error {
	_, err := r.db.Exec("DELETE FROM ListItem WHERE ListId = (?) AND ItemId = (?)", &listItem.ListID, &listItem.ItemID)

	if err != nil {
		return err
	}

	return nil
}

func (listItem *ListItem) Update(r *SQLiteRepository, updated Model) error {
	return nil
}

func (listItem *ListItem) Read(r *SQLiteRepository) (Model, error) {
	res := r.db.QueryRow("SELECT * FROM ListItem WHERE ListId = (?) AND ItemId = (?)", listItem.ListID, listItem.ItemID)

	var updated ListItem
	if err := res.Scan(&updated.ListID, &updated.ItemID); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (listItem *ListItem) ReadAll(r *SQLiteRepository) ([]Model, error) {
	rows, err := r.db.Query("SELECT * FROM ListItem")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lists []Model

	for rows.Next() {
		var li ListItem
		if err := rows.Scan(&li.ListID, &li.ItemID); err != nil {
			return nil, err
		}
		lists = append(lists, &li)
	}

	return lists, nil
}
