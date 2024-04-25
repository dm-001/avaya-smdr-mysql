package main

import (
	"fmt"
	"net"
	"os"
	"strings"

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
	}

	for {
		// Create a buffer to store any received data
		received := make([]byte, 4096)

		// Read data from the connection until an EOF/Null is read
		_, err = conn.Read(received)
		if err != nil {
			fmt.Println(err)
			break
		}
		// Split smdr data in to individual elements, trim whitespace
		smdr := strings.Split(string(received), ",")
		for i := range smdr {
			smdr[i] = strings.TrimSpace(smdr[i])
		}

		// Write call data to the SQL database
		go writeToDatabase(smdr, db)

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

// Write  big old slice of cdr values to the database, filtering for only the data we care about
func writeToDatabase(callData []string, sqldb *sql.DB) {

	execString := "INSERT INTO cdr (CallStart,ConnectedTime,RingTime,Caller,Direction,CalledNumber,DialledNumber,IsInternal,CallID,HoldTime,ParkTime,ExternalTargetingCause,ExternalTargeterId,ExternalTargetedNumber,CallerServerIP,UniqueCallIDCallerExtension,UniqueCallIDCalledParty,SMDRRecordingTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := sqldb.Exec(execString, callData[0], callData[1], callData[2], callData[3], callData[4], callData[5], callData[6], callData[8], callData[9], callData[15], callData[16], callData[27], callData[28], callData[29], callData[30], callData[31], callData[33], callData[34])
	if err != nil {
		fmt.Println(err)
	}

}
