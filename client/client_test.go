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
	var captured string
	monkey.Patch(fmt.Printf,
		func(format string, a ...interface{}) (n int, err error) {
			captured = format
			return 0, nil
		})
	defer monkey.UnpatchAll()

	vk := newVkFromMockApi(api.NewMockApi(friends.SuccessFriendDeleted))

	vk.DeleteFriend(999)

	require.Contains(t, captured, "successfully deleted form friend list")
}
