package client

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"testing"
)

var deletedUserJson = `{
  "response": {
    "count": 3,
    "items": [
      {
        "id": 1,
        "first_name": "Potato",
        "last_name": "Salad",
        "is_closed": false,
        "can_access_closed": true,
        "nickname": "",
        "bdate": "20.2.1985",
        "online": 0
      },
      {
        "id": 2,
        "first_name": "DELETED",
        "last_name": "",
        "deactivated": "deleted",
        "online": 0
      },
      {
        "id": 3,
        "first_name": "Guacamole",
        "last_name": "Spread",
        "deactivated": "banned",
        "online": 0
      }
    ]
  }
}
`

func TestNewVk(t *testing.T) {
	vk := NewVk(&oauth2.Token{})

	require.NotNil(t, vk.api)
}

func TestVk_DeleteFriend_ShouldPrintLineOnSuccessfullDeletion(t *testing.T) {
	//testCases := []struct {
	//	name             string
	//	vkServerResponse string
	//	expectedOutput   string
	//}{
	//	{"Friend deleted from list", requests.SuccessFriendDeleted, "successfully deleted form friend list"},
	//	{"Declined incoming request", requests.SuccessInRequestDeclined, "Successfully declined friend request"},
	//	{"Cancelled outgoing request", requests.SuccessOutRequestCancelled, "Successfully cancelled friend request"},
	//	//{"Deleted friend suggestion", friends.} Was not able to find response example
	//	{"Access Denied", requests.FriedDeleteFailureAccessDenied, "Error removing/declining request"},
	//}
	//
	//var captured string
	//monkey.Patch(fmt.Printf,
	//	func(format string, a ...interface{}) (n int, err error) {
	//		captured = format
	//		return 0, nil
	//	})
	//defer monkey.UnpatchAll()
	//
	//mockApi := requests.NewMockRequestSender()
	//vk := newVkFromMockApi(mockApi)
	//
	//for _, tc := range testCases {
	//	t.Run(tc.name, func(t *testing.T) {
	//		captured = ""
	//		mockApi.SetResponse("friends.delete", tc.vkServerResponse)
	//
	//		vk.DeleteFriend(999)
	//
	//		require.Contains(t, captured, tc.expectedOutput)
	//	})
	//}
}

//func Test_RemoveDeletedFriendsShouldCallDeleteOnlyForUsersWithDeletedFlag(t *testing.T) {
//
//	//TODO finish this test and fix all others. Just make the function accept array of users.
//	//Test by feeding in json?++
//	//Create mock inside of friends package?
//	//Redesign using some sort of pattern?
//	//Export GetResponseValue struct?
//}
