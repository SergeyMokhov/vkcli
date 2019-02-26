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
	requestSender, close := api.MockApi(friendsAddresponse)
	defer close()

	response, err := Add(505, "tst", AsFriend).Perform(requestSender)
	require.Nil(t, err)
	assert.Nil(t, response.Error)
	assert.EqualValues(t, 1, response.Response)
}
