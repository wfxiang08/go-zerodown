package zerodown

import (
	"net"
)

type conn struct {
	net.Conn
	listener *Listener
}

func newConn(c net.Conn, listener *Listener) conn {
	return conn{
		Conn:     c,
		listener: listener,
	}
}

func (w conn) Close() error {
	w.listener.dec()
	return w.Conn.Close()
}
