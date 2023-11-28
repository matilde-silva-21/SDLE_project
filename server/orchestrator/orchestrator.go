package orchestrator

import (
	"fmt"
	"net"
	"log"
	"sync"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"sdle/server/orchestrator/hash"
	"sdle/server/utils/messageStruct"
	"sdle/server/orchestrator/communication/tcp"
	"sdle/server/orchestrator/communication/rabbitMQ"
)

var mutex sync.Mutex

func transferIncomingRabbitMessages(rabbitChannel <-chan amqp.Delivery, hashRing *hash.ConsistentHash) {

	log.Printf("[RabbitMQ] Waiting for Client logs.")
	for msg := range rabbitChannel {
		log.Printf("[RabbitMQ] %s", msg.Body)

		messageObject := messageStruct.JSONToMessage(msg.Body)
		url := messageObject.ListURL

		mutex.Lock()

		nodes, _ := hashRing.GetClosestNodes(url, 3)
		var ip string
		ipList := []string{}

		for _, elem := range(nodes){

			if node, ok := elem.(string); ok {
				ip = hashRing.GetServerIP(node)
				ipList = append(ipList, ip)
			} else {
				log.Print("Unexpected type for nodes key.")
			}
		}

		mutex.Unlock()

		// TODO send message with addresses to the first IP
	}
}

func readTCPConnection(conn *net.TCPConn, hashRing *hash.ConsistentHash) {

	for {

		message := tcp.ReadMessage(conn)

		if len(message) == 0  {
			mutex.Lock()

			hashRing.RemoveNodeByIP(conn.RemoteAddr().String())
			fmt.Println(hashRing.GetNodes())
			mutex.Unlock()
			break
		}

		log.Printf("[TCP] %s", message)
		// TODO handle Server Message
	}
}



func OrchestratorExample() {

	// <------------ Go channel for sharing messages between threads ------------>

	incomingMessageChannel := make(chan []byte, 100)
	defer close(incomingMessageChannel)

	// <------------------------------------------------------------------------->
	


	// <------------ Create Hashing Ring ------------>
	
	hashRing := hash.NewCustomConsistentHash(2, hash.Hash) // MD5 hash
	
	// <--------------------------------------------->



	// sudo docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
	// <------------ RabbitMQ communication channel ------------>

	conn, ch := rabbitmq.CreateChannel()
	
	defer conn.Close()
	defer ch.Close()
	
	exchangeName := "logs"
	
	rabbitmq.DeclareExchange(ch, exchangeName)
	
	q := rabbitmq.DeclareQueue(ch, "")
	
	rabbitmq.BindRoutingKeys(ch, q, exchangeName, "url.*")
	
	messages := rabbitmq.CreateConsumerChannel(ch, q)
	
	go transferIncomingRabbitMessages(messages, hashRing) // Go Routine to handle incoming RabbitMQ messages on a separate thread
	
	// <------------------------------------------>


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

    log.Printf("[TCP] Server is listening on port 8080\n\n")

	
	// Loop through, waiting for connections from the server
	for {
		conn, err := listener.AcceptTCP()
        if err != nil {
			fmt.Println("Error:", err)
            continue
        }

		hashRing.Add(hashRing.GetServerName() , conn.RemoteAddr().String())
		go readTCPConnection(conn, hashRing)
    }
	
	// <--------------------------------------------------------------------------->
}
