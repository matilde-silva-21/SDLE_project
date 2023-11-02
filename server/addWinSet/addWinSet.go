package addWinSet

import (
	"github.com/google/uuid"
)

type void struct{}

type setElem struct {
    element string
    tag  string
}

func CreateSet() map[setElem]void{

	set := make(map[setElem]void)

	return set
}

/** Altera o objeto passado como argumento (valor passado por referência). Retorna o pair (Element, UniqueTag)*/
func Add(element string, elements *map[setElem]void, tombstones *map[setElem]void) setElem {

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

func Remove(element string, elements *map[setElem]void, tombstones *map[setElem]void) bool{

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

func Contains(element string, elements map[setElem]void)bool{

	for item := range elements{
		if item.element == element{
			return true
		}
	}

	return false
}


