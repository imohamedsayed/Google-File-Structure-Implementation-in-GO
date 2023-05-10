package main

import (
	"fmt"
	"net"
	"os"
)

var IPs [1]string // 1 IP for 1 slaves

func main() {
	
	// Connect to the master server
	conn, err := net.Dial("tcp", "192.168.43.232:9090") 
	
	if err != nil {
		fmt.Println(err)
		return
	} else { 
		for i := 0; i < 1; i++ { // we have one slave
			buf := make([]byte, 1024) 
			n, err := conn.Read(buf) // reading the IPs for the salves || n the length of the data 
			if err != nil {  
				fmt.Println(err) 
			}
			IPs[i] = string(buf[:n])
			fmt.Println(IPs[i])
		}
	}

	defer conn.Close() // close the connection after receiving slaves' IPs

	fmt.Println("Enter file the name you need: ")
		var file_name string
		fmt.Scan(&file_name)
	var result = ""
	
	for i := 0; i < 4; i++ { // 4 is the number of slaves
		conn, err := net.Dial("tcp", IPs[i]+":9090") 
		if err != nil {
			fmt.Println(err)
			return
		}
		
		
		_,err = conn.Write([]byte(file_name))
		if err != nil{
			fmt.Println(err)
			return
		}
	
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		result += (string(buf[:n]))
	}

	f, err := os.Create("data.txt")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(result)

	if err2 != nil {
		fmt.Println(err2)
	}

}
