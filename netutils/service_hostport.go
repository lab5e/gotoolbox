package netutils

import (
	"fmt"
	"log"
	"net"
)

// ServiceHostPort returns the listener address as a host:port string. If
// the listener address points at a loopback address it will return the
// address of the loopback adapter.
func ServiceHostPort(addr net.Addr) string {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if ok {
		if tcpAddr.IP.IsUnspecified() {
			publicIP, err := FindPublicIPv4()
			if err != nil {
				log.Printf("Unable to determine public IP: %v", err)
			}
			return fmt.Sprintf("%s:%d", publicIP.String(), tcpAddr.Port)
		}
		return tcpAddr.String()
	}

	udpAddr, ok := addr.(*net.UDPAddr)
	if ok {
		if udpAddr.IP.IsUnspecified() {
			publicIP, err := FindPublicIPv4()
			if err != nil {
				log.Printf("Unable to determine public IP: %v", err)
			}
			return fmt.Sprintf("%s:%d", publicIP.String(), udpAddr.Port)
		}
		return udpAddr.String()
	}
	log.Printf("Listen address isn't TCP or UDP (%T)", addr)
	return ""
}
