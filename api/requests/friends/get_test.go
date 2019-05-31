package friends

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api"
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
	fgr := Get()
	actual := fgr.Method()
	require.EqualValues(t, "friends.get", actual)
}

func TestFriendsGetRequest_UrlValues(t *testing.T) {
	fgr := Get()
	actual := fgr.UrlValues()
	require.NotNil(t, actual)
}

func TestFriendsGetRequest_SetUserId(t *testing.T) {
	fgr := Get().SetUserId(9000)
	actual := fgr.UrlValues().Get("user_id")
	require.EqualValues(t, "9000", actual)
}

func TestFriendsGetRequest_SetOrder(t *testing.T) {
	fgr := Get().SetOrder(Name)
	actual := fgr.UrlValues().Get("order")
	require.EqualValues(t, "name", actual)
}

func TestFriendsGetRequest_SetListId(t *testing.T) {
	fgr := Get().SetListId(7)
	actual := fgr.UrlValues().Get("list_id")
	require.EqualValues(t, "7", actual)
}

func TestFriendsGetRequest_SetCount(t *testing.T) {
	fgr := Get().SetCount(10000)
	actual := fgr.UrlValues().Get("count")
	require.EqualValues(t, "10000", actual)
}

func TestFriendsGetRequest_SetOffset(t *testing.T) {
	fgr := Get().SetOffset(4)
	actual := fgr.UrlValues().Get("offset")
	require.EqualValues(t, "4", actual)
}

func TestFriendsGetRequest_SetFields(t *testing.T) {
	fgr := Get().SetFields(Nickname, Sex, Domain)
	actual := fgr.UrlValues().Get("fields")
	require.EqualValues(t, "nickname,sex,domain", actual)
}

func TestFriendsGetRequest_Perform(t *testing.T) {
	mock := api.NewMockApi().SetResponse("friends.get", fakeFriendsGetResponse)
	defer mock.Shutdown()

	userlist, err := Get().SetOrder(Name).Perform(mock.Api)
	require.Nil(t, err)
	assert.EqualValues(t, "/friends.get", mock.LastRequest.RequestURI)
	assert.EqualValues(t, 1, userlist.Response.Count)
	require.EqualValues(t, 1, len(userlist.Response.Items))
	require.EqualValues(t, 12345, userlist.Response.Items[0].Id)
}

func TestFriendsGetRequest_PerformReturnsErrorResponse(t *testing.T) {
	mock := api.NewMockApi().SetResponse("friends.get", errorResponse)
	defer mock.Shutdown()

	userlist, err := Get().SetOrder(Name).Perform(mock.Api)
	require.Nil(t, err)
	require.EqualValues(t, 0, len(userlist.Response.Items))
	require.EqualValues(t, 5, userlist.Error.ErrorCode)
	require.EqualValues(t, "User authorization failed: no access_token passed.", userlist.Error.ErrorMsg)
}
