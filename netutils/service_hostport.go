package netutils

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
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
				logrus.Fatal("Unable to determine public IP")
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
				logrus.Fatal("Unable to determine public IP")
			}
			return fmt.Sprintf("%s:%d", publicIP.String(), udpAddr.Port)
		}
		return udpAddr.String()
	}
	logrus.WithField("addr", addr).Fatalf("Listen address isn't TCP or UDP (%T)", addr)
	return ""
}
