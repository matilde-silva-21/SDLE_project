/*
*	Implementation of basic (not causal) Add Win Set. Uses UUID for the unique tags
*/

package addWinSet

import (
	"github.com/google/uuid"
)

type void struct{}

type setElem struct {
    element string
    tag  string
}

type addWinSet struct {
	elements map[setElem]void
	tombstones map[setElem]void
}

func CreateSet() addWinSet{

	elements := make(map[setElem]void)
	tombstones := make(map[setElem]void)

	set := addWinSet{elements: elements, tombstones: tombstones}
	
	return set
}

/** Altera o objeto passado como argumento (valor passado por referência). Retorna o pair (Element, UniqueTag)*/
func Add(element string, set *addWinSet) setElem {

	elements := (*set).elements

	var dummy void
	
	// Prepare
	u := uuid.New()
	var newVar setElem = setElem{element: element, tag: u.String()}

	// Effect
	elements[newVar] = dummy

	return newVar
}

func Remove(element string, set *addWinSet) bool{

	var dummy void
	var action bool = false
	var obituaries []setElem

	elements := (*set).elements
	tombstones := (*set).tombstones


	// Prepare
	for item := range elements{
		if item.element == element{
			obituaries = append(obituaries, item)
			action = true
		}
	}

	// Effect
	for _, corpse := range obituaries{
		delete(elements, corpse)
		tombstones[corpse] = dummy
	}

	return action
}

func Contains(element string, set addWinSet) bool{

	elements := set.elements

	for item := range elements{
		if item.element == element{
			return true
		}
	}

	return false
}


func MergeSets(set1 addWinSet, set2 addWinSet) addWinSet{

	var dummy void

	elements1 := set1.elements
	tombstones1 := set1.tombstones

	elements2 := set2.elements
	tombstones2 := set2.tombstones

	newElements := make(map[setElem]void)
	newTombstones := make(map[setElem]void)

	for item := range elements1{
		newElements[item] = dummy
	}

	for item := range tombstones1{
		newTombstones[item] = dummy
	}

	for item := range elements2{
		newElements[item] = dummy
	}

	for item := range tombstones2{
		newTombstones[item] = dummy
	}

	
	newSet := addWinSet{elements: newElements, tombstones: newTombstones}

	return newSet

}


