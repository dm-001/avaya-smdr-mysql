package main

import (
	"fmt"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
)

const (
	HOST = "10.10.5.6"
	PORT = "3000"
	TYPE = "tcp"
)

func main() {

	conn, err := connectToPBX(HOST, PORT)

	// Create a buffer to store any received data
	received := make([]byte, 4096)

	for true {
		// Read data from the connection until an EOF/Null is read
		_, err = conn.Read(received)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(received))
	}

	conn.Close()
}

// Returns a TCP connection interface for the specified PBX ip and port
func connectToPBX(host, port string) (*net.TCPConn, error) {

	// Make sure PBX address can be resolved
	tcpServer, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		fmt.Println("could not resolve PBX address -", err)
		return nil, err
	}

	// Make a basic TCP connection to the PBX SMDR port
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println("could not connect to PBX smdr port -", err)
		return nil, err
	}

	return conn, nil

}
