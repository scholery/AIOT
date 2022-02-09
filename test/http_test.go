package test

import (
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"testing"

	"github.com/sirupsen/logrus"
)

func IndexHandler(w http.ResponseWriter, request *http.Request) {
	logrus.Info("request", request)
	fmt.Fprintln(w, "hello world", request.RequestURI)
}

func HttpServer() bool {
	var server *http.Server
	server = &http.Server{Addr: fmt.Sprintf(":%d", 7000)}
	http.HandleFunc("/test", IndexHandler)
	server.ListenAndServe()
	http.HandleFunc("/test2", IndexHandler)
	signal.Notify(nil, syscall.SIGINT)
	return true
}

func TestHttpServer(t *testing.T) {
	HttpServer()
	for {

	}
}
