package main

import (
	"fmt"
	"sdle/b/CRDT/addWinSet"
	"sdle/b/CRDT/lexCounter"
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

func main() {
	fmt.Println("Hello from server")
	SetExample()
	LexExample()
}
