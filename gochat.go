package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
)

const defaultPort = "8082"
const defaultIP = "127.0.0.1"
const addressConnector = ":"
const senderPrefix = "Me: "
const receiverPrefix = "Received: "

func main() {
	isServer := false
	address := ""
	if len(os.Args) > 1 {
		address, isServer = getAddres()
	} else {
		address = defaultIP + addressConnector + defaultPort
	}

	var connection net.Conn
	var listener net.Listener

	if isServer {
		// listen on all interfaces
		listener, _ = net.Listen("tcp", address)
		connection, _ = listener.Accept()
	} else { // it's a client
		// connect to this socket
		connection, _ = net.Dial("tcp", address)
	}

	textChannel := make(chan string)
	messageReceivedFlag := make(chan bool)

	go closeResources(listener, connection)
	go receiveMessages(connection, messageReceivedFlag)
	go printCurrentText(textChannel, messageReceivedFlag)

	// disable input buffering to be able to read one char at the time, it's
	// needed to be able caching the characters writen when receiving a message
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	fmt.Println("To close correctly gochat, use CTRL + C ... ")
	fmt.Println()

	// read from console and send messages
	var byteRead = make([]byte, 1)
	for { 
		fmt.Print(senderPrefix)
		os.Stdin.Read(byteRead)
		var inputMessage string
		for ( string(byteRead) != "\n") {
			textChannel <- string(byteRead)
			inputMessage += string(byteRead)
			os.Stdin.Read(byteRead)
		}

		// send message via socket
		fmt.Fprintf(connection, inputMessage + "\n")
		inputMessage = ""
		messageReceivedFlag <- false
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
		} else if(strings.ToLower(os.Args[1]) == "-s") {
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
		isErr1 := strings.ToLower(os.Args[1]) == "-s"
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

func receiveMessages(connection net.Conn, messageReceivedFlag chan bool) {
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Printf("%c[2K", 27)
		fmt.Print("\r" + receiverPrefix + string(message))
		messageReceivedFlag <- true
	}
}

func printCurrentText(textChannel chan string, messageReceivedFlag chan bool) {
	var writenUntilNow string
	for {
		select {
		case charReceived := <- textChannel:
			writenUntilNow += charReceived
			charReceived = ""
		case isMEssageReceived := <- messageReceivedFlag:
			if isMEssageReceived {
				fmt.Print(senderPrefix + writenUntilNow)
			}
			writenUntilNow = ""
		default:
		}
	}
}

func closeResources(listener net.Listener, connection net.Conn) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	if listener != nil { listener.Close() }
	if connection != nil { connection.Close() }
	fmt.Println("Program has been interrupted, closing resources...")
	os.Exit(0)
}

func exitProgram() {
	// print usage
	fmt.Println("bash: gochat: " + os.Args[1] + ": invalid option")
	fmt.Println("gochat: usage: gochat [-s] [dir ip] [port]")

	os.Exit(0)
}