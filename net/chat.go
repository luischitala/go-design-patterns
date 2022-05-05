package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Client chan<- string

var (
	incomingClients = make(chan Client)
	leavingClients  = make(chan Client)
	messages        = make(chan string)
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "Port")
)

//Client1 -> Server -> Handle Connection(Client1)

func Handleconnection(conn net.Conn) {
	defer conn.Close()
	message := make(chan string)
	go MessageWrite(conn, message)

	//Provide a name to the client
	clientName := conn.RemoteAddr().String()
	message <- fmt.Sprintf("Welcome to the server, your name is %s\n", clientName)
	//Use the global channel
	messages <- fmt.Sprintf("New client is here, name %s\n", clientName)
	incomingClients <- message

	inputMessage := bufio.NewScanner(conn)
	//Scan the messages to send them
	for inputMessage.Scan() {
		//Bring all the sent messages
		messages <- fmt.Sprintf("%s:%s\n", clientName, inputMessage.Text())
	}

	leavingClients <- message
	messages <- fmt.Sprintf("%s said goodbye!", clientName)
}

func MessageWrite(conn net.Conn, messages <-chan string) {
	for message := range messages {
		//Connection will be in charge of writing the messages
		fmt.Fprintln(conn, message)
	}
}

func Broadcast() {
	clients := make(map[Client]bool)
	for {
		//Use multiplex
		select {
		//Global channel
		case message := <-messages:
			for client := range clients {
				client <- message
			}
			//Add a new client to the map
		case newClient := <-incomingClients:
			clients[newClient] = true
			//Erase the clients that have left the channel
		case leavingClient := <-leavingClients:
			delete(clients, leavingClient)
			close(leavingClient)
		}

	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	//If ther's an error finish the connection
	if err != nil {
		log.Fatal(err)
	}
	go Broadcast()
	//Create loop to listen all the connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go Handleconnection(conn)
	}
}
