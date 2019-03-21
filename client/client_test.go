package client

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/requests/friends"
	"golang.org/x/oauth2"
	"testing"
)

func TestNewVk(t *testing.T) {
	vk := NewVk(&oauth2.Token{})

	require.NotNil(t, vk.api)
}

func TestVk_DeleteFriend_ShouldPrintLineOnSuccessfullDeletion(t *testing.T) {
	testCases := []struct {
		name             string
		vkServerResponse string
		expectedOutput   string
	}{
		{"Friend deleted from list", friends.SuccessFriendDeleted, "successfully deleted form friend list"},
		{"Declined incoming request", friends.SuccessInRequestDeclined, "successfully declined friend request"},
		{"Cancelled outgoing request", friends.SuccessOutRequestCancelled, "successfully cancelled friend request"},
		//{"Deleted friend suggestion", friends.} Was not able to find response example
		{"Access Denied", friends.FriedDeleteFailureAccessDenied, "Error removing/declining request"},
	}

	var captured string
	monkey.Patch(fmt.Printf,
		func(format string, a ...interface{}) (n int, err error) {
			captured = format
			return 0, nil
		})
	defer monkey.UnpatchAll()

	mockApi := api.NewMockApi("")
	vk := newVkFromMockApi(mockApi)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			captured = ""
			mockApi.SetResponse(tc.vkServerResponse)

			vk.DeleteFriend(999)

			require.Contains(t, captured, tc.expectedOutput)
		})
	}
}
