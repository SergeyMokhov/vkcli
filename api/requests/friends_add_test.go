package requests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const friendsAddResponse = `{
 "response": 1
}`

func TestAdd(t *testing.T) {
	req := FriendsAdd(123, "add me please", AsFollower)
	require.EqualValues(t, "123", req.Values.Get("user_id"))
	require.EqualValues(t, "add me please", req.Values.Get("text"))
	require.EqualValues(t, "1", req.Values.Get("follow"))
	require.EqualValues(t, "friends.add", req.MethodStr)
}

func TestFriendsAddRequest_CouldBeSentAndResponseIsParsible(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.add", friendsAddResponse)
	defer mock.Shutdown()
	request := FriendsAdd(123, "add me please", AsFriend)

	err := mock.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsAddResponse)

	require.True(t, ok)
	assert.Equal(t, "/friends.add", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, response.Response)
}
