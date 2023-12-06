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

		err = listenToOrchestrator(conn)
	
		if(err != nil){
			log.Printf("Orchestrator connection unexpectedly shut down. Trying to reconnect.\n")
		}
	
	}

}

func listenToOrchestrator(conn *net.TCPConn) error {

	for {

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {

			log.Printf("Connection with endpoint %s encountered a failure: %s\n", conn.RemoteAddr().String(), err)
			return err
		}

		IPs, payload := messageStruct.ReadServerMessage(buffer[:n])
		//log.Print(IPs, payload)
		
		go StartQuorumConnection(IPs, payload)
	}

}

func ExecuteQuorum(conn *net.TCPConn, connChannel chan []byte, payload messageStruct.MessageStruct) error{

	payload.Action = messageStruct.Read;
	fmt.Println("\n\n", payload, "\n")

	// Mandar o CRDT mais fresh para lá
	conn.Write(payload.ToJSON())


	buffer := make([]byte, 1024)

	// Ler o CRDT de lá
	n, err := conn.Read(buffer)

	fmt.Println("i read something and i am happy")

	if err != nil { return err } // TODO se houver erro A QUALQUER MOMENTO DA CONEXÃO

	// Meter o CRDT de lá no chan
	connChannel <- (buffer[:n])

	// Add a short delay so the other thread can read the channel and send the message
	time.Sleep(100 * time.Millisecond)

	// Mandar o novo CRDT de volta para lá
	messageToSend := <- connChannel

	conn.Write(messageToSend)

	return nil
}


func StartQuorumConnection(IPs []string, payload messageStruct.MessageStruct){

	channelsMap := make(map[string](chan []byte))
	listResponses := make(map[string]([]byte))

	minNumConn := (len(IPs) - len(IPs)%2) + 1
	activeConn := 0

	connections := [](*net.TCPConn){}

	for _, ip := range IPs {
		tcpAddr, err := net.ResolveTCPAddr("tcp", ip)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			log.Printf("Failed to connect to server with IP %s.\n", ip)
			continue
		}

		fmt.Println("\nI connected bro i swear\n")
		defer conn.Close()



		connChannel := make(chan []byte, 3) // Create Channel for the TCP connection know which messages should be sent
		defer close(connChannel)

		channelsMap[ip] = connChannel

		connections = append(connections, conn)
		activeConn += 1




		go ExecuteQuorum(conn, connChannel, payload)

		if (minNumConn <= activeConn) { break }

	}

    // TODO abortar o quorum if connection not successful


	// Set a timeout for the quorum
	timeout := time.After(10 * time.Second)

	// Poll channels
	for {
		fmt.Println("loop 1", channelsMap)
		for key, ch := range channelsMap {
			fmt.Println("loop 2")

			select {
				case data := <-ch:
					log.Printf("My friend %s answered me!", key)
					listResponses[key] = data
				
				case <-timeout:
					// Break out of the loop when the timeout is reached
					log.Print("Timeout reached. Aborting Quorum.")
					// TODO mandar mensagem de erro de volta ao Orch
					return
			}
		}
		if(len(channelsMap) == len(listResponses)) {
			// TODO read and merge all CRDTs and send back responses (nao esquecer de mandar o novo CRDT ao Orch)
			
			fmt.Println("\n\neveryone was here :)\nnow bye-bye!")

			newCRDTMessage := payload.ToJSON() // placeholder
			for _, ch := range channelsMap {
				ch <- newCRDTMessage
			}
			// Add a short delay so the other threads can read the channel and send the message
			time.Sleep(100 * time.Millisecond)
			return
		}
	}

}

func listenToConnection(conn *net.TCPConn) error {

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

		log.Printf("Sending the payload back to see if communication is doing what it should. Payload: %s\n", buffer)
		conn.Write(buffer)
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