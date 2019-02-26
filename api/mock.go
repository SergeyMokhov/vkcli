package api

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
)

//Call the close function once you are done with the mock
func MockApi(vkResponse string) (api *Api, serverCloseFunc func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, vkResponse)
	}))
	requestSender := NewInstance(&oauth2.Token{AccessToken: "000"})
	baseUrl, _ := url.Parse(ts.URL)
	requestSender.BaseUrl = baseUrl
	return requestSender, ts.Close
}
