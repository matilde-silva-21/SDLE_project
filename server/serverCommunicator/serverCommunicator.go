package serverCommunicator

import (
	"net"
	"log"
	"fmt"
	"os"
	"time"
	"sdle/server/utils/messageStruct"
	//"sdle/server/utils/CRDT/shoppingList"
)

type ChanPair struct {
	channel  chan []byte
	ready chan struct{}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().String()

    return localAddr
}


func connectToOrchestrator(outboundIP string) {

	orchestratorAddress := "localhost:8080"
	backupAddress := "localhost:8081"
	
	orchTcpAddr, err := net.ResolveTCPAddr("tcp", orchestratorAddress)

	if err != nil {
		os.Exit(0)
	}

	backupTcpAddr, err := net.ResolveTCPAddr("tcp", backupAddress)
	
	if err != nil {
		os.Exit(0)
	}
	
	orchChannel := make(chan []byte, 10) // There probably won't be more than 10 quorums at the same time.
	
	for {
		// Try to connect to Orchestrator, in case of failure, connect to backup orchestrator. If both connections fail, the program stops
		connOrchestrator, err1 := net.DialTCP("tcp", nil, orchTcpAddr)
		connBackup, err2 := net.DialTCP("tcp", nil, backupTcpAddr)

		conn := connOrchestrator

		if err1 != nil {
			log.Printf("Failed to connect to orchestrator. Connecting to backup.\n")

			if err2 != nil {
				log.Printf("Failed to connect to backup. Re-trying\n")
				continue
			}
			conn = connBackup
		}
	
		defer conn.Close()
		
		conn.Write([]byte(outboundIP))

		err = listenToOrchestrator(conn, &orchChannel)
	
		if(err != nil){
			log.Printf("Orchestrator connection unexpectedly shut down. Trying to reconnect.\n")
		}
	
	}

}

func listenToOrchestrator(conn *net.TCPConn, orchChannel *chan []byte) error {

	buffer := make([]byte, 1024)
	
	// Sempre que recebo uma mensagem do orchestrator, um Quorum é iniciado.
	for {
		select{

			case message := <- (*orchChannel):
				conn.Write(message)
				log.Printf("Sent message to orchestrator: %s\n", string(message))

			default:
				n, err := conn.Read(buffer)
				
				if err != nil {
					
					log.Printf("Connection with endpoint %s encountered a failure: %s\n", conn.RemoteAddr().String(), err)
					return err
				}
				
				IPs, payload := messageStruct.ReadServerMessage(buffer[:n])
				
				go StartQuorumConnection(IPs, payload, orchChannel)
		}
	}

}

func ExecuteQuorum(conn *net.TCPConn, chanPair ChanPair, payload messageStruct.MessageStruct) error{

	remoteIP := conn.RemoteAddr()

	// TODO criar mensagem de pedido leitura 
	payload.Action = messageStruct.Read;

	// Send Message asking to read their copy of the ShoppingList of specified URL
	conn.Write(payload.ToJSON()) // TODO change Placeholder
	log.Printf("Successfully sent read request for URL %s to IP %s.", payload.ListURL, remoteIP)

	buffer := make([]byte, 1024)

	// Wait for the response
	n, err := conn.Read(buffer)

	if err != nil {
		return err // If the operation fails the calling thread will not receive a notification that this thread is ready and thus the quorum will be aborted.
	}

	log.Printf("Received answer to read request from %s: %s", remoteIP, buffer)

	// Send response to channel so the calling thread can handle the CRDT merging
	chanPair.channel <- (buffer[:n])

	// Wait for the signal that channel was read by the calling thread
	<- chanPair.ready

	// Read channel and send the new version of the CRDT to quorum participant
	messageToSend := <- chanPair.channel // TODO nao esquecer de meter Action.Write na outgoing message.
	_, err = conn.Write(messageToSend)

	if err != nil {
		return err // If the operation fails the calling thread will not receive a notification that this thread is ready and thus the quorum will be aborted.
	}

	log.Print("Successfully sent Write command for URL %s to IP %s.", payload.ListURL, remoteIP)
	
	// Indicate to calling thread that quorum is finished and it may be gracefully terminated
	chanPair.ready <- struct{}{}

	conn.Close()

	return nil
}


func PollQuorumChannels(channelsMap *map[string](ChanPair), listResponses *map[string]([]byte), payload messageStruct.MessageStruct, orchChannel *chan []byte){

	// Set a timeout for the quorum
	timeout := time.After(10 * time.Second)

	// Poll channels
	for {
		for key, ch := range *channelsMap {

			select {
				case data := <-(ch.channel):
					ch.ready <- struct{}{}
					(*listResponses)[key] = data
				
				case <-timeout: // Break out of the loop when the timeout is reached. Send error message to orchestrator.
					log.Print("Timeout reached. Aborting Quorum.")

					errorMessage := messageStruct.CreateMessage(payload.ListURL, payload.Username, messageStruct.Error, payload.Body)

					(*orchChannel) <- errorMessage.ToJSON()
					return
			}
		}
		if(len(*channelsMap) == len(*listResponses)) {
			// TODO read and merge all CRDTs and send back responses
			
			log.Print("Merging all CRDTs and sending new version to Quorum Participants.")

			newCRDTMessage := payload.ToJSON() // placeholder
			for _, ch := range *channelsMap {
				ch.channel <- newCRDTMessage
			}

			// End the program only when all threads have already been read
			for _, ch := range *channelsMap {

				select {
					case <- ch.ready:
						continue
					
					case <-timeout: // Break out of the loop when the timeout is reached. Send error message to orchestrator.
						log.Print("Timeout reached. Aborting Quorum.")

						errorMessage := messageStruct.CreateMessage(payload.ListURL, payload.Username, messageStruct.Error, payload.Body)

						(*orchChannel) <- errorMessage.ToJSON()
						return
				}
			}

			log.Printf("Quorum ended succesfully! Sending new version to Orchestrator.")
			// TODO mandar o novo CRDT ao Orch
			return
		}
	}

}


func StartQuorumConnection(IPs []string, payload messageStruct.MessageStruct, orchChannel *chan []byte){

	channelsMap := make(map[string](ChanPair))
	listResponses := make(map[string]([]byte))

	minNumConn := (len(IPs) - len(IPs)%2) + 1
	activeConn := 0

	connections := [](*net.TCPConn){}

	for _, ip := range IPs {
		tcpAddr, err := net.ResolveTCPAddr("tcp", ip)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			log.Printf("Failed to connect to server with outbound IP %s.\n", ip)
			continue
		}

		log.Printf("Succesfully connected to server with outbound IP %s.\n", ip)
		defer conn.Close()



		connChannel := make(chan []byte, 3) // Create Channel for the TCP connection know which messages should be sent
		defer close(connChannel)

		channelsMap[ip] = ChanPair{channel: connChannel, ready: make(chan struct{})}

		connections = append(connections, conn)
		activeConn += 1




		go ExecuteQuorum(conn, channelsMap[ip], payload)

		if (minNumConn <= activeConn) { break }

	}

	// If the minimum number of connections is not reached, abort quorum.
	if (minNumConn > activeConn){

		log.Printf("Minimum number of connections (%d) not reached. Aborting Quorum.", minNumConn)

		errorMessage := messageStruct.CreateMessage(payload.ListURL, payload.Username, messageStruct.Error, payload.Body)

		*orchChannel <- errorMessage.ToJSON()

		return
	}

	PollQuorumChannels(&channelsMap, &listResponses, payload, orchChannel)

	return
}

func listenToConnection(conn *net.TCPConn) error {

	didIt := false
	for {

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {

			if err.Error() == "EOF" {
				log.Printf("Connection %s closed by remote side.", conn.RemoteAddr().String())
				return nil
			} else {
				log.Printf("Connection with endpoint %s encountered a failure: %s\n", conn.RemoteAddr().String(), err)
				return err
			}

		}

		log.Print(string(buffer[:n]))
		if(!didIt){
			// TODO fazer o codigo do recetor do quorum
			log.Printf("Sending the payload back to see if communication is doing what it should. Payload: %s\n", buffer)
			conn.Write(buffer)
			didIt = true
		}
	}

}


func StartServerCommunication() {

	
	outboundIP := GetOutboundIP()

	// <------------ Connect To orchestrator ------------>
	
	go connectToOrchestrator(outboundIP)

	// <------------------------------------------------->
	


	// <------------ Create listener for other servers to connect ------------>
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", outboundIP)
	if err != nil {
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()

    log.Printf("[TCP] Server is listening on address %s\n\n", outboundIP)

	for {
		conn, err := listener.AcceptTCP()
        if err != nil {
			fmt.Println("Error:", err)
            continue
        }
		fmt.Println("\nTHERES SOMEONE AT THE DOOR\n")
		
		defer conn.Close()
		go listenToConnection(conn)
	}

}