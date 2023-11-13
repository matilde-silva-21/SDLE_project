package shoppingList_old

import (
	"fmt"
	LexCounter "sdle/server/CRDT/lexCounter"
	StringStandardizer "sdle/server/utils/stringStandardizer"
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
	
	item = StringStandardizer.StandardizeString(item)

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

		bought = list.CheckIfItemBought(key)

		result += fmt.Sprintf( "\n{item: \"%s\", quantity: %d, bought: %t}", key, value.GetValue(), bought )

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
	
	if (oldQuantity > newQuantity){
		quantity := oldQuantity - newQuantity
		list.list[item].Dec(quantity)
	} else if (oldQuantity < newQuantity) {
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

// Return false if item not bought or if item doesnt exist. Return true if item bought
func (list ShoppingList) CheckIfItemBought(item string) bool{
	
	entry, keyExists := list.state[item]
	
	if (!keyExists) {
		return false
	}

	if (entry.GetValue() >= 1){
		return true
	} else {
		return false;
	}
}

// Returns state value if item exists. Return -1 if item does not exist. 
func (list ShoppingList) GetStateValue(item string) int{
	
	entry, keyExists := list.state[item]
	
	if (!keyExists) {
		return -1
	}

	return entry.GetValue()
}

// Return -1 if item doesn't exist. Returns item quantity if item exists. 
func (list ShoppingList) CheckItemQuantity(item string) int {
	
	entry, keyExists := list.list[item]
	
	if (!keyExists) {
		return -1
	}

	return entry.GetValue()
}


func (list1 ShoppingList) JoinShoppingListHelper(list2 ShoppingList, item2 string) (int, int) {
	
	decreaseQuantityValue := 0
	decreaseStateValue := 0

	bought1 := list1.CheckIfItemBought(item2)
	bought2 := list2.CheckIfItemBought(item2)

	state1 := list1.GetStateValue(item2)
	state2 := list2.GetStateValue(item2)

	quantity1 := list1.CheckItemQuantity(item2)
	quantity2 := list2.CheckItemQuantity(item2)

	if (bought1 && !bought2) {

		if (quantity1 > quantity2) {
			decreaseQuantityValue = quantity2
			decreaseStateValue = 0
		} else if (quantity2 > quantity1) {
			decreaseQuantityValue = 2*quantity1
			decreaseStateValue = state1
		}


	} else if (!bought1 && bought2){

		if (quantity1 > quantity2){ 
			decreaseQuantityValue = 2*quantity2
			decreaseStateValue = state2
		} else if (quantity2 > quantity1) {
			decreaseQuantityValue = quantity1
			decreaseStateValue = 0
		}

	}

	return decreaseQuantityValue, decreaseStateValue

}


func (list1 ShoppingList) JoinShoppingList(list2 ShoppingList) {

	decreaseQuantityValue := 0
	decreaseStateValue := 0

	for item2, lexCounter2 := range list2.list {
		
		_, keyExists := list1.list[item2]

		if (keyExists){

			decreaseQuantityValue, decreaseStateValue = list1.JoinShoppingListHelper(list2, item2)

			list1.list[item2].Join(lexCounter2)
			list1.state[item2].Join(list2.state[item2])

			list1.list[item2].Dec(decreaseQuantityValue)
			list1.state[item2].Dec(decreaseStateValue)

		} else {
			list1.list[item2] = lexCounter2
			list1.state[item2] = list2.state[item2]
		}
	}
}
