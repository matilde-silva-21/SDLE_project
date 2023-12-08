package database

import (
	"database/sql"
	"fmt"
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
		&ShoppingListModel{},
		&Item{},
		&User{},
		&UserList{},
	}

	for i, v := range models {
		err := v.CreateTable(r)
		fmt.Printf("created model %d %v\n", i, v)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteRepository) Seed() error {
	sl1 := ShoppingListModel{Id: 1, Name: "My Shopping List 1", Url: "testurl", List: "", State: ""}
	sl2 := ShoppingListModel{Id: 2, Name: "My Shopping List 1", Url: "test", List: "", State: ""}

	var seedItems = []Model{
		&sl1,
		&sl2,
		&Item{Id: 1, Name: "Bread", Done: false, Quantity: 3, List: sl1},
		&Item{Id: 2, Name: "Cheese", Done: true, Quantity: 1, List: sl2},
	}

	for _, v := range seedItems {
		_, err := v.Create(r)
		fmt.Printf("seeded model %v\n", v)

		if err != nil {
			return err
		}
	}

	return nil
}
