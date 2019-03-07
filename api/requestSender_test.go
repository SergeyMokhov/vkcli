package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestApi_SendRequest_NoRetry(t *testing.T) {
	fakeResponse := `{"response": {"count": 1,"items": [{"id": 12345,"first_name": "Alexander","last_name": "Ivanov",
"bdate": "20.2.1985","online": 0}]}}`
	actualRequest := &http.Request{}

	req := fakeVkRequest{&VkRequestBase{Values: url.Values{"testparam": []string{"valuetest"}}, MethodStr: "testMethod"}}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequest = r
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fakeResponse)
	}))
	defer ts.Close()

	api := getTestApi(t, ts)

	api.SendVkRequestAndRetryOnCaptcha(&req)

	assert.EqualValues(t, "valuetest", req.Values.Get("testparam"))
	assert.EqualValues(t, "5.92", req.Values.Get("v"))
	assert.EqualValues(t, "1", req.Values.Get("https"))
	assert.EqualValues(t, "123", req.Values.Get("access_token"))
	assert.EqualValues(t, "/testMethod", actualRequest.RequestURI)
}

func TestNewDummyVkRequest(t *testing.T) {
	dr := NewVkRequestBase("methodName", &fakeVkResponse{})
	require.NotNil(t, dr.Values)
	require.EqualValues(t, 0, len(dr.UrlValues()))
	require.NotNil(t, dr.ResponseType())
	require.EqualValues(t, "methodName", dr.Method())
}

func TestAddSolvedCaptcha(t *testing.T) {
	dr := NewVkRequestBase("methodName", &fakeVkResponse{})
	vkErr := vkErrors.Error{ErrorInfo: vkErrors.ErrorInfo{CaptchaSid: "98874562"}}
	addSolvedCaptcha(dr, &vkErr, "zQ7a")

	require.EqualValues(t, "98874562", dr.Values.Get("captcha_sid"))
	require.EqualValues(t, "zQ7a", dr.Values.Get("captcha_key"))
}

func TestAddDefaultVkRequestParamsOnlyOnce(t *testing.T) {
	vkr := NewVkRequestBase("tst", &fakeVkResponse{})

	addDefaultParams(vkr, "token")
	addDefaultParams(vkr, "token")
	addDefaultParams(vkr, "token2")

	params := vkr.UrlValues()
	assert.EqualValues(t, 1, len(params["access_token"]))
	assert.EqualValues(t, "token", params.Get("access_token"))
	assert.EqualValues(t, 1, len(params["v"]))
	assert.EqualValues(t, 1, len(params["https"]))
}

func getTestApi(t *testing.T, ts *httptest.Server) *Api {
	api := NewInstance(&oauth2.Token{AccessToken: "123"})
	baseUrl, urlParseErr := url.Parse(ts.URL)
	require.Nil(t, urlParseErr)
	api.BaseUrl = baseUrl
	return api
}
