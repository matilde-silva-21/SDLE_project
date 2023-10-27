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
	sl1 := ShoppingList{Id: 1, Name: "My Shopping List 1",  Url: "testurl"}
	sl2 := ShoppingList{Id: 2, Name: "My Shopping List 1", Url: "test"}  
	
	var seedItems = []Model{
		&sl1,
		&sl2,
		&Item{Id: 1, Name: "Bread", Done: false, List: sl1},
		&Item{Id: 2, Name: "Cheese", Done: true, List: sl1},
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

