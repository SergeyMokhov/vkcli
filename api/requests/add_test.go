package requests

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"testing"
//)
//
//const friendsAddresponse = `{
//  "response": 1
//}`
//
//func TestAdd(t *testing.T) {
//	req := Add(123, "add me please", AsFollower)
//	require.EqualValues(t, "123", req.Values.Get("user_id"))
//	require.EqualValues(t, "add me please", req.Values.Get("text"))
//	require.EqualValues(t, "1", req.Values.Get("follow"))
//	require.EqualValues(t, "friends.add", req.MethodStr)
//}
//
//func TestFriendsAddRequest_Perform(t *testing.T) {
//	mock := NewMockRequestSender().SetResponse("friends.add", friendsAddresponse)
//	defer mock.Shutdown()
//
//	response, err := Perform(mock.Api)
//	require.Nil(t, err)
//	assert.Nil(t, Error)
//	assert.EqualValues(t, 1, Response)
//}
//
//func TestFriendsAddRequest_PerformSetsUri(t *testing.T) {
//	mock := NewMockRequestSender().SetResponse("friends.add", friendsAddresponse)
//	defer mock.Shutdown()
//
//	resp, err := Perform(mock.Api)
//	require.Nil(t, err)
//	assert.Equal(t, "/friends.add", mock.LastRequest.RequestURI)
//	assert.Equal(t, 1, Response)
//}
