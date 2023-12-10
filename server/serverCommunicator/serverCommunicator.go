package serverCommunicator

import (
	"net"
	"log"
	"fmt"
	"os"
	"time"
	"sdle/server/utils/messageStruct"
	"sdle/server/utils/communication/tcp"
	"sdle/server/utils/CRDT/shoppingList"
	"sdle/server/database"

	_ "github.com/mattn/go-sqlite3"
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


func ConnectToOrchestrator(orchestratorAddress, backupAddress, outboundIP string, sqliteRepository *database.SQLiteRepository) {
	
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

		err = ListenToOrchestrator(conn, orchChannel, sqliteRepository)
	
		if(err != nil){
			log.Printf("Orchestrator connection unexpectedly shut down. Trying to reconnect.\n")
		}
	
	}

}

func ListenToOrchestrator(conn *net.TCPConn, orchChannel chan []byte, sqliteRepository *database.SQLiteRepository) error {
	
	// Sempre que recebo uma mensagem do orchestrator, um Quorum é iniciado.
	for {
		select{

			case message := <- orchChannel:
				conn.Write(message)
				log.Printf("Sent message to orchestrator: %s\n", string(message))

			default:
				buffer, err := tcp.ReadMessage(conn, 1)
				
				if(len(buffer) == 0 && err != nil){

					log.Printf("Connection with endpoint %s encountered a failure: %s\n", conn.RemoteAddr().String(), err)
					return err
				} else if (len(buffer) != 0) {

					IPs, payload := messageStruct.ReadServerMessage(buffer)
					
					go StartQuorumConnection(IPs, payload, orchChannel, sqliteRepository)
				}
				
		}
	}

}

func ExecuteQuorum(conn *net.TCPConn, chanPair ChanPair, payload messageStruct.MessageStruct) error{

	remoteIP := conn.RemoteAddr()

	payload.Action = messageStruct.Read;

	// Send Message asking to read their copy of the ShoppingList of specified URL
	conn.Write(payload.ToJSON())
	log.Printf("Successfully sent read request for URL %s to IP %s.", payload.ListURL, remoteIP)

	buffer := make([]byte, 1024)

	// Wait for the response
	n, err := conn.Read(buffer)

	if err != nil {
		return err // If the operation fails the calling thread will not receive a notification that this thread is ready and thus the quorum will be aborted.
	}

	log.Printf("Received answer to read request from %s: %s", remoteIP, buffer[:n])

	// Send response to channel so the calling thread can handle the CRDT merging
	chanPair.channel <- (buffer[:n])

	// Wait for the signal that channel was read by the calling thread
	<- chanPair.ready

	// Read channel and send the new version of the CRDT to quorum participant
	messageToSend := <- chanPair.channel
	_, err = conn.Write(messageToSend)

	if err != nil {
		return err // If the operation fails the calling thread will not receive a notification that this thread is ready and thus the quorum will be aborted.
	}

	log.Printf("Successfully sent %s command for URL %s to IP %s.", payload.Action, payload.ListURL, remoteIP)
	
	// Indicate to calling thread that quorum is finished and it may be gracefully terminated
	chanPair.ready <- struct{}{}

	conn.Close()

	return nil
}

func ReadAndMergeCRDT(listResponses *(map[int]([]byte)), payload messageStruct.MessageStruct, sqliteRepository *database.SQLiteRepository) (shoppingList.ShoppingList, error){
	
	finalList := shoppingList.MessageStructToCRDT(payload)

	if(payload.Action == messageStruct.Read){ // If the action is Read, no need to merge the list that comes in the payload
		finalList.ResetShoppingList()
	}

	id, _ := database.GetIDByURL(sqliteRepository, payload.ListURL)

	var err error
	var nilList shoppingList.ShoppingList
	
	for _, response := range *listResponses {
		mess, err := messageStruct.JSONToMessage(response)

		if (err != nil){
			log.Print("An error occured while merging the lists.")
			return nilList, err
		} else if (mess.Body != ""){
			listResponse := shoppingList.MessageByteToCRDT(response)
			finalList.JoinShoppingList(listResponse)
		}
	}

	if(id != 0){ // If List exists
		dbList := finalList.ToDatabaseShoppingList(id)
		localList, err := dbList.Read(sqliteRepository)
		
		if(err != nil){
			log.Print("An error occured while reading from memory.")
			return nilList, err
		}
		
		localCRDT := shoppingList.DatabaseShoppingListToCRDT(localList.(*database.ShoppingList))
		finalList.JoinShoppingList(localCRDT)
		
		dbList = finalList.ToDatabaseShoppingList(id)
		err = dbList.Update(sqliteRepository, dbList)
	} else {
		dbList := finalList.ToDatabaseShoppingList(id)
		_, err = dbList.Create(sqliteRepository)
	}

	if(err != nil){
		log.Print("An error occured while writing to memory.")
		return nilList, err
	}

	return finalList, nil
}

func PollQuorumChannels(channelsMap *([]ChanPair), listResponses *(map[int]([]byte)), payload messageStruct.MessageStruct, orchChannel chan []byte, sqliteRepository *database.SQLiteRepository){

	// Set a timeout for the quorum
	timeout := time.After(10 * time.Second)

	// Poll channels
	for {
		for i, ch := range *channelsMap {

			select {
				case data := <-(ch.channel):
					ch.ready <- struct{}{}
					(*listResponses)[i] = data
				case <-timeout: // Break out of the loop when the timeout is reached. Send error message to orchestrator.
					log.Print("Timeout reached. Aborting Quorum.")

					errorMessage := messageStruct.CreateMessage(payload.ListURL, payload.Username, messageStruct.Error, payload.Body)

					(orchChannel) <- errorMessage.ToJSON()
					return
			}
		}
		if(len(*channelsMap) == len(*listResponses)) {
			
			mergedCRDT, err := ReadAndMergeCRDT(listResponses, payload, sqliteRepository)
			
			if(err != nil){

				errorMessage := messageStruct.CreateMessage(payload.ListURL, payload.Username, messageStruct.Error, payload.Body)
				(orchChannel) <- errorMessage.ToJSON()

				return
			}

			log.Printf("Merge completed: %v", mergedCRDT)
			
	
			log.Print("Sending new version to Quorum Participants.")

			var action messageStruct.MessageType
			if(payload.Action == messageStruct.Read){
				action = messageStruct.Write
			} else {
				action = payload.Action
			}

			newCRDTMessage := mergedCRDT.ConvertToMessageFormat(payload.Username, action)

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

						(orchChannel) <- errorMessage.ToJSON()
						return
				}
			}

			if(payload.Action != messageStruct.Delete){
				log.Print("Quorum ended succesfully! Sending new version to Orchestrator.")
				orchChannel <- newCRDTMessage
			} else {
				log.Print("Quorum ended succesfully!")

			}

			return
		}
	}

}


func StartQuorumConnection(IPs []string, payload messageStruct.MessageStruct, orchChannel chan []byte, sqliteRepository *database.SQLiteRepository){

	minNumConn := (len(IPs) - len(IPs)%2) + 1
	activeConn := 0

	channelsMap := [](ChanPair) {}
	listResponses := make(map[int]([]byte))

	connections := [](*net.TCPConn){}

	for _, ip := range IPs {
		tcpAddr, _ := net.ResolveTCPAddr("tcp", ip)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			log.Printf("Failed to connect to server with outbound IP %s.\n", ip)
			continue
		}

		log.Printf("Succesfully connected to server with outbound IP %s.\n", ip)
		defer conn.Close()



		connChannel := make(chan []byte, 3) // Create Channel for the TCP connection know which messages should be sent
		defer close(connChannel)

		channelsMap = append(channelsMap, ChanPair{channel: connChannel, ready: make(chan struct{})})

		connections = append(connections, conn)
		activeConn += 1

		if (minNumConn <= activeConn) { break }

	}

	// If the minimum number of connections is not reached, abort quorum.
	if (minNumConn > activeConn){

		log.Printf("Minimum number of connections (%d) not reached. Aborting Quorum.", minNumConn)

		errorMessage := messageStruct.CreateMessage(payload.ListURL, payload.Username, messageStruct.Error, payload.Body)

		orchChannel <- errorMessage.ToJSON()

		return
	} else {
		for i, conn := range connections {
			go ExecuteQuorum(conn, channelsMap[i], payload)
		}
	}

	PollQuorumChannels(&channelsMap, &listResponses, payload, orchChannel, sqliteRepository)

	return
}


func ParticipateInQuorum(conn *net.TCPConn, sqliteRepository *database.SQLiteRepository) error{

	remoteIP := conn.RemoteAddr().String()

	defer conn.Close()

	log.Printf("Participating in quorum invoked by %s.", remoteIP)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		if err != nil {

			if err.Error() == "EOF" {
				log.Printf("Connection %s closed by remote side.", remoteIP)
				return nil
			} else {
				log.Printf("Connection with endpoint %s encountered a failure: %s\n", remoteIP, err)
				return err
			}

		}

		message, _ := messageStruct.JSONToMessage(buffer[:n])
		
		id, _ := database.GetIDByURL(sqliteRepository, message.ListURL)
		
		CRDT := shoppingList.MessageByteToCRDT(buffer[:n])
		dbList := CRDT.ToDatabaseShoppingList(id)
		
		switch message.Action {
			case messageStruct.Write:
				log.Printf("Received write command from %s: %s", remoteIP, message)

				if(id == 0){ // If List doesn't yet exist, create it
					_, err = dbList.Create(sqliteRepository)
				} else {
					err = dbList.Update(sqliteRepository, dbList)
				}

				if(err != nil){
					log.Print("An error occured while writing to memory.")
					return err
				}

				log.Print("Successfully wrote to memory.")
				return nil

			case messageStruct.Read:
				log.Printf("Received read request from %s: %s", remoteIP, message)
				messageToSend := []byte{}
				
				if(id == 0){ // If List doesn't exist, send empty body
					messageToSend = messageStruct.CreateMessage(message.ListURL, message.Username, messageStruct.Read, "").ToJSON()
				} else {
					localList, err := dbList.Read(sqliteRepository)
					
					if(err != nil){
						log.Print("An error occured while reading from memory.")
						return err
					}

					localCRDT := shoppingList.DatabaseShoppingListToCRDT(localList.(*database.ShoppingList))
					messageToSend = localCRDT.ConvertToMessageFormat(message.Username, messageStruct.Read)
				}

				conn.Write([]byte(messageToSend))

				log.Print("Successfully sent payload.")
			case messageStruct.Delete:
				log.Printf("Received delete command from %s: %s", remoteIP, message)
				
				if(id != 0) {
					err := dbList.Delete(sqliteRepository)

					if(err != nil){
						log.Print("An error occured while deleting from memory.")
						return err
					}
				}
				
				log.Printf("Successfully deleted from memory.")
				return nil
		}
	}
}


func StartServerCommunication(orchestratorAddress, backupAddress string, sqliteRepository *database.SQLiteRepository) {

	
	outboundIP := GetOutboundIP()

	// <------------ Connect To orchestrator ------------>
	
	go ConnectToOrchestrator(orchestratorAddress, backupAddress, outboundIP, sqliteRepository)

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
		
		go ParticipateInQuorum(conn, sqliteRepository)
	}

}