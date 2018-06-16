package vkcli

import (
	"golang.org/x/oauth2"
	"log"
	"net"
	"net/http"
	"strconv"
	"os"
	"time"
	"path/filepath"
)

var dataFolder = "./data"

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
	host := "localhost,127.0.0.1"
	validFrom := ""
	validFor := 365 * 24 * time.Hour
	isCA := true
	rsaBits := 2048
	ecdsaCurve := ""
	cert := filepath.Join(dataFolder, "cert.pem")
	key := filepath.Join(dataFolder, "key.pem")

	if !certificateExist() {
		GenerateCert(host, validFrom, validFor, isCA, rsaBits, ecdsaCurve, dataFolder)
	}

	srv = &http.Server{}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to start listener %v", err)
	}

	srv.Addr = "localhost:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	go func() {
		if err := srv.ServeTLS(listener, cert, key); err != nil {
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

	cert = fileExists(filepath.Join(dataFolder, "cert.pem"))
	key = fileExists(filepath.Join(dataFolder, "key.pem"))

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
