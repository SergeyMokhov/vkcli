package requests

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// Just regular Api, but sends requests to mock http server instead of real one.
// Also records last request and can respond with response you want.
// You can set different responses for different methods.
// Keeps track of number of requests received per method.
// For usage see tests
type MockRequestSender struct {
	VkRequestSender *VkRequestSender
	Server          *httptest.Server
	LastRequest     *http.Request
	Response        map[string]string
	requestCounter  map[string]int
}

func (m *MockRequestSender) SetResponse(requestUrlPath string, response string) *MockRequestSender {
	m.Response["/"+requestUrlPath] = response
	return m
}

func (m *MockRequestSender) NumberOfRequestsReceived(requestUrlPath string) int {
	return m.requestCounter["/"+requestUrlPath]
}

func (m *MockRequestSender) Shutdown() {
	m.Server.Close()
}

//Call the shutdown function once you are done with the mock
func NewMockRequestSender() *MockRequestSender {
	mock := &MockRequestSender{}
	mock.Response = make(map[string]string)
	mock.requestCounter = make(map[string]int)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.requestCounter[r.RequestURI] = mock.requestCounter[r.RequestURI] + 1
		mock.LastRequest = r
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", mock.Response[r.RequestURI])
	}))

	mock.Server = server
	requestSender := NewVkRequestSender(&oauth2.Token{AccessToken: "000"})
	baseUrl, _ := url.Parse(server.URL)
	requestSender.BaseUrl = baseUrl
	mock.VkRequestSender = requestSender
	return mock
}
