package http

import (
	"crypto/tls"
	"github.com/googollee/go-zerodown"
	"net/http"
	"time"
)

var server *Server

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

func DupFd() (int, error) {
	return server.DupFd()
}

func Close() error {
	return server.Close()
}

func Wait(timeout time.Duration) error {
	return server.Wait(timeout)
}

type Server struct {
	server   *http.Server
	listener *zerodown.Listener
}

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

func (s *Server) ListenAndServe() error {
	return s.server.Serve(s.listener)
}

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

func (s *Server) DupFd() (int, error) {
	return s.listener.DupFd()
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Wait(timeout time.Duration) error {
	return s.listener.Wait(timeout)
}
