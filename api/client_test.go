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
