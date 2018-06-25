package auth

import (
	"gitlab.com/g00g/vkcli/tools"
	"golang.org/x/oauth2"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"errors"
	"fmt"
)

var dataFolder = "./data"

type TokenListener struct {
	token  *oauth2.Token
	server *http.Server
}

func NewTokenListener() (*TokenListener, error) {
	tl := TokenListener{}
	tl.token = nil
	srv, err := startServer()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to start token listener: %s", err))
	}
	tl.server = srv
	return &tl, nil
}

func startServer() (srv *http.Server, err error) {
	host := "localhost,127.0.0.1"
	validFrom := ""
	validFor := 365 * 24 * time.Hour
	isCA := true
	rsaBits := 2048
	ecdsaCurve := ""
	cert := filepath.Join(dataFolder, "cert.pem")
	key := filepath.Join(dataFolder, "key.pem")

	if !certificateExists() {
		err := GenerateCert(host, validFrom, validFor, isCA, rsaBits, ecdsaCurve, dataFolder)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot generate Certificate: %s", err))
		}
	}

	srv = &http.Server{}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listening port %v", err))
	}

	srv.Addr = "localhost:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	go func() {
		if err := srv.ServeTLS(listener, cert, key); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error staring or stopping server: %v", err)
		}
	}()

	return srv, nil
}

func certificateExists() bool {
	var cert bool
	var key bool

	cert = tools.FileExists(filepath.Join(dataFolder, "cert.pem"))
	key = tools.FileExists(filepath.Join(dataFolder, "key.pem"))

	return cert && key
}

func (tl *TokenListener) Stop() error {
	return tl.server.Close()
}

func (tl *TokenListener) Addr() (addr string) {
	return tl.server.Addr
}

func (*TokenListener) Token() (token *oauth2.Token, err error) {
	return nil, nil
	//TODO return real token
}
