package main

import (
	"flag"
	"fmt"
	"github.com/googollee/go-zerodown"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	log.SetPrefix(fmt.Sprintf("[%d] ", syscall.Getpid()))

	var fd int
	flag.IntVar(&fd, "fd", -1, "the already-open fd to listen on")
	flag.Parse()

	server := &http.Server{Addr: ":8000"}

	var listener *zerodown.Listener
	var err error
	if fd < 0 {
		log.Println("Listening on a new fd")
		listener, err = zerodown.Listen("tcp", server.Addr)
	} else {
		log.Println("Listening to existing fd", fd)
		listener, err = zerodown.FdListen(fd)
	}
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	server.Handler = mux

	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello world\n")
	})
	mux.HandleFunc("/sleep", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Second * 10)
	})
	mux.HandleFunc("/upgrade", func(w http.ResponseWriter, req *http.Request) {
		newFd, err := listener.DupFd()
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command(os.Args[0],
			fmt.Sprintf("-fd=%d", newFd))
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		log.Println("starting cmd:", cmd.Args)
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		listener.Close()
	})

	log.Println("Serving on", server.Addr)
	err = server.Serve(listener)

	if listener.IsClosed() {
		if err := listener.Wait(time.Second * 10); err != nil {
			log.Println("wait error:", err)
		} else {
			log.Println("quit")
		}
	} else if err != nil {
		log.Fatal(err)
	}
}
