package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	// Connect to the master server
	conn, err := net.Dial("tcp", "192.168.43.232:9090") // Replace with the master's IP address
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Send the slave ID to the master
	slaveID := 2 // Replace with your desired slave ID
	_, err = conn.Write([]byte(strconv.Itoa(slaveID)))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send commands to the master and read responses
	//for {
	var cmd string = "192.168.43.114"
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err)
		return
	}
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Slave " + string(slaveID) + " waiting")
	conn1, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Client Connected")
	}
	buf := make([]byte, 1024)
	n, err := conn1.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	file := (string(buf[:n]))
	// read file path
	filePath := file

	// Read the contents of the file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	// Convert the byte slice to a string
	var msg string = string(contentBytes)

	_, err = conn1.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}

}
