package cli

import (
	"fmt"
	"os"
	"net"
	"github.com/spf13/pflag"
	"github.com/pilebones/backdoorGolang/core/common"
)

var (
	host    	= pflag.IPP("host", "h", net.ParseIP("localhost"), "Set hostname to use")
	port    	= pflag.IntP("port", "p", 9876, "Set port number to use")
	isListenMode 	= pflag.BoolP("listen", "l", false, "Enable listen mode (server socket mode)")
	isVerboseMode 	= pflag.BoolP("verbose", "v", false, "Enable mode verbose")
	isDebugMode 	= pflag.BoolP("debug", "d", false, "Enable mode debug")
	isVersionMode 	= pflag.BoolP("version", "", false, "Display version number")
)

func InitFlags() {
	pflag.Parse()
	// If no value defined
	// pflag.Lookup("isListenMode").NoOptDefVal = true
	// pflag.Lookup("isVerboseMode").NoOptDefVal = true
	// pflag.Lookup("isDebugMode").NoOptDefVal = true

	if *isVersionMode {
		fmt.Printf("%s : Build %s Version %s\nAuthor : %s (see: %s)", common.PRODUCT_NAME, common.BUILD, common.VERSION, common.AUTHOR, common.CONTACT)
		os.Exit(0)
	}

	fmt.Printf("Host : %s\n")
	fmt.Printf("Port : %d\n")
	fmt.Printf("isListenMod : %d\n", *isListenMode)
	fmt.Printf("isVerboseMode : %d\n", *isVerboseMode)
	fmt.Printf("isDebugMode : %d\n", *isDebugMode)
}