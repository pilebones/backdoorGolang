package cli

import (
	"fmt"
	"github.com/pilebones/backdoorGolang/core/socket"
)

type Context struct {
	Target *socket.TargetWrapper
	UseListenMode bool
	UseDebugMode bool
	UseVerboseMode bool
}

/** Convert context struct type as string */
func (c Context) PrettyString() string {

	return fmt.Sprintf(`Host : %s
Port : %d
IPv4 : %v
IPv6 : %v
isListenMode : %t
isVerboseMode : %t
isDebugMode : %t
`,
		c.Target.Host,
		c.Target.Port,
		c.Target.Ipv4,
		c.Target.Ipv6,
		c.UseListenMode,
		c.UseVerboseMode,
		c.UseDebugMode,
	)
}
