package client

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api/requests"
	"golang.org/x/oauth2"
	"reflect"
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

func TestVk_DeleteFriend_ShouldPrintLineOnSuccessfulDeletion(t *testing.T) {
	testCases := []struct {
		name             string
		vkServerResponse string
		expectedOutput   string
	}{
		{"Friend deleted from list", requests.SuccessFriendDeleted, "successfully deleted form friend list"},
		{"Declined incoming request", requests.SuccessInRequestDeclined, "Successfully declined friend request"},
		{"Cancelled outgoing request", requests.SuccessOutRequestCancelled, "Successfully cancelled friend request"},
		//{"Deleted friend suggestion", friends.} Was not able to find response example
		{"Access Denied", requests.FriedDeleteFailureAccessDenied, "Error removing/declining request"},
	}

	var captured string
	monkey.Patch(fmt.Printf,
		func(format string, a ...interface{}) (n int, err error) {
			captured = format
			return 0, nil
		})
	defer monkey.UnpatchAll()

	mockApi := requests.NewMockRequestSender()
	vk := NewVkFromMock(mockApi)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			captured = ""
			mockApi.SetResponse("friends.delete", tc.vkServerResponse)

			vk.DeleteFriend(999)

			require.Contains(t, captured, tc.expectedOutput)
		})
	}
}

func Test_RemoveDeletedFriends_ShouldCallDeleteOnlyForUsersWithDeletedFlag(t *testing.T) {
	deleteFriendActualCallCounter := 0
	var deletedFriendId int
	mock := requests.NewMockRequestSender().SetResponse("friends.get", deletedUserJson)
	vk := NewVkFromMock(mock)
	monkey.PatchInstanceMethod(reflect.TypeOf(vk), "DeleteFriend",
		func(vk *Vk, id int) {
			deleteFriendActualCallCounter++
			deletedFriendId = id
		})
	defer monkey.UnpatchAll()

	vk.RemoveDeletedFriends()

	require.Equal(t, 1, deleteFriendActualCallCounter)
	require.Equal(t, 2, deletedFriendId)
}
