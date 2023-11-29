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

func transferIncomingRabbitMessages(rabbitChannel <-chan amqp.Delivery, hashRing *hash.ConsistentHash, TCPchannels *map[string](chan []byte) ) {

	log.Printf("[RabbitMQ] Waiting for Client logs.")
	for msg := range rabbitChannel {
		
		messageObject := messageStruct.JSONToMessage(msg.Body)
		url := messageObject.ListURL
		
		log.Printf("[RabbitMQ] Received message (url.%s): %s\n", url, msg.Body)

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
		
		(*TCPchannels)[ipList[0]]<-((messageObject).BuildMessageForServer(ipList)) // Send message body to TCP
		
		mutex.Unlock()
	}
}

func readTCPConnection(conn *net.TCPConn, hashRing *hash.ConsistentHash, messagesToSend chan []byte) {

	ip := conn.RemoteAddr().String()

	for {

		select {

			case payload := <-(messagesToSend): // When channel has a message, route it to the server
				tcp.SendMessage(conn, string(payload))
				log.Printf("[TCP] Sent message to %s: %s\n", ip, string(payload))

			default:
				message, err := tcp.ReadMessage(conn)
		
				if (err != nil) {
					mutex.Lock()
					hashRing.RemoveNodeByIP(ip)
					mutex.Unlock()
					return
				} else if (len(message) != 0) {
					log.Printf("[TCP] Received message from %s: %s\n", ip, message)
				}
		

				// TODO reencaminhar mensagem para os clientes
		}
		
	}
}



func OrchestratorExample() {

	// <------------ Create a map with the channels corresponding to each TCP connection. The channels will share messages between threads ------------>

	channelsMap := make(map[string](chan []byte)) // Key is the IP address, Value is the channel

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
	
	go transferIncomingRabbitMessages(messages, hashRing, &channelsMap) // Go Routine to handle incoming RabbitMQ messages on a separate thread
	
	// <-------------------------------------------------------->


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

		ipString := conn.RemoteAddr().String()
		
		hashRing.Add(hashRing.GetServerName() , ipString) // Add server to hash ring
		
		
		incomingMessageChannel := make(chan []byte, 100) // Create Channel for the TCP connection know which messages should be sent
		defer close(incomingMessageChannel)
		
		channelsMap[ipString] = incomingMessageChannel // Add channel to channel map

		go readTCPConnection(conn, hashRing, incomingMessageChannel) // Call thread to continuously poll TCP connection / outgoing messages channel 
    }
	
	// <--------------------------------------------------------------------------->
}
