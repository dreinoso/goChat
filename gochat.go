package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)

const defaultPort = "8080"
const defaultIP = "127.0.0.1"
const addressConnector = ":"

func main() {
	isServer := false
	address := ""
	if len(os.Args) > 1 {
		address, isServer = getAddres()
	} else {
		address = defaultIP + addressConnector + defaultPort
	}

	var connection net.Conn 
	if isServer {
		// listen on all interfaces
		listenter, _ := net.Listen("tcp", address)
	
		// accept connection on port
		connection, _ = listenter.Accept()
	} else { // it's a client
		// connect to this socket
		connection, _ = net.Dial("tcp", address)
	}

	go receiveMessages(connection)

	for { 
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(connection, text + "\n")
	}
}

func getAddres()(string, bool) {
	var address string
	isServer := false
	switch len(os.Args) {
    case 2:
    	_, err := strconv.Atoi(os.Args[1]);
    	if err == nil {
        	address = os.Args[1] + addressConnector + defaultPort
    	} else if(strings.ToLower(os.Args[1]) == "-c") {
			address = defaultIP + addressConnector + defaultPort
			isServer = true
      	} else {
			exitProgram()
		}
    case 3:
		_, err1 := strconv.Atoi(os.Args[1]);
    	_, err2 := strconv.Atoi(os.Args[2]);
      	if err1 == nil && err2 == nil {
        	address = os.Args[1] + addressConnector + os.Args[1]
      	} else {
         exitProgram()
     	}
    case 4:
      	isErr1 := strings.ToLower(os.Args[1]) == "-c"
      	_, err2 := strconv.Atoi(os.Args[2]);
      	_, err3 := strconv.Atoi(os.Args[3]);
      	if isErr1 && err2 == nil && err3 == nil {
			address = os.Args[1] + addressConnector + os.Args[1]
			isServer = true
      	} else {
         exitProgram()
      	}
    default: exitProgram()
  }

  return address, isServer
}

func receiveMessages(connection net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Printf("%c[2K", 27)
		fmt.Println("\rMessage received: " + string(message))
		var inputMsgCached []byte
		reader.Read(inputMsgCached)
		fmt.Print(inputMsgCached)
	}
}

func exitProgram() {
	// print usage
	fmt.Println("bash: gochat: " + os.Args[1] + ": invalid option")
	fmt.Println("gochat: usage: gochat [-c] [dir ip] [port]")

	os.Exit(0)
}