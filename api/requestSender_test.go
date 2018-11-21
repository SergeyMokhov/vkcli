package api

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	captchaNeededResponse = `{
  "error": {
    "error_code": 14,
    "error_msg": "Captcha needed",
    "request_params": [
      {
        "key": "oauth",
        "value": "1"
      },
      {
        "key": "method",
        "value": "friends.add"
      },
      {
        "key": "follow",
        "value": "0"
      },
      {
        "key": "https",
        "value": "1"
      },
      {
        "key": "text",
        "value": "lol"
      },
      {
        "key": "user_id",
        "value": "155633421"
      },
      {
        "key": "v",
        "value": "5.85"
      }
    ],
    "captcha_sid": "583944678566",
    "captcha_img": "https:\/\/api.vk.com\/captcha.php?sid=583944678566&s=1"
  }
}`
)

func TestNewInstance(t *testing.T) {
	api := NewInstance(&oauth2.Token{AccessToken: "123"})

	require.NotNil(t, api.client)
	require.NotNil(t, api.token)
	require.NotNil(t, api.BaseUrl)
	require.EqualValues(t, "123", api.token.AccessToken)
}

func TestApi_SendRequest_NoRetry(t *testing.T) {
	tc := []struct {
		Name string
		F    func(rb *Api, request vkRequest) (err error)
	}{
		{"sendVkRequestAndRetyOnCaptcha", sendVkRequestAndRetyOnCaptcha},
	}

	for _, test := range tc {
		t.Run(test.Name, func(t *testing.T) {
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

			api := getTestApi(t, ts)

			test.F(api, &req)

			assert.EqualValues(t, "valuetest", req.values.Get("testparam"))
			assert.EqualValues(t, "5.85", req.values.Get("v"))
			assert.EqualValues(t, "1", req.values.Get("https"))
			assert.EqualValues(t, "123", req.values.Get("access_token"))
			assert.EqualValues(t, "/testMethod", actualRequest.RequestURI)
		})
	}
}

func TestApi_SendRequest_AndRetry(t *testing.T) {
	requestCounter := 0
	fakeResponse := captchaNeededResponse
	actualRequestFirst := &http.Request{}
	var r2Body string
	req := fakeVkRequest{values: url.Values{"testparam": []string{"valuetest"}}, method: "testMethod"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter++
		switch requestCounter {
		case 1:
			actualRequestFirst = r
		case 2:
			rb, _ := ioutil.ReadAll(r.Body)
			r2Body = string(rb)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fakeResponse)
	}))
	defer ts.Close()

	patch := monkey.Patch(promptForCaptcha, func(vkErr *vkErrors.Error) (answer string) {
		return "ABC123"
	})
	defer patch.Unpatch()

	api := getTestApi(t, ts)

	sendVkRequestAndRetyOnCaptcha(api, &req)

	secondRequestParams, err := url.ParseQuery(r2Body)

	assert.Nil(t, err)
	assert.EqualValues(t, "583944678566", secondRequestParams.Get("captcha_sid"))
	assert.EqualValues(t, "ABC123", secondRequestParams.Get("captcha_key"))
	assert.EqualValues(t, 2, requestCounter)
	assert.EqualValues(t, "valuetest", req.values.Get("testparam"))
	assert.EqualValues(t, "5.85", req.values.Get("v"))
	assert.EqualValues(t, "1", req.values.Get("https"))
	assert.EqualValues(t, "123", req.values.Get("access_token"))
	assert.EqualValues(t, "/testMethod", actualRequestFirst.RequestURI)
}

func TestNewDummyVkRequest(t *testing.T) {
	dr := NewVkRequestBase("methodName", &struct{}{})
	require.NotNil(t, dr.Values)
	require.EqualValues(t, 0, len(dr.UrlValues()))
	require.NotNil(t, dr.ResponseType())
	require.EqualValues(t, "methodName", dr.Method())
}

func TestAddSolvedCaptcha(t *testing.T) {
	dr := NewVkRequestBase("methodName", &struct{}{})
	vkErr := vkErrors.Error{ErrorInfo: vkErrors.ErrorInfo{CaptchaSid: "98874562"}}
	addSolvedCaptcha(dr, &vkErr, "zQ7a")

	require.EqualValues(t, "98874562", dr.Values.Get("captcha_sid"))
	require.EqualValues(t, "zQ7a", dr.Values.Get("captcha_key"))
}

func TestAddDefaultVkRequestParamsOnlyOnce(t *testing.T) {
	vkr := NewVkRequestBase("tst", struct{}{})

	addDefaultParams(vkr, "token")
	addDefaultParams(vkr, "token")
	addDefaultParams(vkr, "token2")

	params := vkr.UrlValues()
	assert.EqualValues(t, 1, len(params["access_token"]))
	assert.EqualValues(t, "token", params.Get("access_token"))
	assert.EqualValues(t, 1, len(params["v"]))
	assert.EqualValues(t, 1, len(params["https"]))
}

type fakeVkRequest struct {
	values url.Values
	method string
}

func getTestApi(t *testing.T, ts *httptest.Server) *Api {
	api := NewInstance(&oauth2.Token{AccessToken: "123"})
	baseUrl, urlParseErr := url.Parse(ts.URL)
	require.Nil(t, urlParseErr)
	api.BaseUrl = baseUrl
	return api
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
