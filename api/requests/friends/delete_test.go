package friends

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api"
	"testing"
)

func TestDeleteShouldSetUserId(t *testing.T) {
	req := Delete(101)
	require.Equal(t, "101", req.Values.Get("user_id"))
}

func TestFriendsDelete_PerformSetsRequestUri(t *testing.T) {
	mock := api.NewMockApi(SuccessFriendDeleted)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, delResponse.Response.Success)
	assert.Equal(t, 1, delResponse.Response.FriendDeleted)
}

func TestFriendsDelete_PerformInRequestDeleted(t *testing.T) {
	mock := api.NewMockApi(SuccessInRequestDeleted)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, 1, delResponse.Response.Success)
	assert.Equal(t, 1, delResponse.Response.InRequestDeleted)
}

func TestFriendsDelete_PerformOutRequestDeleted(t *testing.T) {
	mock := api.NewMockApi(SuccessOutRequestDeleted)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, 1, delResponse.Response.Success)
	assert.Equal(t, 1, delResponse.Response.OutRequestDeleted)
}

func TestFriendsDelete_PerformAccessDenied(t *testing.T) {
	mock := api.NewMockApi(FriedDeleteFailureAccessDenied)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, 0, delResponse.Response.Success)
	assert.Equal(t, 0, delResponse.Response.OutRequestDeleted)
	assert.Equal(t, 0, delResponse.Response.InRequestDeleted)
	assert.Equal(t, 0, delResponse.Response.FriendDeleted)

	assert.Equal(t, 15, delResponse.ErrorCode)
	assert.Equal(t, "Access denied: No friend or friend request found.", delResponse.ErrorMsg)
}
