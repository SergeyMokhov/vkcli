package requests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteShouldSetUserId(t *testing.T) {
	req := FriendsDelete(101)
	require.Equal(t, "101", req.Values.Get("user_id"))
}

//TODO Provide generic function to run similar tests
func TestFriendsDeleteRequest_CouldBeSentAndResponseIsParsible_FriendDeleted(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.delete", SuccessFriendDeleted)
	defer mock.Shutdown()
	request := FriendsDelete(123)

	err := mock.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsDeleteResponse)

	require.True(t, ok)
	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, response.Response.Success)
	assert.Equal(t, 1, response.Response.FriendDeleted)
}

func TestFriendsDeleteRequest_CouldBeSentAndResponseIsParsible_InRequest(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.delete", SuccessInRequestDeclined)
	defer mock.Shutdown()
	request := FriendsDelete(123)

	err := mock.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsDeleteResponse)

	require.True(t, ok)
	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, response.Response.Success)
	assert.Equal(t, 1, response.Response.InRequestDeleted)
}

func TestFriendsDeleteRequest_CouldBeSentAndResponseIsParsible_OutRequest(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.delete", SuccessOutRequestCancelled)
	defer mock.Shutdown()
	request := FriendsDelete(123)

	err := mock.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsDeleteResponse)

	require.True(t, ok)
	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, response.Response.Success)
	assert.Equal(t, 1, response.Response.OutRequestDeleted)
}

func TestFriendsDeleteRequest_CouldBeSentAndResponseIsParsible_Error(t *testing.T) {
	mock := NewMockRequestSender().SetResponse("friends.delete", FriedDeleteFailureAccessDenied)
	defer mock.Shutdown()
	request := FriendsDelete(123)

	err := mock.SendVkRequestAndRetryOnCaptcha(request)
	require.Nil(t, err)

	response, ok := request.ResponseStructPointer.(*FriendsDeleteResponse)

	require.True(t, ok)
	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
	assert.Equal(t, 0, response.Response.Success)
	assert.Equal(t, 0, response.Response.FriendDeleted)
	assert.Equal(t, 0, response.Response.InRequestDeleted)
	assert.Equal(t, 0, response.Response.OutRequestDeleted)
	assert.Equal(t, 15, response.ErrorCode)
	assert.Equal(t, "Access denied: No friend or friend request found.", response.ErrorMsg)
}
