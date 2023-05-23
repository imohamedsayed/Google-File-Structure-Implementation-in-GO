package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var IPs [1]string // 1 IP for 1 slaves


type KeyValue struct {
    Key   string
    Value int
}
type WordCount struct {
	word  string
	count int
}


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
	
	for i := 0; i < 1; i++ { // 4 is the number of slaves
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

		println(result)






		println("------------------------------------")

	}
	f, err := os.Create("data.txt")


	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()


	// Input data (list of sentences)
	input := []string{
		result,
	}

	// Split input into chunks for mapping
	chunks := splitInput(input, 2)

	// Map stage
	mapOutput := performMap(chunks)

	// Reduce stage
	reduceOutput := performReduce(mapOutput)


	_, err2 := f.WriteString(result)


	// Print the word count results
	for _, wc := range reduceOutput {

		// fmt.Printf("%s: %d\n", wc.word, wc.count)

		f.WriteString("\r")


		count:= strconv.Itoa(wc.count)

		
		f.WriteString(wc.word + " : " +count)



		// str := wc.word + " : " +  string(wc.count)

		// f.WriteString(str)
	}




	if err2 != nil {
		fmt.Println(err2)
	}

}


// Splits the input data into chunks for mapping
func splitInput(input []string, numChunks int) [][]string {
	chunkSize := (len(input) + numChunks - 1) / numChunks
	chunks := make([][]string, numChunks)
	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(input) {
			end = len(input)
		}
		chunks[i] = input[start:end]
	}
	return chunks
}

// Map function
func mapFunc(input string) []WordCount {
	words := strings.Fields(input)
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	var wordCounts []WordCount
	for word, count := range counts {
		wordCounts = append(wordCounts, WordCount{word, count})
	}
	return wordCounts
}

// Performs the mapping phase
func performMap(chunks [][]string) []WordCount {
	var mapOutput []WordCount
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for _, chunk := range chunks {
		go func(chunk []string) {
			defer wg.Done()
			for _, input := range chunk {
				mapOutput = append(mapOutput, mapFunc(input)...)
			}
		}(chunk)
	}
	wg.Wait()
	return mapOutput
}

// Reduce function
func reduceFunc(input []WordCount) WordCount {
	word := input[0].word
	count := 0
	for _, wc := range input {
		count += wc.count
	}
	return WordCount{word, count}
}

// Performs the reducing phase
func performReduce(mapOutput []WordCount) []WordCount {
	wordCounts := make(map[string][]WordCount)
	for _, wc := range mapOutput {
		word := wc.word
		wordCounts[word] = append(wordCounts[word], wc)
	}
	var reduceOutput []WordCount
	var wg sync.WaitGroup
	wg.Add(len(wordCounts))
	for _, counts := range wordCounts {
		go func(counts []WordCount) {
			defer wg.Done()
			reduceOutput = append(reduceOutput, reduceFunc(counts))
		}(counts)
	}
	wg.Wait()
	return reduceOutput
}
