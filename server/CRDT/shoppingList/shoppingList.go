package main

import (
	"fmt"
	LexCounter "sdle/server/CRDT/lexCounter"
)


type ShoppingList struct {
	
	list map[string](LexCounter.LexCounter[string, int])
	state map[string](LexCounter.LexCounter[string, int])  // If item state == 0, not bought. If item state >= 1, bought

}

func Create() ShoppingList {

	list := make(map[string] LexCounter.LexCounter[string, int])
	state := make(map[string] LexCounter.LexCounter[string, int])

	return ShoppingList{list: list, state: state}
}

func (list ShoppingList) AddItem(item string, quantity int) bool {
	
	_, keyExists := list.list[item]

	if(keyExists) {
		return false
	}

	itemObject := LexCounter.Create[string, int](item)
	itemState := LexCounter.Create[string, int](item)

	itemObject.Inc(quantity)

	list.list[item] = itemObject
	list.state[item] = itemState

	return true
}


/*
{
	{
		item: apple,
		quantity: 5,
		bought: true 
	},
	{
		item: pear,
		quantity: 1,
		bought: false 
	}
}

*/
func (list ShoppingList) JSON() string{

	var bought, comma = false, false

	result := "{"

	for key, value := range list.list {
		
		if (comma) { 
			result += ","
		} else {
			comma = true
		}

		if (list.state[key].GetValue() >= 1){
			bought = true
		} else {
			bought = false
		}

		result += fmt.Sprintf( "\n{item: %s, quantity: %d, bought: %t}", key, value.GetValue(), bought )

	}

	result += "\n}"

	return result
}

// Returns false if item doesn't exist or if item already bought. Returns true if alteration was succesful
func (list ShoppingList) BuyItem(item string) bool {
	
	_, keyExists := list.list[item]

	if(!keyExists || (list.state[item].GetValue() >= 1)) {
		return false
	}

	list.state[item].Inc(1)

	return true
}

// Returns false if item is already bought or item doesn't exist. Returns true if alteration was succesful
func (list ShoppingList) AlterItemQuantity(item string, newQuantity int) bool {

	_, keyExists := list.list[item]

	if ((list.state[item].GetValue() >= 1) || !keyExists){
		return false
	}
	
	oldQuantity := list.list[item].GetValue()
	
	if(oldQuantity >= newQuantity){
		quantity := oldQuantity - newQuantity
		list.list[item].Dec(quantity)
	} else {
		quantity := newQuantity - oldQuantity
		list.list[item].Inc(quantity)
	}

	return true
}

// Returns false if item doesn't exist. Returns true if deletion was succesful
func (list ShoppingList) DeleteItem(item string) bool{

	_, keyExists := list.list[item]

	if (!keyExists){
		return false
	}

	delete(list.list, item)
	delete(list.state, item)

	return true
}


func (list1 ShoppingList) JoinShoppingList(list2 ShoppingList) {

	decreaseQuantityValue := 0
	decreaseStateValue := 0

	for item2, lexCounter2 := range list2.list {
		
		_, keyExists := list1.list[item2]

		if (keyExists){

			state1Value := list1.state[item2].GetValue()
			state2Value := list2.state[item2].GetValue()

			if (state1Value >= 1 && state2Value == 0) {

				// FIXME perceber como fazer as coisas de decrementar os valores e tal

				if (lexCounter2.GetValue() < list1.list[item2].GetValue()){ 
					decreaseStateValue = 0
					decreaseQuantityValue = lexCounter2.GetValue()

				} else if (lexCounter2.GetValue() > list1.list[item2].GetValue()) {
					decreaseStateValue = state1Value
					decreaseQuantityValue = list1.list[item2].GetValue()
				}


			} else if (state1Value == 0 && state2Value >= 1){


				if (lexCounter2.GetValue() < list1.list[item2].GetValue()){ 
					decreaseStateValue = state2Value
					decreaseQuantityValue = lexCounter2.GetValue()
				} else if (lexCounter2.GetValue() > list1.list[item2].GetValue()) {
					decreaseStateValue = 0
					decreaseQuantityValue = list1.list[item2].GetValue()
				}

			}

			list1.list[item2].Join(lexCounter2)
			list1.state[item2].Join(list2.state[item2])

			fmt.Println(list1.list[item2])
			list1.list[item2].Dec(decreaseQuantityValue)
			list1.state[item2].Dec(decreaseStateValue)
			fmt.Println(list1.list[item2])


		} else {
			list1.list[item2] = lexCounter2
			list1.state[item2] = list2.state[item2]
		}
	}
}


func main() {
	
	shopList1 := Create()
	shopList2 := Create()

	shopList1.AddItem("apple", 3)
	shopList1.AddItem("rice", 5)
	
	shopList2.AddItem("pear", 2)
	shopList2.AddItem("rice", 3)
	shopList2.BuyItem("rice")

	fmt.Println("\n\n")
	fmt.Println(shopList1.JSON())
	fmt.Println("\n\n")

	//fmt.Println(shopList1)

	shopList1.JoinShoppingList(shopList2)

	//fmt.Println(shopList1.JSON())
	fmt.Println(shopList1.JSON())

}


/*

 	if (item2.bought && !item1.bought)
		join as normal but decrease the already bought quantity
	if (!item2.bought && item1.bought)
		if item1.quantity >= item2.quantity 
*/
