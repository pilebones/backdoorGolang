package socket

import (
	"fmt"
	"net"
	"bufio"
	"strings"
)

type ServerProvider struct {
	Host string
	IsHostIsResolved bool
	Ipv4 net.IP
	Ipv6 net.IP
	Port int
	Clients map[string]SocketWrapper
}

func CreateServer(wrapper SocketWrapper) ServerProvider {
	server := new(ServerProvider)
	server.Host = wrapper.Host
	server.Port = wrapper.Port
	server.IsHostIsResolved = wrapper.IsHostIsResolved
	server.Ipv4 = wrapper.Ipv4
	server.Ipv6 = wrapper.Ipv6
	server.Clients = map[string]SocketWrapper{}

	return * server
}

func (s ServerProvider) Start() {
	// listen on all <host>:<port>
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		panic(err.Error())
	}

	// accept connection on port
	conn, err := listener.Accept()
	if err != nil {
		panic(err.Error())
	}

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
