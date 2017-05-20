// Package udp provides transport over UDP.
package udp

import "net"

// Conn implements the transport.Conn interface for UDP.
type Conn struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	*net.UDPConn
}

// New creates a connection with the remote peer over UDP.
func New(local, remote string) (*Conn, error) {
	localAddr, err := net.ResolveUDPAddr("udp", local)
	if nil != err {
		return nil, err
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", remote)
	if nil != err {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if nil != err {
		return nil, err
	}

	return &Conn{localAddr, remoteAddr, conn}, nil
}

// Read reads data from the peer.
func (c *Conn) Read(b []byte) (int, error) {
	n, addr, err := c.ReadFromUDP(b)
	if err != nil {
		return n, err
	}

	// Discard packets not from the peer
	if addr.String() != c.remoteAddr.String() {
		return 0, nil
	}

	return n, err
}

// Write writes data to the peer.
func (c *Conn) Write(b []byte) (int, error) {
	return c.WriteToUDP(b, c.remoteAddr)
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.localAddr
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}
