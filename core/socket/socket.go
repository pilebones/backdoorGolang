package socket

import (
	"net"
	"fmt"
)

type TargetWrapper struct {
	Host string
	Port int
	Ipv4 net.IP
	Ipv6 net.IP
}

type TargetWrapperInterface interface {
	ResolvHost()
	HostCanBeResolv() bool

}

func (c *TargetWrapper) ResolvHost() {
	// Resolv hostname as net.IP
	ip := net.ParseIP(c.Host)
	if ip == nil { // Host isn't IP => check if the hostname can be resolved
		ips, err := net.LookupIP(c.Host)
		if err != nil {
			panic(fmt.Sprintf(`Couln't resolv hostname "%s"`, c.Host))
		}

		if (0 == len(ips)) {
			panic(fmt.Sprintf(`Unknow error after resolving hostname "%s"`, c.Host))
		}

		c.Ipv4 = ips[1]
		c.Ipv6 = ips[0]
	} else {
		c.Ipv4 = ip.To4()
		c.Ipv6 = ip.To16()
	}

}

func (c *TargetWrapper) HostCanBeResolv() bool {
	var hostHasBeenResolved bool = false
	defer func() {
		if err := recover(); err != nil {
			hostHasBeenResolved = false
		}
	}()

	c.ResolvHost()
	hostHasBeenResolved = true;

	// Return true if Host has been resolv
	return hostHasBeenResolved
}