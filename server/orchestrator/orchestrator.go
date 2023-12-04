package main

import (
	"fmt"
	"net"
	"log"
	"sync"
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"sdle/server/orchestrator/hash"
	"sdle/server/utils/messageStruct"
	"sdle/server/orchestrator/communication/tcp"
	"sdle/server/orchestrator/communication/rabbitMQ"
)

var mutex sync.Mutex

func handleIncomingRabbitMessages(rabbitChannel <-chan amqp.Delivery, hashRing *hash.ConsistentHash, TCPchannels *map[string](chan []byte) ) {

	log.Printf("[RabbitMQ] Waiting for Client logs.")
	for msg := range rabbitChannel {
		
		messageObject, _ := messageStruct.JSONToMessage(msg.Body)
		url := messageObject.ListURL
		
		log.Printf("[RabbitMQ] Received message (url.%s): %s\n", url, msg.Body)

		mutex.Lock()
		// TODO meter uma flag para evitar fazer esta logica, se nao houver servers ligados (por causa do backup)

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


func handleOutgoingRabbitMessages(messagesToSend chan []byte, ch *amqp.Channel, exchangeName string) {

	for {
		select {
			case payload := <-messagesToSend:
				messageObject, err := messageStruct.JSONToMessage(payload)
				if(err == nil){
					url := messageObject.ListURL
	
					rabbitmq.PublishMessage("text/json", string(payload), ch, exchangeName, "url."+url)
					log.Printf("[RabbitMQ] Sent message (%s): %s\n", url, payload)
				}
		}
	}

}


func readTCPConnection(conn *net.TCPConn, hashRing *hash.ConsistentHash, outboundIP string, messagesToSend chan []byte, rabbitChannel chan []byte) {

	for {

		select {

			case payload := <-(messagesToSend): // When channel has a message, route it to the server
				tcp.SendMessage(conn, string(payload))
				log.Printf("[TCP] Sent message to %s: %s\n", outboundIP, string(payload))

			default:
				message, err := tcp.ReadMessage(conn, 1)
		
				if (err != nil) {
					mutex.Lock()
					hashRing.RemoveNodeByIP(outboundIP)
					mutex.Unlock()
					return
				} else if (len(message) != 0) {
					log.Printf("[TCP] Received message from %s: %s\n", outboundIP, message)
					rabbitChannel <- message
				}
		}
		
	}
}



func main() {

	// <------------ Create a map with the channels corresponding to each TCP connection. The channels will share messages between threads ------------>

	channelsMap := make(map[string](chan []byte)) // Key is the IP address, Value is the channel
	outgoingRabbitChannel := make(chan []byte, 100)

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
	
	rabbitmq.BindRoutingKeys(ch, q, exchangeName, "server/url.*")
	
	messages := rabbitmq.CreateConsumerChannel(ch, q)
	
	go handleIncomingRabbitMessages(messages, hashRing, &channelsMap) // Go Routine to handle incoming RabbitMQ messages on a separate thread
	go handleOutgoingRabbitMessages(outgoingRabbitChannel, ch, exchangeName) // Go Routine to handle outgoing RabbitMQ messages on a separate thread
	
	// <-------------------------------------------------------->


	// <------------ Create TCP Listener For Servers To join Hash Ring ------------>
	
	port := "8080" // Default orchestrator port
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	address := "localhost:"+port // Orchestrator address
	
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

    log.Printf("[TCP] Orchestrator is listening on port %s\n\n", port)

	
	// Loop through, waiting for connections from the server
	for {
		conn, err := listener.AcceptTCP()
        if err != nil {
			fmt.Println("Error:", err)
            continue
        }
		
		// Read the very first message that contains the Outbound IP (waits 1 second)
		outboundIP, err := tcp.ReadMessage(conn, 1000)

		if (len(outboundIP) != 0 && err == nil) {
			
			outboundIP := string(outboundIP)

			hashRing.Add(hashRing.GetServerName(), outboundIP) // Add server to hash ring
			
			incomingMessageChannel := make(chan []byte, 100) // Create Channel for the TCP connection know which messages should be sent
			defer close(incomingMessageChannel)
			
			channelsMap[outboundIP] = incomingMessageChannel // Add channel to channel map
			
			log.Printf("Established connection to server with outbound IP: %s", outboundIP)

			go readTCPConnection(conn, hashRing, outboundIP, incomingMessageChannel, outgoingRabbitChannel) // Call thread to continuously poll TCP connection and outgoing messages channel 

		}

    }
	
	// <--------------------------------------------------------------------------->
}
