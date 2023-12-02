package main

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"sdle/server/database"
	"sdle/server/CRDT/lexCounter"
	shoppingList "sdle/server/CRDT/shoppingList"
)

func LexExample(){
	x := lexCounter.Create[string, int]("a")
	y := lexCounter.Create[string, int]("b")

	x.Inc(4)
	x.Dec(1)

	y.Inc(2)

	fmt.Println(x.GetValue())
	fmt.Println(y.GetValue())

	x.Join(y)

	lexCounter.Print(x, y)
}


func ShopListExample() {
	
	shopList1 := shoppingList.Create("My List 1")
	shopList2 := shoppingList.Create("My List 2")

	fmt.Println(shopList1.GetURL())

	shopList1.AddItem("apple", 3)
	shopList1.AddItem("rice", 5)
	
	shopList2.AddItem("pear", 2)
	shopList2.AddItem("rice", 3)
	shopList2.BuyItem("rice")

	fmt.Println("\nShop List 1")
	fmt.Println(shopList1.JSON())
	fmt.Println("\n")

	fmt.Println("\nShop List 2")
	fmt.Println(shopList2.JSON())
	fmt.Println("\n")


	shopList1.JoinShoppingList(shopList2)

	fmt.Println("\nShop List 1 after merging with Shop List 2")
	fmt.Println(shopList1.JSON())
	fmt.Println("\n")

	shopList1.JoinShoppingList(shopList2)

	fmt.Println("\nShop List 1 after merging with Shop List 2 (again)")
	fmt.Println(shopList1.JSON())
	fmt.Println("\n")

}

func main() {
	fmt.Println("Hello from server")
	//SetExample()
	//LexExample()
	ShopListExample()

	const filename = "server.db"
	db, err := sql.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}

	sqliteRepository := database.NewSQLiteRepository(db)
	createError := sqliteRepository.CreateTables()
	if createError != nil {
		fmt.Println(createError.Error())
		return
	}

	seedError := sqliteRepository.Seed()
	if seedError != nil {
		fmt.Println(seedError.Error())
	}
}
