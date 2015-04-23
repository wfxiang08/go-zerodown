package zerodown

import "net"

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

func (c conn) Close() error {
	c.listener.dec()
	return c.Conn.Close()
}
