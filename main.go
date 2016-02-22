package main

import (
	// "fmt"
	"github.com/pilebones/backdoorGolang/core/cli"
	"fmt"
)

func main() {

	// Parsing arguments before running
	var context cli.Context = cli.InitFlags()

	// Display target resolution before create socket provider
	if (context.IsHostIsResolved) {
		fmt.Printf("Target resolved to %s\n", context.Ipv4)
	}

	if context.UseListenMode {
		fmt.Printf("Init server mode : feature not fully implemented yet work in progress\n")
		fmt.Printf("Listening on %s:%d\n", context.Host, context.Port)
	} else {
		fmt.Printf("Init client mode : feature not fully implemented yet work in progress\n")

	}
}
