package cli

import (
	"fmt"
	"os"
	"net"
	"github.com/spf13/pflag"
	"github.com/pilebones/backdoorGolang/core/common"
)

var (
	host    	= pflag.StringP("host", "h", "localhost", "Set hostname to use")
	port    	= pflag.IntP("port", "p", 9876, "Set port number to use")
	isListenMode 	= pflag.BoolP("listen", "l", false, "Enable listen mode (server socket mode)")
	isVerboseMode 	= pflag.BoolP("verbose", "v", false, "Enable mode verbose")
	isDebugMode 	= pflag.BoolP("debug", "d", false, "Enable mode debug")
	isVersionMode 	= pflag.BoolP("version", "", false, "Display version number")
)

/** Print data when debug mode is enabled */
func DisplayAsDebug(message string) {
	if UseDebugMode() {
		fmt.Println(message)
	}
}

/** Return true if mode use from parameter */
func UseMode(mode *bool) bool {
	var useMode bool = false
	if *mode {
		useMode = true
	}
	return useMode
}

/** Return true if debug mode is enabled */
func UseDebugMode() bool {
	return UseMode(isDebugMode)
}

/** Return true if listen mode is enabled */
func UseListenMode() bool {
	return UseMode(isListenMode)
}

/** Return true if the user want to see the program version */
func UseVersionMode() bool {
	return UseMode(isVerboseMode)
}
/** Parse arguments and check value is allowed */
func InitFlags() {
	pflag.Parse()
	if UseVersionMode() {
		fmt.Printf("%s : Build %s Version %s\nAuthor : %s (see: %s)", common.PRODUCT_NAME, common.BUILD, common.VERSION, common.AUTHOR, common.CONTACT)
		os.Exit(0)
	}

	DisplayAsDebug("Debug mode enabled")

	// Resolv hostname as net.IP
	DisplayAsDebug(fmt.Sprintf("Host : %s", *host))
	var ip net.IP = net.ParseIP(*host)
	if (ip == nil) { // if argument isn't IP => check if the hostname can be resolved
		ips, err := net.LookupIP(*host)
		if err != nil {
			fmt.Errorf("Couln't resolv hostname \"%s\", error : %v", *host, err)
			os.Exit(1);
		} else {
			if 1 < len(ips) {
				var ipv6 net.IP = ips[0]
				var ipv4 net.IP = ips[1]
				DisplayAsDebug(fmt.Sprintf("IPv6: %s\nIPv4: %s", ipv6.String(), ipv4.String()))
			} else {
				fmt.Errorf("Couln't resolv from hostname IPv4  \"%s\"", *host)
				os.Exit(1)
			}
		}
	} else {
		DisplayAsDebug(fmt.Sprintf("IP %s", ip.String()))
	}

	DisplayAsDebug(fmt.Sprintf("Port : %d", *port))
	DisplayAsDebug(fmt.Sprintf("isListenMode : %t", UseListenMode()))
	DisplayAsDebug(fmt.Sprintf("isVerboseMode : %t", UseVersionMode()))
}