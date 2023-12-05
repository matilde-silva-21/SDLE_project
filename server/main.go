package main

import (
	"fmt"
	//"sdle/server/orchestrator"
	"sdle/server/utils/messageStruct"
	"sdle/server/utils/CRDT/lexCounter"
	shoppingList "sdle/server/utils/CRDT/shoppingList"
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

	messageFormat := shopList1.ConvertToMessageFormat("john.doe", messageStruct.Add)

	fmt.Println("\n", string(messageFormat))

	fmt.Println("\n", shoppingList.MessageFormatToCRDT(messageFormat))
	shoplistCopy := shoppingList.CreateFromStrings(shopList1.GetURL(), shopList1.GetListName(), shopList1.ListFormatForDatabase(), shopList1.StateFormatForDatabase())

	fmt.Println("\n", shoplistCopy)
}

func main() {

	fmt.Println("Hello from server")

	ShopListExample()

	//orchestrator.OrchestratorExample();
}
