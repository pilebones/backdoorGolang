package socket

import (
	"net"
	"container/list"
	"fmt"
)

/** Client structure */
type Client struct {
	Incoming chan string // Message Channel : Input
	Outgoing chan string // Message Channel : Output
	Connection net.Conn // Socket between client and server
	Quit chan bool // Message Channel which contains the state of connection
	Clients *list.List // List of clients connected from server
}

/**
 * Read data from client socket
 * @param buffer
 */
func (c Client) Read(buffer []byte) bool {
	bytesRead, error := c.Connection.Read(buffer)
	if (error != nil) {
		c.Close()
		fmt.Errorf("Unable to read buffer client, connection close (%v)\n", error)
		return false
	}
	fmt.Println(bytesRead, " bytes read")
	return true
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