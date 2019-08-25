package requests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	fakeFriendsGetResponse = `{
  "response": {
    "count": 1,
    "items": [
      {
        "id": 12345,
        "first_name": "Alexander",
        "last_name": "Ivanov",
        "bdate": "20.2.1985",
        "online": 0,
		"is_closed": false,
		"can_access_closed": true
      }
    ]
  }
}`
	errorResponse = `{
	 "error": {
	   "error_code": 5,
	   "error_msg": "User authorization failed: no access_token passed.",
	   "request_params": [
	     {
	       "key": "oauth",
	       "value": "1"
	     },
	     {
	       "key": "method",
	       "value": "friends.get"
	     },
	     {
	       "key": "fields",
	       "value": "bdate,has_mobile,nickname"
	     },
	     {
	       "key": "https",
	       "value": "1"
	     },
	     {
	       "key": "user_id",
	       "value": "70007"
	     },
	     {
	       "key": "v",
	       "value": "5.84"
	     }
	   ]
	 }
	}`
)

func TestFriendsGetRequest_Method(t *testing.T) {
	fgr := FriendsGet()
	actual := fgr.Method()
	require.EqualValues(t, "friends.get", actual)
}

func TestFriendsGetRequest_UrlValues(t *testing.T) {
	fgr := FriendsGet()
	actual := fgr.UrlValues()
	require.NotNil(t, actual)
}

func TestFriendsGetRequest_SetUserId(t *testing.T) {
	fgr := FriendsGet().SetUserId(9000)
	actual := fgr.UrlValues().Get("user_id")
	require.EqualValues(t, "9000", actual)
}

func TestFriendsGetRequest_SetOrder(t *testing.T) {
	fgr := FriendsGet().SetOrder(Name)
	actual := fgr.UrlValues().Get("order")
	require.EqualValues(t, "name", actual)
}

func TestFriendsGetRequest_SetListId(t *testing.T) {
	fgr := FriendsGet().SetListId(7)
	actual := fgr.UrlValues().Get("list_id")
	require.EqualValues(t, "7", actual)
}

func TestFriendsGetRequest_SetCount(t *testing.T) {
	fgr := FriendsGet().SetCount(10000)
	actual := fgr.UrlValues().Get("count")
	require.EqualValues(t, "10000", actual)
}

func TestFriendsGetRequest_SetOffset(t *testing.T) {
	fgr := FriendsGet().SetOffset(4)
	actual := fgr.UrlValues().Get("offset")
	require.EqualValues(t, "4", actual)
}

func TestFriendsGetRequest_SetFields(t *testing.T) {
	fgr := FriendsGet().SetFields([]FriendsGetFields{Nickname, Sex, Domain})
	actual := fgr.UrlValues().Get("fields")
	require.EqualValues(t, "nickname,sex,domain", actual)
}

func TestFriendsGetRequest_CouldBeSentAndResponseIsProperlyParsed(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.get", fakeFriendsGetResponse)
	defer mock.Shutdown()
	request := FriendsGet().SetOrder(Name)
	err := mock.VkRequestSender.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsGetResponse)

	require.True(t, ok)
	assert.EqualValues(t, "/friends.get", mock.LastRequest.RequestURI)
	assert.EqualValues(t, 1, response.Response.Count)
	require.EqualValues(t, 1, len(response.Response.Items))
	require.EqualValues(t, 12345, response.Response.Items[0].Id)
}

func TestFriendsGetRequest_PerformReturnsErrorResponse(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.get", errorResponse)
	defer mock.Shutdown()
	request := FriendsGet().SetOrder(Name)
	err := mock.VkRequestSender.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsGetResponse)

	require.True(t, ok)
	require.Nil(t, err)
	require.EqualValues(t, 0, len(response.Response.Items))
	require.EqualValues(t, 5, response.Error.ErrorCode)
	require.EqualValues(t, "User authorization failed: no access_token passed.", response.Error.ErrorMsg)
}
