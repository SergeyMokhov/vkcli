package friends

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api"
	"testing"
)

var successFriendDeleted = `{"response":{"success":1,"friend_deleted":1}}`
var successInRequestDeleted = `{"response":{"success":1,"in_request_deleted":1}}`
var successOutRequestDeleted = `{"response":{"success":1,"out_request_deleted":1}}`
var failureAccessDenied = `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found.","request_params":[{"key":"oauth","value":"1"},{"key":"method","value":"friends.delete"},{"key":"https","value":"1"},{"key":"user_id","value":"2916112"},{"key":"v","value":"5.92"}]}}`

func TestDeleteShouldSetUserId(t *testing.T) {
	req := Delete(101)
	require.Equal(t, "101", req.Values.Get("user_id"))
}

func TestFriendsDelete_PerformSetsRequestUri(t *testing.T) {
	mock := api.NewMockApi(successFriendDeleted)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
	assert.Equal(t, 1, delResponse.Response.Success)
	assert.Equal(t, 1, delResponse.Response.FriendDeleted)
}

func TestFriendsDelete_PerformInRequestDeleted(t *testing.T) {
	mock := api.NewMockApi(successInRequestDeleted)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, 1, delResponse.Response.Success)
	assert.Equal(t, 1, delResponse.Response.InRequestDeleted)
}

func TestFriendsDelete_PerformOutRequestDeleted(t *testing.T) {
	mock := api.NewMockApi(successOutRequestDeleted)
	defer mock.Shutdown()

	delResponse, err := Delete(1).Perform(mock.Api)
	require.Nil(t, err)
	assert.Equal(t, 1, delResponse.Response.Success)
	assert.Equal(t, 1, delResponse.Response.OutRequestDeleted)
}

func TestFriendsDelete_PerformAccessDenied(t *testing.T) {
	mock := api.NewMockApi(failureAccessDenied)
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
