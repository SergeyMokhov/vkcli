package friends

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFriendsGetRequest_Method(t *testing.T) {
	fgr := Get()
	actual := fgr.Method()
	require.EqualValues(t, "friends.get", actual)
}

func TestFriendsGetRequest_UrlValues(t *testing.T) {
	fgr := Get()
	actual := fgr.UrlValues()
	require.NotNil(t, actual)
}

func TestFriendsGetRequest_SetUserId(t *testing.T) {
	fgr := Get().SetUserId(9000)
	actual := fgr.UrlValues().Get("user_id")
	require.EqualValues(t, "9000", actual)
}

func TestFriendsGetRequest_SetOrder(t *testing.T) {
	fgr := Get().SetOrder(Name)
	actual := fgr.UrlValues().Get("order")
	require.EqualValues(t, "name", actual)
}
