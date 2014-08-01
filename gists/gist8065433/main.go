package gist8065433

import (
	"fmt"
	"net"

	. "github.com/shurcooL/go/gists/gist5286084"
)

// TODO: Better interface. Perhaps a slice of strings.
// Get a string of non-loopback IPs.
func Something() (out string) {
	netInterfaces, err := net.Interfaces()
	CheckError(err)
	for _, netInterface := range netInterfaces {
		addrs, err := netInterface.Addrs()
		CheckError(err)
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil && !ip4.IsLoopback() {
					//netInterface.Name
					out += fmt.Sprintln(ipnet.IP.String())
				}
			}
		}
	}
	return out
}

func main() {
	fmt.Print(Something())
}
