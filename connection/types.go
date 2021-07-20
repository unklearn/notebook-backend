package connection

import (
	"net"
)

// Proxy interface to net.Conn
type Channel interface {
	net.Conn
	GetId() string
}
