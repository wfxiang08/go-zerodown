package zerodown

import (
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestListener(t *testing.T) {
	Convey("Test listener", t, func() {
		Convey("Smoking", func() {
			addr := RandAddr()
			sync := make(chan int)
			l, err := Listen("tcp", addr)
			So(err, ShouldBeNil)
			go func() {
				time.Sleep(time.Second / 10)
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					t.Fatal(err)
				}
				conn.Close()
				sync <- 1
			}()
			conn, err := l.Accept()
			So(err, ShouldBeNil)
			conn.Close()
			<-sync
			fd, err := l.DupFd()
			So(err, ShouldBeNil)
			err = l.Close()
			So(err, ShouldBeNil)
			err = l.Wait(time.Second)
			So(err, ShouldBeNil)

			l, err = FdListen(fd)
			So(err, ShouldBeNil)
			go func() {
				time.Sleep(time.Second / 10)
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					t.Fatal(err)
				}
				conn.Close()
				sync <- 1
			}()
			conn, err = l.Accept()
			So(err, ShouldBeNil)
			conn.Close()
			<-sync
		})

		Convey("Test wait", func() {
			Convey("Wait before close", func() {
				addr := RandAddr()
				l, err := Listen("tcp", addr)
				So(err, ShouldBeNil)
				err = l.Wait(time.Second)
				So(err, ShouldNotBeNil)
				l.Close()
				err = l.Wait(time.Second)
				So(err, ShouldBeNil)
			})

			Convey("Wait connection", func() {
				addr := RandAddr()
				sync := make(chan int)
				l, err := Listen("tcp", addr)
				So(err, ShouldBeNil)
				go func() {
					conn, err := l.Accept()
					if err != nil {
						t.Fatal(err)
					}
					defer conn.Close()
					sync <- 1
					time.Sleep(time.Second / 2)
				}()
				conn, err := net.Dial("tcp", addr)
				So(err, ShouldBeNil)
				defer conn.Close()
				<-sync
				l.Close()
				err = l.Wait(time.Second)
				So(err, ShouldBeNil)
			})

			Convey("Wait timeout", func() {
				addr := RandAddr()
				sync := make(chan int)
				l, err := Listen("tcp", addr)
				So(err, ShouldBeNil)
				go func() {
					conn, err := l.Accept()
					if err != nil {
						t.Fatal(err)
					}
					defer conn.Close()
					sync <- 1
					time.Sleep(time.Second / 2)
				}()
				conn, err := net.Dial("tcp", addr)
				So(err, ShouldBeNil)
				defer conn.Close()
				<-sync
				l.Close()
				err = l.Wait(time.Second / 10)
				So(err, ShouldNotBeNil)
				err = l.Wait(time.Second)
				So(err, ShouldBeNil)
			})
		})
	})
}

func RandAddr() string {
	s := httptest.NewServer(nil)
	s.Close()
	u, err := url.Parse(s.URL)
	if err != nil {
		panic(err)
	}
	return u.Host
}
