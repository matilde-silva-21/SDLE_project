package communicator

import (
	"net"
	"log"
	"fmt"
	//"sdle/server/utils/messageStruct"
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

func connectToOrchestrator() (*net.TCPConn, error){

	orchestratorAddress := "localhost:8080"
	backupAddress := "localhost:8081"
	
	orchTcpAddr, err := net.ResolveTCPAddr("tcp", orchestratorAddress)

	if err != nil {
		return nil, err
	}

	backupTcpAddr, err := net.ResolveTCPAddr("tcp", backupAddress)
	
	if err != nil {
		return nil, err
	}
	

	// Try to connect to Orchestrator, in case of failure, connect to backup orchestrator. If both connections fail, the program stops
	conn, err := net.DialTCP("tcp", nil, orchTcpAddr)
	if err != nil {
		log.Printf("Failed to connect to orchestrator. Trying to talk to backup.\n")
		conn, err = net.DialTCP("tcp", nil, backupTcpAddr)
		if err != nil {
			log.Printf("Failed to connect to backup. Shutting off...\n")
			return nil, err
		}
	}

	return conn, nil

}

func listenToConnection(conn *net.TCPConn) {

	for {

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {

			log.Print("Error reading message: ", err)
			return
		}

		log.Print(string(buffer[:n]))
	}

}


func StartServerCommunication() {

	
	outboundIP := GetOutboundIP()
	fmt.Println(outboundIP)

	// <------------ Connect To orchestrator ------------>
	
	orchestrator, err := connectToOrchestrator()
	
	if err != nil {
		return
	}
	
	defer orchestrator.Close()
	
	orchestrator.Write([]byte(outboundIP))
	
	go listenToConnection(orchestrator)

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

		go listenToConnection(conn)
	}


}