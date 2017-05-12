package tuntap

import (
	"net"

	"github.com/vishvananda/netlink"
)

// An Addr represents an IP address and its associated routing prefix.
type Addr netlink.Addr

// ParseCIDR parses s as a CIDR notation IP address and prefix length.
func ParseCIDR(s string) (*Addr, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	ipNet.IP = ip

	return &Addr{IPNet: ipNet}, nil
}

// Network returns the network name, "ip+net".
func (addr *Addr) Network() string {
	return addr.IPNet.Network()
}

// String returns the CIDR notation of the address.
func (addr *Addr) String() string {
	return addr.IPNet.String()
}
