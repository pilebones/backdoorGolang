package server

import (
	"bytes"
	"container/list"
	"fmt"
	"net"
	"strings"
)

/** Client structure */
type Client struct {
	Incoming   chan string // Message Channel : Input
	Outgoing   chan string // Message Channel : Output
	Connection net.Conn    // Socket between client and server
	Quit       chan bool   // Message Channel which contains the state of connection
	Clients    *list.List  // List of clients connected from server
}

/**
 * Read data from client socket
 * @param buffer
 */
func (c Client) Read(buffer []byte) bool {
	_, error := c.Connection.Read(buffer)
	if error != nil {
		c.Close()
		fmt.Errorf("Unable to read buffer client, connection close (%v)\n", error)
		return false
	}
	// fmt.Println(bytesRead, " bytes read")
	return true
}

/**
 * Send data from client socket
 * @param buffer - []byte
 * @return int - nb bytes sent
 */
func (c Client) Send(buffer []byte) int {
	count := 0
	// Add End-line char to buffer before sending
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == 0x00 {
			break
		}
		count++
	}
	message := strings.TrimSpace(string(buffer)) // Just for debug and server output
	fmt.Printf("[MESSAGE][%s] %s {%d bytes}\n", c.Connection.RemoteAddr().String(), message, count)
	c.Connection.Write([]byte(buffer)[0:count]) // Flush message to client socket

	return count
}

/**
 * Send string message from client socket
 * @param message - string
 */
func (c Client) SendMessage(message string) {
	if len(message) != 0 {
		buffer := bytes.NewBufferString(message).Bytes()
		c.Send(buffer)
	}
}

/**
 * Manage client logout
 */
func (c Client) Close() {
	c.Quit <- true
	c.Connection.Close()
	c.RemoveMe()
	fmt.Println("Client disconnected")
}

/**
 * Remove client from the list of connected users
 */
func (c Client) RemoveMe() {
	for element := c.Clients.Front(); element != nil; element = element.Next() {
		client := element.Value.(Client)
		if c.Connection.RemoteAddr() == client.Connection.RemoteAddr() {
			c.Clients.Remove(element)
		}
	}
}

/**
 * Getting uniq identifier of a client
 */
func (c Client) GetId() string {
	return c.Connection.RemoteAddr().String()
}
