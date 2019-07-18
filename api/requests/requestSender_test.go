package requests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"net/url"
	"testing"
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
	assert.EqualValues(t, "5.95", req.Values.Get("v"))
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
