package main

import (
	"fmt"
	"sdle/b/addWinSet"
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


func main() {
	fmt.Println("Hello from server")
	SetExample()
}
