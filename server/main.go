package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sdle/server/database"
	"sdle/server/serverCommunicator"
	"sdle/server/utils/CRDT/lexCounter"
	shoppingList "sdle/server/utils/CRDT/shoppingList"
	"sdle/server/utils/messageStruct"

	_ "github.com/mattn/go-sqlite3"
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

	fmt.Println("\nShop List 2")
	fmt.Println(shopList2.JSON())


	shopList1.JoinShoppingList(shopList2)

	fmt.Println("\nShop List 1 after merging with Shop List 2")
	fmt.Println(shopList1.JSON())

	shopList1.JoinShoppingList(shopList2)

	fmt.Println("\nShop List 1 after merging with Shop List 2 (again)")
	fmt.Println(shopList1.JSON())

	messageFormat := shopList1.ConvertToMessageFormat("john.doe", messageStruct.Write)


	fmt.Println("\n", string(messageFormat))


	fmt.Println("\n", shoppingList.MessageByteToCRDT(messageFormat))
	shoplistCopy := shoppingList.CreateFromStrings(shopList1.GetURL(), shopList1.GetListName(), shopList1.ListFormatForDatabase(), shopList1.StateFormatForDatabase())

	fmt.Println("\n", shoplistCopy)
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Not enough arguments.\nUsage: go run main.go <orchestrator_address> <backup_orchestrator_address>.")
		return
	}

	fmt.Println("Hello from server")

	db_name := os.Args[3]
	
	filename := fmt.Sprintf("./dbs/server-%s.db", db_name)
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

	serverCommunicator.StartServerCommunication(os.Args[1], os.Args[2], sqliteRepository)
}
