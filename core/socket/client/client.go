package client

import (
	"github.com/pilebones/backdoorGolang/core/socket"
)


/** Server structure */
type ClientProvider struct {
	Token string
	SocketWrapper socket.SocketWrapper
}

type ServerProvider struct {
	// SocketWrapper socket.SocketWrapper
	Host string
	Port int
	UseDebugMode bool
	Clients *list.List
}
