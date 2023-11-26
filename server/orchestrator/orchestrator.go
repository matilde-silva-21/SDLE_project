package orchestrator

import (
	"fmt"
	"net"
	"log"
	"strconv"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"server-utils/orchestrator/hash"
	//"server-utils/messageStruct"
	"server-utils/orchestrator/communication/tcp"
	"server-utils/orchestrator/communication/rabbitMQ"
)


func transferIncomingRabbitMessages(rabbitChannel <-chan amqp.Delivery) {

	log.Printf("[RabbitMQ] Waiting for Client logs.")
	for msg := range rabbitChannel {
	   log.Printf("[x] %s", msg.Body)
	   // TODO handle Client Message
	}
}

func readTCPConnection(conn *net.TCPConn) {

	for {

		message := tcp.ReadMessage(conn)

		if len(message) == 0  {
			// TODO Connection closed or error occurred
			break
		}

		log.Printf("[o] %s", message)
	   // TODO handle Server Message
	}
}



func OrchestratorExample() {

	// <------------ Go channel for sharing messages between threads ------------>

	incomingMessageChannel := make(chan []byte, 100)
	defer close(incomingMessageChannel)

	// <------------------------------------------------------------------------->
	

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
	
	go transferIncomingRabbitMessages(messages) // Go Routine to handle incoming RabbitMQ messages on a separate thread
	
	// <------------------------------------------>


	// <------------ Create Hashing Ring ------------>
	
	hashRing := hash.NewCustomConsistentHash(2, hash.Hash) // MD5 hash
	
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

    log.Printf("[TCP] Server is listening on port 8080\n\n")

	
	// Loop through, waiting for connections from the server
	for {
		conn, err := listener.AcceptTCP()
        if err != nil {
			fmt.Println("Error:", err)
            continue
        }
		
		hashRing.Add("server " + strconv.Itoa(hashRing.GetNumberOfKeys()) , conn)
		go readTCPConnection(conn)
    }
	
	// <--------------------------------------------------------------------------->
}
