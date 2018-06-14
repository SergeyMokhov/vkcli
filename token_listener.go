package vkcli

import (
	"golang.org/x/oauth2"
	"log"
	"net"
	"net/http"
	"strconv"
	"os"
)

var data = "./data"

type TokenListener struct {
	token  *oauth2.Token
	server *http.Server
}

func NewTokenListener() (*TokenListener, error) {
	tl := TokenListener{}
	tl.token = nil
	srv := startServer()
	tl.server = srv
	return &tl, nil
}

func startServer() (srv *http.Server) {
	srv = &http.Server{}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to start listener %v", err)
	}

	srv.Addr = "localhost:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Printf("%v", err)
		}
	}()

	return srv
}

func fileExists(path string) (bool) {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func certificateExist() (bool) {
	var cert bool
	var key bool

	cert = fileExists(data + "/cert.pem")
	key = fileExists(data + "/key.pem")

	return cert && key
}

func (tl *TokenListener) Stop() (error) {
	return tl.server.Close()
}

func (tl *TokenListener) Addr() (addr string) {
	return tl.server.Addr
}

func (*TokenListener) Token() (token *oauth2.Token, err error) {
	return nil, nil
	//TODO return real token
}
