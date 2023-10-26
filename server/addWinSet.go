package addWinSet

import (
	"fmt"
	"github.com/google/uuid"
)

type void struct{}

type setElem struct {
    element string
    tag  string
}

func createSet() map[setElem]void{

	set := make(map[setElem]void)

	return set
}

/** Altera o objeto passado como argumento (valor passado por referência). Retorna o pair (Element, UniqueTag)*/
func add(element string, elements *map[setElem]void, tombstones *map[setElem]void) setElem {

	var dummy void
	
	// Prepare
	u := uuid.New()
	var newVar setElem = setElem{element: element, tag: u.String()}

	// Effect
	(*elements)[newVar] = dummy

	for item := range *tombstones{
		delete(*elements, item)
	}

	return newVar
}

func remove(element string, elements *map[setElem]void, tombstones *map[setElem]void) bool{

	var dummy void
	var action bool = false
	var obituaries []setElem

	// Prepare
	for item := range *elements{
		if item.element == element{
			obituaries = append(obituaries, item)
			action = true
		}
	}

	// Effect
	for _, corpse := range obituaries{
		delete(*elements, corpse)
		(*tombstones)[corpse] = dummy
	}

	return action
}

func contains(element string, elements map[setElem]void)bool{

	for item := range elements{
		if item.element == element{
			return true
		}
	}

	return false
}


func main() {
	
	mySet := createSet()
	myTombstone := createSet()

	fmt.Println(mySet)

	add("apple", &mySet, &myTombstone)
	add("pear", &mySet, &myTombstone)


	fmt.Println("\nOp1", mySet, myTombstone)
	fmt.Println(contains("apple", mySet))
	
	remove("apple", &mySet, &myTombstone)
	
	fmt.Println("\nOp2", mySet, myTombstone)
	fmt.Println(contains("apple", mySet))

}
