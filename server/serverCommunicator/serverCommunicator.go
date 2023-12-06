package serverCommunicator

import (
	"net"
	"log"
	"fmt"
	"os"
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
		log.Print(IPs, payload)
		
		go StartQuorumConnection(IPs)
	}

}

func StartQuorumConnection(IPs []string){

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
		connections = append(connections, conn)
		activeConn += 1

		go listenToConnection(conn) // Posso mandar a mensagem de read logo, para apressar as coisas ? depois de ter o numero certo de mensagens enviadas (e sabendo que essas ligaçoes ainda estao em pé), masndar a mensagem de escrita, maybe?

		if (minNumConn <= activeConn) { break }

	}
	fmt.Println(connections)

    // TODO abortar o quorum if connection not successful
    //fmt.Println(connections)

    // Keep function running until quorum is over or else the connections break down (calling a thread2 inside a thread1 -> thread2 terminates when thread1 terminates)
    for {
        
    }

}

func listenToConnection(conn *net.TCPConn) error {

	for {

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {

			log.Printf("Connection with endpoint %s encountered a failure: %s\n", conn.RemoteAddr().String(), err)
			return err
		}

		log.Print(string(buffer[:n]))
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