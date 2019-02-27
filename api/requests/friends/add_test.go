package friends

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api"
	"testing"
)

const friendsAddresponse = `{
  "response": 1
}`

func TestAdd(t *testing.T) {
	req := Add(123, "add me please", AsFollower)
	require.EqualValues(t, "123", req.Values.Get("user_id"))
	require.EqualValues(t, "add me please", req.Values.Get("text"))
	require.EqualValues(t, "1", req.Values.Get("follow"))
	require.EqualValues(t, "friends.add", req.MethodStr)
}

func TestFriendsAddRequest_Perform(t *testing.T) {
	mock := api.NewMockApi(friendsAddresponse)
	defer mock.Shutdown()

	response, err := Add(505, "tst", AsFriend).Perform(mock.Api)
	require.Nil(t, err)
	assert.Nil(t, response.Error)
	assert.EqualValues(t, 1, response.Response)
}

func TestFriendsAddRequest_PerformSetsUri(t *testing.T) {
	mock := api.NewMockApi(friendsAddresponse)
	defer mock.Shutdown()

	resp, err := Add(1, "sfsdf", AsFriend).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, "/friends.add", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, resp.Response)
}
