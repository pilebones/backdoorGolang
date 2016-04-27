package cli

import (
	"fmt"
	"os"

	"github.com/pilebones/backdoorGolang/core/common"
	"github.com/pilebones/backdoorGolang/core/socket"
	"github.com/spf13/pflag"
)

var (
	host          = pflag.StringP("host", "h", "localhost", "Set hostname to use")
	port          = pflag.IntP("port", "p", 9876, "Set port number to use")
	isListenMode  = pflag.BoolP("listen", "l", false, "Enable listen mode (server socket mode)")
	isVerboseMode = pflag.BoolP("verbose", "v", false, "Enable mode verbose")
	isDebugMode   = pflag.BoolP("debug", "d", false, "Enable mode debug")
	isVersionMode = pflag.BoolP("version", "V", false, "Display version number")
)

/** Print data when debug mode is enabled */
func DisplayAsDebug(message string) {
	if UseDebugMode() {
		fmt.Println(message)
	}
}

/** Return true if mode use from parameter */
func UseMode(mode *bool) bool {
	if *mode {
		return true
	}
	return false
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
	return UseMode(isVersionMode)
}

/** Return true if the user want to see the program version */
func UseVerboseMode() bool {
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
		fmt.Printf("%s : Build %s Version %f\nAuthor : %s (see: %s)", common.PRODUCT_NAME, common.BUILD, common.VERSION, common.AUTHOR, common.CONTACT)
		os.Exit(0)
	} else {
		fmt.Printf("%v", isVerboseMode)
	}

	return generateContextFromFlags()
}

/** Parse, Validate flags and generate CLI context Object from CLI Arguments */
func generateContextFromFlags() Context {

	target := new(socket.TargetWrapper)
	target.Host = *host
	target.Port = *port

	context := new(Context)
	context.Target = target
	context.UseDebugMode = UseDebugMode()
	context.UseListenMode = UseListenMode()
	context.UseVerboseMode = UseVersionMode()

	DisplayAsDebug("Debug mode enabled")

	// Resolv hostname as net.IP if possible
	if !target.HostCanBeResolv() {
		panic(fmt.Sprintf(`Invalid host : Couln't resolv "%s"`, *host))
	}

	DisplayAsDebug(context.PrettyString())

	return *context
}
