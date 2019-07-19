package requests

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"testing"
//)
//
//func TestDeleteShouldSetUserId(t *testing.T) {
//	req := Delete(101)
//	require.Equal(t, "101", req.Values.Get("user_id"))
//}
//
//func TestFriendsDelete_PerformSetsRequestUri(t *testing.T) {
//	mock := NewMockRequestSender().SetResponse("friends.delete", SuccessFriendDeleted)
//	defer mock.Shutdown()
//
//	delResponse, err := Perform(mock.Api)
//	require.Nil(t, err)
//	assert.Equal(t, "/friends.delete", mock.LastRequest.RequestURI)
//	assert.Equal(t, 1, Success)
//	assert.Equal(t, 1, FriendDeleted)
//}
//
//func TestFriendsDelete_PerformInRequestDeleted(t *testing.T) {
//	mock := NewMockRequestSender().SetResponse("friends.delete", SuccessInRequestDeclined)
//	defer mock.Shutdown()
//
//	delResponse, err := Perform(mock.Api)
//	require.Nil(t, err)
//	assert.Equal(t, 1, Success)
//	assert.Equal(t, 1, InRequestDeleted)
//}
//
//func TestFriendsDelete_PerformOutRequestDeleted(t *testing.T) {
//	mock := NewMockRequestSender().SetResponse("friends.delete", SuccessOutRequestCancelled)
//	defer mock.Shutdown()
//
//	delResponse, err := Perform(mock.Api)
//	require.Nil(t, err)
//	assert.Equal(t, 1, Success)
//	assert.Equal(t, 1, OutRequestDeleted)
//}
//
//func TestFriendsDelete_PerformAccessDenied(t *testing.T) {
//	mock := NewMockRequestSender().SetResponse("friends.delete", FriedDeleteFailureAccessDenied)
//	defer mock.Shutdown()
//
//	delResponse, err := Perform(mock.Api)
//	require.Nil(t, err)
//	assert.Equal(t, 0, Success)
//	assert.Equal(t, 0, OutRequestDeleted)
//	assert.Equal(t, 0, InRequestDeleted)
//	assert.Equal(t, 0, FriendDeleted)
//
//	assert.Equal(t, 15, delResponse.ErrorCode)
//	assert.Equal(t, "Access denied: No friend or friend request found.", delResponse.ErrorMsg)
//}
