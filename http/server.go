package http

import (
	"crypto/tls"
	"github.com/googollee/go-zerodown"
	"net/http"
	"time"
)

var server *Server

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections.
// See http.ListenAndServe for more details.
func ListenAndServe(addr string, handler http.Handler) error {
	if addr == "" {
		addr = ":80"
	}
	var err error
	server, err = New(&http.Server{
		Addr:    addr,
		Handler: handler,
	})
	if err != nil {
		return err
	}
	return server.ListenAndServe()
}

// FdListenAndServe listens on the address which described by file descrptor fd and then calls Serve with handler to handle requests on incoming connections.
func FdListenAndServe(fd int, addr string, handler http.Handler) error {
	var err error
	server, err = FromFd(fd, &http.Server{
		Addr:    addr,
		Handler: handler,
	})
	if err != nil {
		return err
	}
	return server.ListenAndServe()
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it expects HTTPS connections.
// See http.ListenAndServeTLS for more details.
func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	var err error
	server, err = New(&http.Server{
		Addr:    addr,
		Handler: handler,
	})
	if err != nil {
		return err
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}

// FdListenAndServeTLS acts identically to FdListenAndServe, except that it expects HTTPS connections.
func FdListenAndServeTLS(fd int, addr, certFile, keyFile string, handler http.Handler) error {
	var err error
	server, err = FromFd(fd, &http.Server{
		Addr:    addr,
		Handler: handler,
	})
	if err != nil {
		return err
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}

// DupFd returns the integer Unix file descriptor duplicated from Listener.
func DupFd() (int, error) {
	return server.DupFd()
}

// IsClosed returns a boolean to indicate whether Listener is closed.
func IsClosed() bool {
	return server.IsClosed()
}

// Close will closed Listener. It won't Accept connection any more.
func Close() error {
	return server.Close()
}

// Wait waits all connections created by Listener closed.
// If after timeout reach but not all connection closed, it will return time out error.
func Wait(timeout time.Duration) error {
	return server.Wait(timeout)
}

// A Server with closable listener.
type Server struct {
	server   *http.Server
	listener *zerodown.Listener
}

// New creates a Server with zerodown.Listener, and use defines parameters in server for running an HTTP server.
// It will listen on server.Addr
func New(server *http.Server) (*Server, error) {
	listener, err := zerodown.Listen("tcp", server.Addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		server:   server,
		listener: listener,
	}, nil
}

// New creates a Server with zerodown.Listener, and use defines parameters in server for running an HTTP server.
// It will listen on the adress described in file descriptor fd.
func FromFd(fd int, server *http.Server) (*Server, error) {
	listener, err := zerodown.FdListen(fd)
	if err != nil {
		return nil, err
	}
	return &Server{
		server:   server,
		listener: listener,
	}, nil
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
func (s *Server) ListenAndServe() error {
	return s.server.Serve(s.listener)
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming TLS connections.
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	config := &tls.Config{}
	if s.server.TLSConfig != nil {
		*config = *s.server.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(s.listener, config)
	return s.server.Serve(tlsListener)
}

// DupFd returns the integer Unix file descriptor duplicated from Listener.
func (s *Server) DupFd() (int, error) {
	return s.listener.DupFd()
}

// IsClosed returns a boolean to indicate whether Listener is closed.
func (s *Server) IsClosed() bool {
	return s.listener.IsClosed()
}

// Close will closed Listener. It won't Accept connection any more.
func (s *Server) Close() error {
	return s.listener.Close()
}

// Wait waits all connections created by Listener closed.
// If after timeout reach but not all connection closed, it will return time out error.
func (s *Server) Wait(timeout time.Duration) error {
	return s.listener.Wait(timeout)
}
