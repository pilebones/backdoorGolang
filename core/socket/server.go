package socket

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"container/list"
	"bytes"
)

const (
	BUFFER_SIZE = 1024
)

/** Server structure */
type ServerProvider struct {
	Host string
	Port int
	UseDebugMode bool
	Clients *list.List
}

/** Init server instance */
func CreateServer(host string, port int, useDebugMode bool) ServerProvider {
	server := new(ServerProvider)
	server.Host = host
	server.Port = port
	server.UseDebugMode = useDebugMode
	server.Clients = list.New()

	return * server
}

func (s ServerProvider) Start() {

	in := make(chan string) // Pipe between all clients, this pipe relay all message receive from clients to all clients
	go IOHandler(in, s.Clients) // Init channel for clients (async management)

	// Listen on all <host>:<port>
	serverAddr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		panic(err.Error())
	}

	// Registration socket closure process at the end of func (close socket properly)
	defer listener.Close()

	// Infinite loop (until CTRL+C)
	for {
		fmt.Println("[SERVER][" + serverAddr + "] Waiting for client...")

		// accept connection on port
		connection, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}

		fmt.Print("[SERVER][" + serverAddr + "] Connection receive from ", connection.RemoteAddr(), "\n")
		go ClientHandler(connection, in, s.Clients)
	}
}

/**
 * Manage asynchronously IO from all connected clients
 *
 * @param Incoming - A pipe between each client to relay message
 * @param clients - A list of connected client
 */
func IOHandler(Incoming <- chan string, clients *list.List) {
	for {
		// fmt.Println("IOHandler : Waiting for input")
		input := <- Incoming
		// fmt.Println("IOHandler : Handling ", input)
		for element := clients.Front(); element != nil; element = element.Next() {
			client := element.Value.(Client)
			client.Incoming <-input
		}
	}
}

func ClientHandler(connection net.Conn, ch chan string, clients *list.List) {
	newClient := & Client{ make(chan string), ch, connection, make(chan bool), clients }
	go ClientSender(newClient)
	go ClientReceiver(newClient)
	clients.PushBack(*newClient) // Register client to server list of connected clients

	ch <- fmt.Sprintf("Another client as joined the server %s\n", connection.RemoteAddr().String())
}

func ClientReceiver(client *Client) {
	clientId := client.Connection.RemoteAddr()

	buffer := make([]byte, BUFFER_SIZE)
	for client.Read(buffer) {
		if (bytes.Equal(buffer, []byte("/quit"))) { // Logout Instruction
			client.Close()
			break;
		}

		message := string(buffer)
		fmt.Printf("ClientReader receiver : " + clientId.String() + " : " + message)

		client.Outgoing <- fmt.Sprintf("[%s] %s", clientId.String(), message)
		for i:= 0; i < BUFFER_SIZE; i++ {
			buffer[i] = 0x00; // Char End-Line
		}
	}

	client.Outgoing <- fmt.Sprintf("[%s] has left the server", clientId.String())
}

func ClientSender(client *Client) {
	clientId := client.Connection.RemoteAddr()
	for {
		select {
			// Message standard reçu par le client
			case buffer := <- client.Incoming:
				count := 0
				// Add End-line char to buffer before sending
				for i := 0; i < len(buffer); i++ {
					if buffer[i] == 0x00 {
						break
					}
					count++
				}
				fmt.Println("[SERVER] Sending to ", clientId.String(), ":")
				fmt.Print("[MESSAGE] " + string(buffer))
				fmt.Println("[MESSAGE] Size of payload : ", count)
				client.Connection.Write([]byte(buffer)[0:count]) // Flush message to client socket
			// Logout instruction
			case <- client.Quit:
				fmt.Println("Client ", clientId.String(), " quitting")
				client.Connection.Close()
				break

		}
	}
}



func (s ServerProvider) manageClient(conn net.Conn) {
	// infinite loop (until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
}
