package shoppingList

import (
	"fmt"
	LexCounter "sdle/server/CRDT/lexCounter"
	StringStandardizer "sdle/server/utils/stringStandardizer"
	"github.com/google/uuid"
)


type ShoppingList struct {
	
	url string
	name string
	list LexCounter.LexCounter[string, int]
	state LexCounter.LexCounter[string, int]  // If item state == 0, not bought. If item state >= 1, bought

}

func Create(listName string) ShoppingList {

	list := LexCounter.Create[string, int]("list")
	state := LexCounter.Create[string, int]("state")

	u := uuid.New()

	return ShoppingList{url: u.String(), name: listName, list: list, state: state}
}


func (list ShoppingList) GetURL() string {
	return list.url
}

func (list ShoppingList) GetListName() string {
	return list.name
}

func (list ShoppingList) AddItem(item string, quantity int) bool {
	
	item = StringStandardizer.StandardizeString(item)

	_, keyExists := list.list.Map[item]

	if(keyExists) {
		return false
	}

	itemObj := LexCounter.Create[string, int](item)

	itemObj.Inc(quantity)

	list.list.Join(itemObj)

	return true
}



// Returns false if item doesn't exist or if item already bought. Returns true if alteration was succesful
func (list ShoppingList) BuyItem(item string) bool {
	
	_, keyExists1 := list.list.Map[item]
	_, keyExists2 := list.state.Map[item]
	
	if(!keyExists1 || keyExists2) {
		return false
	}
	
	itemObj := LexCounter.Create[string, int](item)

	itemObj.Inc(1)

	list.state.Join(itemObj)
	
	return true
}

// Returns false if item is already bought or item doesn't exist. Returns true if alteration was succesful
func (list ShoppingList) AlterItemQuantity(item string, newQuantity int) bool {
	
	_, keyExists := list.list.Map[item]
	
	if ((list.state.Map[item].Second >= 1) || !keyExists){
		return false
	}
	
	oldQuantity := list.list.Map[item].Second
	
	if (oldQuantity > newQuantity){
		quantity := oldQuantity - newQuantity

		itemObj := LexCounter.Create[string, int](item)

		itemObj.Dec(quantity)

		list.list.Join(itemObj)

	} else if (oldQuantity < newQuantity) {
		quantity := newQuantity - oldQuantity
		
		itemObj := LexCounter.Create[string, int](item)

		itemObj.Inc(quantity)

		list.list.Join(itemObj)
	}
	
	return true
}


// Returns false if item doesn't exist. Returns true if deletion was succesful
func (list ShoppingList) DeleteItem(item string) bool{
	
	_, keyExists := list.list.Map[item]
	
	if (!keyExists){
		return false
	}
	
	delete(list.list.Map, item)
	delete(list.state.Map, item)
	
	return true
}

// Return false if item not bought or if item doesnt exist. Return true if item bought
func (list ShoppingList) CheckIfItemBought(item string) bool{
	
	val, keyExists := list.state.Map[item]
	
	if(keyExists){
		return (val.Second != 0)	
	}
	return false
}


// Return -1 if item doesn't exist. Returns item quantity if item exists (doesnt matter if bought or not). 
func (list ShoppingList) CheckItemQuantity(item string) int {
	
	entry, keyExists := list.list.Map[item]
	
	if (!keyExists) {
		return -1
	}
	
	return entry.Second
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

	for key, _ := range list.list.Map {
		
		if (comma) { 
			result += ","
		} else {
			comma = true
		}

		bought = list.CheckIfItemBought(key)
		quantity := list.CheckItemQuantity(key)

		result += fmt.Sprintf( "{item: \"%s\", quantity: %d, bought: %t}", key, quantity, bought)

	}

	result += "}"

	return result
}


func (list1 ShoppingList) JoinShoppingListHelper(list2 ShoppingList, item2 string) (int, int) {
	
	decreaseQuantityValue := 0
	decreaseStateValue := 0

	bought1 := list1.CheckIfItemBought(item2)
	bought2 := list2.CheckIfItemBought(item2)

	state1 := list1.state.Map[item2].Second
	state2 := list2.state.Map[item2].Second

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
	
	/*list1.list.Join(list2.list)
	list1.state.Join(list2.state)*/

	quantityMap := make(map[string] int)
	stateMap := make(map[string] int)

	for key, value := range list2.list.Map {
		
		pr, keyExists := list1.list.Map[key]
		if (keyExists && value.First == pr.First){

			decreaseQuantityValue, decreaseStateValue := list1.JoinShoppingListHelper(list2, key)

			quantityMap[key] = decreaseQuantityValue
			stateMap[key] = decreaseStateValue

		}
	}

	list1.list.Join(list2.list)
	list1.state.Join(list2.state)

	for key, value := range quantityMap {

		itemObj := LexCounter.Create[string, int](key)
		itemObj.Dec(value)
	
		pair := list1.list.Map[key]
		pair.Second = itemObj.Map[key].Second
		itemObj.Map[key] = pair
	
		stateObj := LexCounter.Create[string, int](key)
		stateObj.Dec(stateMap[key])
	
		statePair := list1.state.Map[key]
		statePair.Second = stateObj.Map[key].Second
		stateObj.Map[key] = statePair
	
		list1.list.Join(itemObj)
		list1.state.Join(stateObj)

	}

}
