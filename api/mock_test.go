package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMockApi_ShouldUseKeepLatestRequest(t *testing.T) {
	expectedResponse := `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found."}}`
	mock := NewMockApi(expectedResponse)
	defer mock.Shutdown()
	req := fakeVkRequest{VkRequestBase: NewVkRequestBase("testMethod", &fakeVkResponse{})}
	response, err := sendRequest(mock.Api, &req)

	require.Nil(t, err)
	require.Equal(t, "/testMethod", mock.LastRequest.RequestURI)
	require.Equal(t, expectedResponse, string(response))
}

func TestMockApi_SetResponse(t *testing.T) {
	expectedResponse := `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found."}}`

	mock := NewMockApi("Some Random String")
	defer mock.Shutdown()
	req := fakeVkRequest{VkRequestBase: NewVkRequestBase("testMethod", &fakeVkResponse{})}

	mock.SetResponse(expectedResponse)
	response, err := sendRequest(mock.Api, &req)

	require.Nil(t, err)
	require.Equal(t, expectedResponse, string(response))
}
