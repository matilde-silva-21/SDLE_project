package orchestrator

import (
	"fmt"
	"net"
	"log"
	"server-utils/orchestrator/hash"
	"server-utils/orchestrator/communication/tcp"
	"server-utils/orchestrator/communication/rabbitMQ"
)

func OrchestratorExample() {

	// <------------ RabbitMQ channel ------------>
	
	conn, ch := rabbitmq.CreateChannel()
	
	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	rabbitmq.DeclareExchange(ch, exchangeName)
	
	q := rabbitmq.DeclareQueue(ch, "")
	
	rabbitmq.BindRoutingKeys(ch, q, exchangeName, "url.*")
	
	messages := rabbitmq.CreateConsumerChannel(ch, q)
	
	go rabbitmq.HandleIncomingMessages(messages) // Go Routine to handle incoming RabbitMQ messages on a separate thread
	
	// <------------------------------------------>


	// <------------ Create Hashing Ring ------------>
	
	hashRing := hash.NewCustomConsistentHash(2, hash.Hash) // MD5 hash

	/*hashRing.Add("server 1")
	hashRing.Add("server 2")
	hashRing.Add("server 3")
	hashRing.Add("server 4")
	hashRing.Add("server 5")*/
	
	fmt.Println(hashRing.GetNodes())
	fmt.Println(hashRing.GetRing())

	//fmt.Println(hashRing.GetClosestNodes("url123", 3))
	
	// <--------------------------------------------->


	// <------------ Create TCP Listener For Servers To join Hash Ring ------------>
	
	address := "localhost:8080" // Orchestrator address
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is listening on port 8080")


	for {
        // Accept incoming connections
        conn, err := listener.AcceptTCP()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Handle client connection.....
		message := tcp.ReadMessage(conn)
		log.Printf(message)
    }
	
	// <--------------------------------------------------------------------------->

}
