package socket

import (
	"net"
)

type SocketWrapperInterface interface {
	Token() string
	IsHostIsResolved() bool
	Host() string
	Ipv4() net.IP
	Ipv6() net.IP
	Port() int
}

type SocketWrapper struct {
	Token string
	IsHostIsResolved bool
	Host string
	Ipv4 net.IP
	Ipv6 net.IP
	Port int
}