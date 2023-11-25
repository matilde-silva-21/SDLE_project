package main

import (
	"fmt"
	"sdle/server/CRDT/addWinSet"
	"sdle/server/CRDT/lexCounter"
	shoppingList "sdle/server/CRDT/shoppingList"
)

func SetExample(){

	CRDT1 := addWinSet.CreateSet()
	CRDT2 := addWinSet.CreateSet()


	fmt.Println(CRDT1)

	addWinSet.Add("apple", &CRDT1)
	addWinSet.Add("pear", &CRDT1)

	addWinSet.Add("cheese", &CRDT2)
	addWinSet.Add("milk", &CRDT2)


	fmt.Println("\nOp1", CRDT1)
	fmt.Println(addWinSet.Contains("apple", CRDT1))
	
	addWinSet.Remove("apple", &CRDT1)
	addWinSet.Remove("basil", &CRDT2)

	
	fmt.Println("\nOp2", CRDT1)
	fmt.Println(addWinSet.Contains("apple", CRDT1))

	set3 := addWinSet.MergeSets(CRDT1, CRDT2)

	fmt.Println("\nOp3", set3)

	addWinSet.Remove("milk", &set3)

	fmt.Println("\nOp4", set3)
	fmt.Println(addWinSet.Contains("milk", set3))

}

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
}
