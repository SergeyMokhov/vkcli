package vk

import (
	"testing"
	"time"
	"net/http"
	"strings"
)

func TestNewTokenListener(t *testing.T) {
	startAndVerifyTokenListener(t, "First and the only.")
}

func TestShouldFindFreePort(t *testing.T) {
	startAndVerifyTokenListener(t, "First")
	startAndVerifyTokenListener(t, "Second")
	startAndVerifyTokenListener(t, "Third")
}

func TestShouldStopListener(t *testing.T) {
	tl := startAndVerifyTokenListener(t, "Listener to stop")
	err := tl.Stop()
	if err != nil {
		t.Errorf("Failed to stop listener. %v", err)
	}
	err2 := post(tl.Addr())
	if err2 == nil {
		t.Errorf("Successfully posted to listener afer calling Stop()")
	}
}

func post(addr string) (error) {
	_, err := http.Post("http://"+addr, "text/html", strings.NewReader("test"))
	return err
}

func startAndVerifyTokenListener(t *testing.T, errMsg string) (*TokenListener) {
	timeout := 10 * time.Millisecond

	tl, err := NewTokenListener()
	if err != nil {
		t.Fatalf("Failed to start Token Listener. %v %v", errMsg, err)
	}

	now := time.Now()
	deadline := now.Add(timeout)
	for time.Now().Before(deadline) {
		_, err = http.Post("http://"+tl.Addr(), "text/html", strings.NewReader("test"))
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
