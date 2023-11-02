package addWinSet

import (
	"github.com/google/uuid"
)

type void struct{}

type SetElem struct {
    element string
    tag  string
}

func CreateSet() map[SetElem]void{

	set := make(map[SetElem]void)

	return set
}

/** Altera o objeto passado como argumento (valor passado por referência). Retorna o pair (Element, UniqueTag)*/
func Add(element string, elements *map[SetElem]void, tombstones *map[SetElem]void) SetElem {

	var dummy void
	
	// Prepare
	u := uuid.New()
	var newVar SetElem = SetElem{element: element, tag: u.String()}

	// Effect
	(*elements)[newVar] = dummy

	for item := range *tombstones{
		delete(*elements, item)
	}

	return newVar
}

func Remove(element string, elements *map[SetElem]void, tombstones *map[SetElem]void) bool{

	var dummy void
	var action bool = false
	var obituaries []SetElem

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

func Contains(element string, elements map[SetElem]void)bool{

	for item := range elements{
		if item.element == element{
			return true
		}
	}

	return false
}
