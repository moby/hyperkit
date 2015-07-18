// Package gist8065433 contains helpers for getting list of IPs on the machine.
package gist8065433

import "net"

// GetPublicIps returns a string slice of non-loopback IPs.
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
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip4 := ipNet.IP.To4()
			if ip4 == nil || ip4.IsLoopback() {
				continue
			}
			publicIps = append(publicIps, ipNet.IP.String())
		}
	}
	return publicIps, nil
}

// GetAllIps returns a string slice of all IPs.
func GetAllIps() (ips []string, err error) {
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
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip4 := ipNet.IP.To4()
			if ip4 == nil {
				continue
			}
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips, nil
}
