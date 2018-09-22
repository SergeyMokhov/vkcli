package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewInstance(t *testing.T) {
	api := NewInstance(&oauth2.Token{AccessToken: "123"})

	require.NotNil(t, api.client)
	require.NotNil(t, api.token)
	require.NotNil(t, api.BaseUrl)
	require.EqualValues(t, "123", api.token.AccessToken)
}

func TestApi_PerformShouldRatainUrlValuesAndAddDefaultOnes(t *testing.T) {
	fakeResponse := `{"response": {"count": 1,"items": [{"id": 12345,"first_name": "Alexander","last_name": "Ivanov",
"bdate": "20.2.1985","online": 0}]}}`
	actualRequest := &http.Request{}
	req := fakeVkRequest{values: url.Values{"testparam": []string{"valuetest"}}, method: "testMethod"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequest = r
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fakeResponse)
	}))
	defer ts.Close()

	api := NewInstance(&oauth2.Token{AccessToken: "123"})
	baseUrl, urlParseErr := url.Parse(ts.URL)
	require.Nil(t, urlParseErr)
	api.BaseUrl = baseUrl

	api.SendRequest(&req)

	assert.EqualValues(t, "valuetest", req.values.Get("testparam"))
	assert.EqualValues(t, "5.85", req.values.Get("v"))
	assert.EqualValues(t, "1", req.values.Get("https"))
	assert.EqualValues(t, "123", req.values.Get("access_token"))
	assert.EqualValues(t, "/testMethod", actualRequest.RequestURI)
}

func TestNewDummyVkRequest(t *testing.T) {
	dr := NewDummyVkRequest("methodName", &struct{}{})
	require.NotNil(t, dr.Values)
	require.EqualValues(t, 0, len(dr.UrlValues()))
	require.NotNil(t, dr.ResponseType())
	require.EqualValues(t, "methodName", dr.Method())
}

type fakeVkRequest struct {
	values url.Values
	method string
}

func (fg *fakeVkRequest) UrlValues() url.Values {
	return fg.values
}

func (fg *fakeVkRequest) Method() string {
	return fg.method
}

func (fg *fakeVkRequest) ResponseType() interface{} {
	return &struct{}{}
}
