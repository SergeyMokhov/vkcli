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
        "value": "5.92"
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

func TestApi_SendRequest_AndRetry(t *testing.T) {
	requestCounter := 0
	fakeResponse := captchaNeededResponse
	actualRequestFirst := &http.Request{}
	var r2Body string
	req := fakeVkRequest{&VkRequestBase{
		Values:                url.Values{"testparam": []string{"valuetest"}},
		MethodStr:             "testMethod",
		ResponseStructPointer: &fakeVkResponse{}}}

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

	api.SendVkRequestAndRetryOnCaptcha(&req)

	secondRequestParams, err := url.ParseQuery(r2Body)

	assert.Nil(t, err)
	assert.EqualValues(t, "583944678566", secondRequestParams.Get("captcha_sid"))
	assert.EqualValues(t, "ABC123", secondRequestParams.Get("captcha_key"))
	assert.EqualValues(t, 2, requestCounter)
	assert.EqualValues(t, "valuetest", req.Values.Get("testparam"))
	assert.EqualValues(t, "5.92", req.Values.Get("v"))
	assert.EqualValues(t, "1", req.Values.Get("https"))
	assert.EqualValues(t, "123", req.Values.Get("access_token"))
	assert.EqualValues(t, "/testMethod", actualRequestFirst.RequestURI)
}
