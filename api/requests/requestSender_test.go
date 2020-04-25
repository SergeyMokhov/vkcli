package requests

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

func TestRequseSender_SendRequest_NoRetry(t *testing.T) {
	fakeResponse := `{"response": {"count": 1,"items": [{"id": 12345,"first_name": "Alexander","last_name": "Ivanov",
"bdate": "20.2.1985","online": 0}]}}`
	testMethod := "testMethod"
	mock := NewMockRequestSender()
	mock.SetResponse(testMethod, fakeResponse)
	req := FakeVkRequest{&VkRequestBase{Values: url.Values{"testparam": []string{"valuetest"}}, MethodStr: testMethod}}
	defer mock.Shutdown()

	requestSender := mock.VkRequestSender

	requestSender.SendVkRequestAndRetryOnCaptcha(&req)

	assert.EqualValues(t, "valuetest", req.Values.Get("testparam"))
	assert.EqualValues(t, "5.103", req.Values.Get("v"))
	assert.EqualValues(t, "1", req.Values.Get("https"))
	assert.EqualValues(t, "000", req.Values.Get("access_token"))
	assert.EqualValues(t, "/testMethod", mock.LastRequest.RequestURI)
}

func TestNewDummyVkRequest(t *testing.T) {
	dr := NewVkRequestBase("methodName", &FakeVkResponse{})
	require.NotNil(t, dr.Values)
	require.EqualValues(t, 0, len(dr.UrlValues()))
	require.NotNil(t, dr.ResponseType())
	require.EqualValues(t, "methodName", dr.Method())
}

func TestAddSolvedCaptcha(t *testing.T) {
	dr := NewVkRequestBase("methodName", &FakeVkResponse{})
	vkErr := vkErrors.Error{ErrorInfo: vkErrors.ErrorInfo{CaptchaSid: "98874562"}}
	addSolvedCaptcha(dr, &vkErr, "zQ7a")

	require.EqualValues(t, "98874562", dr.Values.Get("captcha_sid"))
	require.EqualValues(t, "zQ7a", dr.Values.Get("captcha_key"))
}

func TestAddDefaultVkRequestParamsOnlyOnce(t *testing.T) {
	vkr := NewVkRequestBase("tst", &FakeVkResponse{})

	addDefaultParams(vkr, "token")
	addDefaultParams(vkr, "token")
	addDefaultParams(vkr, "token2")

	params := vkr.UrlValues()
	assert.EqualValues(t, 1, len(params["access_token"]))
	assert.EqualValues(t, "token", params.Get("access_token"))
	assert.EqualValues(t, 1, len(params["v"]))
	assert.EqualValues(t, 1, len(params["https"]))
}

func TestNewVkRequestSender_ShouldInitializeAllFields(t *testing.T) {
	api := NewVkRequestSender(&oauth2.Token{AccessToken: "123"})

	require.NotNil(t, api.client)
	require.NotNil(t, api.token)
	require.NotNil(t, api.BaseUrl)
	require.EqualValues(t, "123", api.token.AccessToken)
}

func TestRequestSender_SendRequest_AndRetry_ShouldSendCaptchaWhenRetries(t *testing.T) {
	testMethod := "testMethod"
	mock := NewMockRequestSender()
	defer mock.Shutdown()
	requestCounter := 0
	fakeResponse := captchaNeededResponse
	actualRequestFirst := &http.Request{}
	var r2Body string
	fakeRequest := FakeVkRequest{&VkRequestBase{
		Values:                url.Values{"testparam": []string{"valuetest"}},
		MethodStr:             testMethod,
		ResponseStructPointer: &FakeVkResponse{}}}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	mock.SetTestServer(testServer)
	requestSender := mock.VkRequestSender

	patch := monkey.Patch(promptForCaptcha, func(vkErr *vkErrors.Error) (answer string) {
		return "ABC123"
	})
	defer patch.Unpatch()

	requestSender.SendVkRequestAndRetryOnCaptcha(&fakeRequest)

	secondRequestParams, err := url.ParseQuery(r2Body)

	assert.Nil(t, err)
	assert.EqualValues(t, "583944678566", secondRequestParams.Get("captcha_sid"))
	assert.EqualValues(t, "ABC123", secondRequestParams.Get("captcha_key"))
	assert.EqualValues(t, 2, requestCounter)
	assert.EqualValues(t, "valuetest", fakeRequest.Values.Get("testparam"))
	assert.EqualValues(t, "5.103", fakeRequest.Values.Get("v"))
	assert.EqualValues(t, "1", fakeRequest.Values.Get("https"))
	assert.EqualValues(t, "000", fakeRequest.Values.Get("access_token"))
	assert.EqualValues(t, "/testMethod", actualRequestFirst.RequestURI)
}
