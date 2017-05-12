// Package tuntap provides access to TUN/TAP device.
package tuntap

import (
	"io"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// TUN and TAP represent the two device types.
const (
	TUN = water.TUN
	TAP = water.TAP
)

// Device represents a TUN/TAP device.
type Device struct {
	Index int
	Type  int
	io.ReadWriteCloser
}

// New creates a new TUN/TAP device.
func New(devType int) (*Device, error) {
	iface, err := water.New(water.Config{
		DeviceType: water.DeviceType(devType),
	})
	if err != nil {
		return nil, err
	}

	link, err := netlink.LinkByName(iface.Name())
	if err != nil {
		return nil, err
	}

	return &Device{link.Attrs().Index, devType, iface}, nil
}

// NewTUN creates a new TUN device.
func NewTUN() (*Device, error) {
	return New(TUN)
}

// NewTAP creates a new TAP device.
func NewTAP() (*Device, error) {
	return New(TAP)
}

// Attrs returns the link attributes.
func (dev *Device) Attrs() *netlink.LinkAttrs {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return nil
	}

	return link.Attrs()
}

// Up changes the state of the interface to UP.
func (dev *Device) Up() error {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return err
	}

	return netlink.LinkSetUp(link)
}

// Down changes the state of the interface to DOWN.
func (dev *Device) Down() error {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return err
	}

	return netlink.LinkSetDown(link)
}

// SetMTU sets the maximum transmission unit of the device.
func (dev *Device) SetMTU(mtu int) error {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return err
	}

	return netlink.LinkSetMTU(link, mtu)
}

// Addrs returns all IPv4 and IPv6 addresses assigned to the interface.
func (dev *Device) Addrs() ([]*Addr, error) {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return nil, err
	}
	list, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return nil, err
	}

	addrs := make([]*Addr, len(list))
	for i := range list {
		a := Addr(list[i])
		addrs[i] = &a
	}
	return addrs, nil
}

// AddAddr adds the specified IP address to the interface.
func (dev *Device) AddAddr(addr *Addr) error {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return err
	}

	a := netlink.Addr(*addr)
	return netlink.AddrAdd(link, &a)
}

// DelAddr removes the specified IP address from the interface.
func (dev *Device) DelAddr(addr *Addr) error {
	link, err := netlink.LinkByIndex(dev.Index)
	if err != nil {
		return err
	}

	a := netlink.Addr(*addr)
	return netlink.AddrDel(link, &a)
}
