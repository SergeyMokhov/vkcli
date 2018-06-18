package vkcli_test

import (
	"crypto/tls"
	"github.com/SergeyMokhov/vkcli"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewTokenListener(t *testing.T) {
	startTokenListenerAndWaitStarted(t, "First and the only.")
}

func TestShouldFindFreePort(t *testing.T) {
	startTokenListenerAndWaitStarted(t, "First")
	startTokenListenerAndWaitStarted(t, "Second")
	startTokenListenerAndWaitStarted(t, "Third")
}

func TestShouldStopListener(t *testing.T) {
	tl := startTokenListenerAndWaitStarted(t, "Listener to stop")
	err := tl.Stop()
	if err != nil {
		t.Errorf("Failed to stop listener. %v", err)
	}
	_, postErr := unsafeHttpsPost(tl.Addr())
	if postErr == nil {
		t.Errorf("Successfully posted to listener afer calling Stop()")
	}
}

func unsafeHttpsPost(addr string) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(100 * time.Millisecond),
	}

	resp, err = client.Post("https://"+addr, "text/html", strings.NewReader("Hello?"))
	return
}

func startTokenListenerAndWaitStarted(t *testing.T, errMsg string) (listener *vkcli.TokenListener) {
	timeout := 10 * time.Millisecond

	tl, err := vkcli.NewTokenListener()
	if err != nil {
		t.Fatalf("Failed to start Token Listener. %v %v", errMsg, err)
	}

	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		_, err := unsafeHttpsPost(tl.Addr())
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("%v. Token Listener with addres: %v, did not start in %v. Error: %v", errMsg,
			tl.Addr(), timeout, err)
	}
	return tl
}
