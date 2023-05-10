package main

import (
	"fmt"
	"net"
	"strconv"
)
const size int = 4
var slaves = make(map[int]net.Conn)
var IPs[size] string
var count int=0
var result = make(chan string)

func main() {
	// Start the master server
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Master server started")

	// Initialize a map to store connected slaves

	conn1, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("Client Connected")
		// Wait for slave connections
	for {

		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		slaveID, err := strconv.Atoi(string(buf[:n]))
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Read the slave ID from the connection


		// Store the connection in the map
		slaves[slaveID] = conn

		fmt.Printf("Slave %d connected\n", slaveID)

		// Start a new goroutine to handle commands from the slave
		count++
		go handleSlave(slaveID, conn)
		IPs[slaveID-1] = <- result
		if count == size{
			break
		}
	}
}
for i := 0 ; i < size; i++{
	fmt.Println(IPs[i])
	_, err = conn1.Write([]byte(IPs[i]))
	if err != nil {
		fmt.Println(err)
		return
	}
}

}

func handleSlave(slaveID int, conn net.Conn) {
	defer conn.Close()


		// Read commands from the slave
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)

		}
		cmd := string(buf[:n])
		result <- cmd
		//fmt.Println(cmd)

}
