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

func TestFriendsGetRequest_SetListId(t *testing.T) {
	fgr := Get().SetListId(7)
	actual := fgr.UrlValues().Get("list_id")
	require.EqualValues(t, "7", actual)
}

func TestFriendsGetRequest_SetCount(t *testing.T) {
	fgr := Get().SetCount(10000)
	actual := fgr.UrlValues().Get("count")
	require.EqualValues(t, "10000", actual)
}

func TestFriendsGetRequest_SetOffset(t *testing.T) {
	fgr := Get().SetOffset(4)
	actual := fgr.UrlValues().Get("offset")
	require.EqualValues(t, "4", actual)
}

func TestFriendsGetRequest_SetFields(t *testing.T) {
	fgr := Get().SetFields(Nickname, Sex, Domain)
	actual := fgr.UrlValues().Get("fields")
	require.EqualValues(t, "nickname,sex,domain", actual)
}
