package main

import (
	"fmt"
	"net"
	"os"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

const (
	HOST = "10.10.5.6"
	PORT = "3000"
	TYPE = "tcp"
)

func main() {

	conn, err := connectToPBX(HOST, PORT)
	if err != nil {
		fmt.Println("could not connect to PBXD -", err)
	}

	db, err := connectToDatabase()
	if err != nil {
		fmt.Println("could not open database connection -", err)
		return
	} else {
		db.Ping()
	}

	// Create a buffer to store any received data
	received := make([]byte, 4096)

	for {
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

func connectToDatabase() (*sql.DB, error) {

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS"), // e.g ."127.0.0.1:3306"
		DBName:               "AvayaCdr",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("could not open database connection -", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("could not ping database -", err)
		return nil, err
	}
	fmt.Println("Connected!")

	return db, nil

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
