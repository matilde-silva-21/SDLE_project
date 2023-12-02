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
	var models = []Model {
		&ShoppingList{},
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
	sl1 := ShoppingList{Id: 1, Name: "Shopping List 1",  Url: "test1"}
	
	var seedItems = []Model{
		&sl1,
		&Item{Id: 1, Name: "apple", Done: false, Quantity: 3, List: sl1},
		&Item{Id: 1, Name: "rice", Done: false, Quantity: 5, List: sl1},
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

