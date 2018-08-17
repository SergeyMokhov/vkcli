package api

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"testing"
)

func TestNewVk(t *testing.T) {
	vk := NewVk(&oauth2.Token{})

	require.NotNil(t, vk.token)
	require.NotNil(t, vk.client)
	require.NotNil(t, vk.baseUrl)
}

func TestNewRequestShouldReturnUrlAndParams(t *testing.T) {
	token := &oauth2.Token{AccessToken: "123"}
	vk := NewVk(token)

	methodUrl, params, err := vk.newRequest("friends.get")

	require.Nil(t, err)
	require.EqualValues(t, "https://api.vk.com/method/friends.get", methodUrl.String())
	require.EqualValues(t, "123", params.Get("access_token"))
	require.EqualValues(t, "1", params.Get("https"))
	require.NotEqual(t, "", params.Get("v"))
}
