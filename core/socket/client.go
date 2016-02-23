package socket

import (
	"net"
	"container/list"
	"fmt"
)

/** Client structure */
type Client struct {
	Incoming chan string
	Outgoing chan string
	Connection net.Conn
	Quit chan bool
	Clients *list.List
}

/** */
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

func (c Client) Close() {
	c.Quit <- true
	c.Connection.Close()
	c.RemoveMe()
	fmt.Println("Client disconnected")
}

func (c Client) RemoveMe() {
	for element := c.Clients.Front(); element != nil; element = element.Next() {
		client := element.Value.(Client)
		if c.Connection.RemoteAddr() == client.Connection.RemoteAddr() {
			c.Clients.Remove(element)
		}
	}
}