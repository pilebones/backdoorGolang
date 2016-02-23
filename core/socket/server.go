package socket

import (
	"fmt"
	"net"
	"container/list"
	"bytes"
)

const (
	BUFFER_CLIENT_NB_MESSAGE_THRESHOLD = 200
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

	// Pipe between all clients, this pipe relay all message receive from clients to all clients
	in := make(chan string, BUFFER_CLIENT_NB_MESSAGE_THRESHOLD) // Limit number of simultaneous message with a threshold to avoid server DOS
	go IOHandler(in, s.Clients) // Init clients channel for message (async management)

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
		go ClientHandler(connection, in, s.Clients) // Manage all clients asychronously
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
		// Waiting incoming message from a client
		input := <- Incoming
		// fmt.Println("IOHandler : Handling ", input)
		for element := clients.Front(); element != nil; element = element.Next() {
			client := element.Value.(Client)
			client.Incoming <-input // Notify all clients for this message
		}
	}
}

/**
 * Manage I/O client from server socket
 *
 * @param connection - Socket between client and server
 * @param messageChannel - The shared bus message between all clients
 * @param clients - The list of all clients connected
 */
func ClientHandler(connection net.Conn, messageChannel chan string, clients *list.List) {
	// Create new client instance
	newClient := & Client{ make(chan string), messageChannel, connection, make(chan bool), clients }
	go ClientSender(newClient) // Manage sending message
	go ClientReceiver(newClient) // Manage receiving message
	clients.PushBack(*newClient) // Register client to server list of connected clients

	// Notify all clients for the new connection
	messageChannel <- fmt.Sprintf("Another client as joined the server %s\n", connection.RemoteAddr().String())
}

/**
 * Manage client input from server-side socket
 *
 * @param client - Client object
 */
func ClientReceiver(client *Client) {
	clientId := client.Connection.RemoteAddr()

	buffer := make([]byte, BUFFER_SIZE)
	for client.Read(buffer) { // While read data from server-side socket
		if (bytes.Equal(buffer, []byte("/quit"))) { // Logout Instruction
			client.Close()
			break;
		}

		message := string(buffer)
		fmt.Printf("ClientReader receiver : " + clientId.String() + " : " + message)

		// Sending message to client
		client.Outgoing <- fmt.Sprintf("[%s] %s", clientId.String(), message)
		for i:= 0; i < BUFFER_SIZE; i++ {
			buffer[i] = 0x00; // Char End-Line
		}
	}

	client.Outgoing <- fmt.Sprintf("[%s] has left the server", clientId.String())
}

/**
 * Manage output from server-side socket
 *
 * @param client - Client object
 */
func ClientSender(client *Client) {
	clientId := client.Connection.RemoteAddr()
	for {
		select {
			// Standard message receive by the server from a client
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
