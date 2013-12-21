package main

import (
	"net"

	. "gist.github.com/5286084.git"
)

func main() {
	netInterfaces, err := net.Interfaces()
	CheckError(err)
	for _, netInterface := range netInterfaces {
		addrs, err := netInterface.Addrs()
		CheckError(err)
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil && !ip4.IsLoopback() {
					println(netInterface.Name, addr.String())
				}
			}
		}
	}
}
