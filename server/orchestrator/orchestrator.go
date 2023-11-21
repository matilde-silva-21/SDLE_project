package orchestrator

import (
	"fmt"
	"./hash"
)

func OrchestratorExample() {

	consistentHash := hash.NewCustomConsistentHash(2, hash.Hash)

	consistentHash.Add("server 1")
	consistentHash.Add("server 2")
	consistentHash.Add("server 3")
	consistentHash.Add("server 4")
	consistentHash.Add("server 5")


	/*fmt.Println(consistentHash.GetNodes())
	fmt.Println(consistentHash.GetRing())*/

	fmt.Println(consistentHash.GetClosestNodes("url123", 3))


}
