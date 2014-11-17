package zerodown

import (
	"errors"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

type fileListener interface {
	net.Listener
	File() (f *os.File, err error)
}

type Listener struct {
	fileListener
	quit          chan struct{}
	closeLocker   sync.RWMutex
	isClosed      bool
	counterLocker sync.RWMutex
	counter       int
}

func New(l net.Listener) (*Listener, error) {
	fl, ok := l.(fileListener)
	if !ok {
		return nil, errors.New("listener doesn't have file")
	}
	return &Listener{
		fileListener: fl,
		quit:         make(chan struct{}),
	}, nil
}

func Listen(lnet, laddr string) (*Listener, error) {
	l, err := net.Listen(lnet, laddr)
	if err != nil {
		return nil, err
	}
	return New(l)
}

func FdListen(fd int) (*Listener, error) {
	f := os.NewFile(uintptr(fd), "listen socket")
	defer f.Close()
	l, err := net.FileListener(f)
	if err != nil {
		return nil, err
	}
	return New(l)
}

func (l *Listener) Accept() (net.Conn, error) {
	c, err := l.fileListener.Accept()
	if err != nil {
		return nil, err
	}

	l.inc()

	return newConn(c, l), nil
}

func (l *Listener) DupFd() (int, error) {
	f, err := l.fileListener.File()
	if err != nil {
		return 0, err
	}

	return syscall.Dup(int(f.Fd()))
}

func (l *Listener) IsClosed() bool {
	l.closeLocker.RLock()
	defer l.closeLocker.RUnlock()
	return l.isClosed
}

func (l *Listener) Close() error {
	l.closeLocker.Lock()
	defer l.closeLocker.Unlock()
	l.isClosed = true
	return l.fileListener.Close()
}

func (l *Listener) Wait(timeout time.Duration) error {
	if !l.IsClosed() {
		return errors.New("not closed")
	}
	if l.count() == 0 {
		return nil
	}

	select {
	case <-l.quit:
	case <-time.After(timeout):
		return errors.New("time out")
	}

	return nil
}

func (l *Listener) inc() {
	l.counterLocker.Lock()
	defer l.counterLocker.Unlock()
	l.counter++
}

func (l *Listener) dec() {
	l.counterLocker.Lock()
	defer l.counterLocker.Unlock()
	l.counter--
	if l.IsClosed() && l.counter == 0 {
		l.quit <- struct{}{}
	}
}

func (l *Listener) count() int {
	l.counterLocker.RLock()
	defer l.counterLocker.RUnlock()
	return l.counter
}
