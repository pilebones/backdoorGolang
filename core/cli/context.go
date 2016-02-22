package cli

import (
	"fmt"
	"net"
	"github.com/pilebones/backdoorGolang/core/socket"
)

type Context struct {
	Host string
	IsHostIsResolved bool
	Ipv4 net.IP
	Ipv6 net.IP
	Port int
	UseListenMode bool
	UseDebugMode bool
	UseVerboseMode bool
}

/** Convert context struct type as string */
func (c Context) ToString() string {

	return fmt.Sprintf(`Host : %s
Port : %d
IP v4 : %s
IP v6 : %s
isHostIsResolved : %t
isListenMode : %t
isVerboseMode : %t
isDebugMode : %t
`,
		c.Host,
		c.Port,
		c.Ipv4.String(),
		c.Ipv6.String(),
		c.IsHostIsResolved,
		c.UseListenMode,
		c.UseVerboseMode,
		c.UseDebugMode,
	)
}

func (c Context) CastToSocketWrapper() socket.SocketWrapper {
	var wrapper = new(socket.SocketWrapper)
	wrapper.Host = c.Host
	wrapper.Port = c.Port
	wrapper.IsHostIsResolved = c.IsHostIsResolved
	wrapper.Ipv4 = c.Ipv4
	wrapper.Ipv6 = c.Ipv6

	return * wrapper
}
