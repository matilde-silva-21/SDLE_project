package communicator

import (
	"net"
	"log"
	"fmt"
	"os"
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
				log.Printf("Failed to connect to backup. Shutting off...\n")
				os.Exit(0)
			}
			conn = connBackup
		}
	
		defer conn.Close()
		
		conn.Write([]byte(outboundIP))

		err = listenToConnection(conn)
	
		if(err != nil){
			log.Printf("Orchestrator connection unexpectedly shut down. Trying to reconnect.\n")
		}
	
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

		go listenToConnection(conn)
	}

	// TODO quando orchestrator morrer, conectar ao backup

}