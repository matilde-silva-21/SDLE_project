package database

import (
	"database/sql"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) CreateTables() error {
	var models = []Model{
		&Item{},
		&User{},
		&ShoppingList{},
		&UserList{},
		&ListItem{},
	}

	for _, v := range models {
		err := v.CreateTable(r)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteRepository) Seed() error {

	var seedItems = []Model{
		&Item{Id: 1, Name: "Bread", Done: false},
		&Item{Id: 2, Name: "Cheese", Done: true},
		&User{Username: "user1"},
		&ShoppingList{Id: 1},
	}

	for _, v := range seedItems {
		_, err := v.Create(r)

		if err != nil {
			return err
		}
	}

	return nil
}

