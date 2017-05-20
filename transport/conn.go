package transport

import "net"

// Conn represents a generic network connection.
type Conn interface {
	// Read reads data from the connection.
	Read(b []byte) (n int, err error)

	// Write writes data to the connection.
	Write(b []byte) (n int, err error)

	// Close closes the connection.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
}
