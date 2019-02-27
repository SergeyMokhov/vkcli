package api

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type mockApi struct {
	Api         *Api
	Server      *httptest.Server
	LastRequest *http.Request
	Response    *string
}

func (m *mockApi) SetResponse(s string) {
	m.Response = &s
}

func (m *mockApi) Shutdown() {
	m.Server.Close()
}

//Call the shutdown function once you are done with the mock
func NewMockApi(response string) *mockApi {
	mock := &mockApi{}
	mock.SetResponse(response)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.LastRequest = r
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", *mock.Response)
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