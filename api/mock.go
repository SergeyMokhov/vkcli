package api

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
)

//Just regular Api, but sends requests to mock http server instead of real one.
// Also records last request and can respond with response you want.
type MockApi struct {
	Api            *Api
	Server         *httptest.Server
	LastRequest    *http.Request
	Response       map[string]string
	requestCounter map[string]int
}

func (m *MockApi) SetResponse(requestUrlPath string, response string) *MockApi {
	m.Response["/"+requestUrlPath] = response
	return m
}

func (m *MockApi) NumberOfRequestsReceived(requestUrlPath string) int {
	return m.requestCounter["/"+requestUrlPath]
}

func (m *MockApi) Shutdown() {
	m.Server.Close()
}

//Call the shutdown function once you are done with the mock
func NewMockApi() *MockApi {
	mock := &MockApi{}
	mock.Response = make(map[string]string)
	mock.requestCounter = make(map[string]int)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.requestCounter[r.RequestURI] = mock.requestCounter[r.RequestURI] + 1
		mock.LastRequest = r
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", mock.Response[r.RequestURI])
	}))

	mock.Server = server
	requestSender := NewInstance(&oauth2.Token{AccessToken: "000"})
	baseUrl, _ := url.Parse(server.URL)
	requestSender.BaseUrl = baseUrl
	mock.Api = requestSender
	return mock
}

type fakeVkRequest struct {
	*VkRequestBase
}

func (fg *fakeVkRequest) UrlValues() url.Values {
	return fg.Values
}

func (fg *fakeVkRequest) Method() string {
	return fg.MethodStr
}

func (fg *fakeVkRequest) ResponseType() vkResponse {
	return fg.ResponseStructPointer
}

type fakeVkResponse struct {
	*vkErrors.Error
}

func (fr *fakeVkResponse) GetError() *vkErrors.Error {
	return fr.Error
}
