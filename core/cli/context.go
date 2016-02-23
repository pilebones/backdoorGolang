package cli

import (
	"fmt"
	"net"
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
