# go-zerodown

[![GoDoc](http://godoc.org/github.com/googollee/go-zerodown?status.svg)](http://godoc.org/github.com/googollee/go-zerodown) [![Build Status](https://travis-ci.org/googollee/go-zerodown.svg)](https://travis-ci.org/googollee/go-zerodown)

go-zerodown provides a listener which can shutdown gracefully and relaunch in another process. It mainly use in upgrading http server without completely shut down the service.

You can run the example at [`example/main.go`](https://github.com/googollee/go-zerodown/blob/master/example/main.go) like this:

```
go build example/main.go
./main &
curl http://localhost:8000/upgrade
```
