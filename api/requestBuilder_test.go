package api

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vkcli/api/requests/friends"
	"golang.org/x/oauth2"
	"testing"
)

func TestNewRequestShouldReturnUrlAndParams(t *testing.T) {
	token := &oauth2.Token{AccessToken: "123"}
	vk := NewVk(token)

	methodUrl, params, err := NewRequestBuilder(vk).NewRequest(friends.Get())

	require.Nil(t, err)
	require.EqualValues(t, "https://api.vk.com/method/friends.get", methodUrl.String())
	require.EqualValues(t, "123", params.Get("access_token"))
	require.EqualValues(t, "1", params.Get("https"))
	require.NotEqual(t, "", params.Get("v"))
}
