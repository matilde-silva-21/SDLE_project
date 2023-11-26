package tcp

import (
	"log"
	"net"
	"time"
	"fmt"
)

// CreateConnection creates a TCP connection to the specified address.
func CreateConnection(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}


// CreateListenerConnection creates a TCP listener on the specified address.
func CreateListenerConnection(address string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// EndConnection closes the TCP connection.
func EndConnection(conn *net.TCPConn) {
	conn.Close()
}


// EndListenerConnection closes the TCP listener.
func EndListenerConnection(listener *net.TCPListener) {
	listener.Close()
}

// SendMessage sends a message over the TCP connection.
func SendMessage(conn *net.TCPConn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Print("Error sending message:", err)
	}
}

// ReadMessage reads a message from the TCP connection.
func ReadMessage(conn *net.TCPConn) []byte {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Print("Error reading message:", err)
		return []byte{}
	}
	return buffer[:n]
}

// RespondPing responds to a ping message with a pong.
func RespondPing(conn *net.TCPConn) {
	for {
		message := ReadMessage(conn)
		log.Printf("Received ping: %s, now sending pong...", message)
		SendMessage(conn, "PONG")
		time.Sleep(1 * time.Second) // Adjust the delay as needed
	}
}


// SendPing sends a ping message and measures the round trip time.
func SendPing(conn *net.TCPConn) {
	for {
		log.Println("Sending ping...")
		startTime := time.Now()
		SendMessage(conn, "PING")
		response := ReadMessage(conn)
		elapsed := time.Since(startTime)
		log.Printf("Received: %s, RTT: %s", response, elapsed)
		time.Sleep(1 * time.Second) // Adjust the delay as needed
	}
}

func SenderExample() {
	address := "localhost:12345"
	conn, err := CreateConnection(address)
	if err != nil {
		log.Fatal("Error creating connection:", err)
	}
	defer EndConnection(conn)

	SendPing(conn)
}

func ReceiverExample() {
	address := "localhost:12345"
	listener, err := CreateListenerConnection(address)
	if err != nil {
		log.Fatal("Error creating listener:", err)
	}
	defer EndListenerConnection(listener)

	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Fatal("Error accepting connection:", err)
	}
	defer conn.Close()

	RespondPing(conn)
}


func ServerSocketExample(){
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is listening on port 8080")

	for {
        // Accept incoming connections
        conn, err := listener.AcceptTCP()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Handle client connection.....
		message := ReadMessage(conn)
		log.Printf("[x] %s", message)
    }
}
