package gist8065433

import "net"

// Get a string of non-loopback IPs.
func GetPublicIps() (publicIps []string, err error) {
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, ifi := range ifis {
		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil && !ip4.IsLoopback() {
					//ifi.Name
					publicIps = append(publicIps, ipnet.IP.String())
				}
			}
		}
	}
	return publicIps, nil
}
