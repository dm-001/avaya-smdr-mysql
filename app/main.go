package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := connectToDatabase()
	if err != nil {
		fmt.Println("could not open database connection -", err)
		return
	}

	// Get port number from env vars
	targetPort := os.Getenv("LISTEN_PORT")

	listener, err := listenOnPort(":" + targetPort)
	if err != nil {
		fmt.Println("could not bind to port:", err)
		return
	}

	for {
		// Wait for and accept a connection on the port listener
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			// Handle the connection in a goroutine so more can be accepted
			go handleConnection(conn, db)
		}

	}

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

	// Determine port to listen on

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

func listenOnPort(port string) (net.Listener, error) {

	listener, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println("error establishing listener on port: ", err)
		return nil, err
	}
	fmt.Println("Listening for SMDR data on: ", port)
	return listener, nil

}

// Write  big old slice of cdr values to the database, filtering for only the data we care about
func writeToDatabase(callData []string, sqldb *sql.DB) {

	execString := "INSERT INTO cdr (CallStart,ConnectedTime,RingTime,Caller,Direction,CalledNumber,DialledNumber,IsInternal,CallID,HoldTime,ParkTime,ExternalTargetingCause,ExternalTargeterId,ExternalTargetedNumber,CallerServerIP,UniqueCallIDCallerExtension,UniqueCallIDCalledParty,SMDRRecordingTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := sqldb.Exec(execString, callData[0], callData[1], callData[2], callData[3], callData[4], callData[5], callData[6], callData[8], callData[9], callData[15], callData[16], callData[27], callData[28], callData[29], callData[30], callData[31], callData[33], callData[34])
	if err != nil {
		fmt.Println(err)
	}

}

func handleConnection(conn net.Conn, db *sql.DB) error {

	// Create a buffer to store any received data
	received := make([]byte, 4096)

	// Read data from the connection until an EOF/Null is read
	_, err := conn.Read(received)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Split smdr data in to individual elements, trim whitespace and newlines
	smdr := strings.Split(string(received), ",")
	for i := range smdr {
		smdr[i] = strings.TrimSpace(smdr[i])
		smdr[i] = strings.Replace(smdr[i], "\r\n", "", -1)
	}

	// Debug
	//fmt.Println("RAW:  ", string(received))

	// Write call data to the SQL database
	go writeToDatabase(smdr, db)
	conn.Close()

	return nil

}
