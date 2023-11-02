package main

import (
	"fmt"
	"sdle/b/addWinSet"
)

func SetExample(){
	mySet := addWinSet.CreateSet()
	myTombstone := addWinSet.CreateSet()

	fmt.Println(mySet)

	addWinSet.Add("apple", &mySet, &myTombstone)
	addWinSet.Add("pear", &mySet, &myTombstone)


	fmt.Println("\nOp1", mySet, myTombstone)
	fmt.Println(addWinSet.Contains("apple", mySet))
	
	addWinSet.Remove("apple", &mySet, &myTombstone)
	
	fmt.Println("\nOp2", mySet, myTombstone)
	fmt.Println(addWinSet.Contains("apple", mySet))
}


func main() {
	fmt.Println("Hello from server")
	SetExample()
}
