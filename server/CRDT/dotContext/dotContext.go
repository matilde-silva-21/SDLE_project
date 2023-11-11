package dotContext

import (
	"reflect"
)


type void struct{}

type Pair[K comparable] struct{
	first K
	second int
}


type DotContext[K comparable] struct {
	cc map[K](int)
	dc map[Pair[K]]void
}


func Create[K comparable]() DotContext[K]{
	finalCC := make(map[K]int)
	finalDC := make(map[Pair[K]]void)

	finalContext := DotContext[K] {cc: finalCC, dc: finalDC}

	return finalContext
}

func Dotin[K comparable] (d Pair[K], context *DotContext[K]) bool{

	cc := (*context).cc
	dc := (*context).dc

	itm, keyExists1 := cc[d.first]
	
	if (keyExists1 && d.second <= itm) {return true}

	_, keyExists2 := dc[d]

	if (keyExists2) {return true}

	return false
}


func Makedot[K comparable] (id K, context *DotContext[K]) Pair[K]{

	cc := (*context).cc

	itm, keyExists1 := cc[id]

	if keyExists1 {
		cc[id] = itm + 1
		return Pair[K]{first: id, second: (itm + 1)}
	} else {
		cc[id] = 1
		return Pair[K]{first: id, second: 1}
	}

}


func Compact[K comparable] (context *DotContext[K]){

	flag := true
	dc := (*context).dc
	cc := (*context).cc

	for (flag) {

		flag = false;

		for pr, _ := range dc{
		
			itm, keyExists := cc[pr.first]

			//FIXME será possivel que apagar um elemento do mapa enquanto o estamos a iterar vai dar merda?
			// If pr not in CC
			if (!keyExists){
				if (pr.second == 1){ // Can compact
					cc[pr.first] = pr.second
					delete(dc, pr)
					flag = true
				}
			} else {

				if (pr.second == (itm + 1)){
					cc[pr.first] = (itm + 1)
					delete(dc, pr)
					flag = true
				} else {

					if (pr.second <= itm) {
						delete(dc, pr)
					} 
				}
			}

		}
	
	}
}


func Insertdot[K comparable] (d Pair[K], context *DotContext[K], compactNow bool) {

	var dummy void

	dc := (*context).dc

	dc[d] = dummy

	if (compactNow) {
		Compact(context)
	}

}


func Join[K comparable] (context1 *DotContext[K], context2 *DotContext[K]) DotContext[K]{

	cc1 := (*context1).cc
	dc1 := (*context1).dc

	cc2 := (*context2).cc
	dc2 := (*context2).dc

	ccEqual := reflect.DeepEqual(cc2, cc1)
	dcEqual := reflect.DeepEqual(dc2, dc1)

	finalContext := Create[K]()

	// Join is idempotent, so just dont do it.
	if(ccEqual && dcEqual) { return *context2 }


	// Loop through cc1, if key not in cc2, add to finalContext.cc. If key in cc2, keep the maximum of the two values, and add it to finalContext.cc
	for first, second := range cc1{
		
		item, keyExists := cc2[first]

		if (!keyExists) {
			finalContext.cc[first] = second
		} else {
			finalContext.cc[first] = max(item, second)
		}

	}

	// Loop through cc2, if key not in cc1, add to finalContext.cc. All other keys that intersect with cc1, have already been added
	for first, second := range cc2{
		
		_, keyExists := cc1[first]

		if (!keyExists) {
			finalContext.cc[first] = second
		}
	}


	// Loop through dc1, add all entries to finalContext
	for pr, _ := range dc1{
		
		Insertdot(pr, &finalContext, false)
	}

	// Loop through dc2, add all entries to finalContext
	for pr, _ := range dc2{

		Insertdot(pr, &finalContext, false)
	}

	Compact(&finalContext)

	return finalContext
}


