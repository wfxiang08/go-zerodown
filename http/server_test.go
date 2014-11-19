package http

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second / 2)
	})

	Convey("Test server", t, func() {
		Convey("Test relaunch", func() {
			addr := RandAddr()
			go ListenAndServe(addr, mux)
			time.Sleep(time.Second / 2)

			fd, err := DupFd()
			So(err, ShouldBeNil)

			err = Close()
			So(err, ShouldBeNil)
			err = Wait(time.Second)
			So(err, ShouldBeNil)

			go FdListenAndServe(fd, addr, mux)
			time.Sleep(time.Second / 2)

			err = Close()
			So(err, ShouldBeNil)
			err = Wait(time.Second)
			So(err, ShouldBeNil)
		})
	})

	Convey("Test server tls", t, func() {
		Convey("Test relaunch", func() {
			addr := RandAddr()
			go ListenAndServeTLS(addr, "certs/server.pem", "certs/server.key", mux)
			time.Sleep(time.Second / 2)

			fd, err := DupFd()
			So(err, ShouldBeNil)

			err = Close()
			So(err, ShouldBeNil)
			err = Wait(time.Second)
			So(err, ShouldBeNil)

			go FdListenAndServeTLS(fd, addr, "certs/server.pem", "certs/server.key", mux)
			time.Sleep(time.Second / 2)

			err = Close()
			So(err, ShouldBeNil)
			err = Wait(time.Second)
			So(err, ShouldBeNil)
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
