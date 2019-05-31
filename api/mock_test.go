package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var fakeRequestMethod = "testMethod"

func TestNewMockApi_ShouldUseKeepLatestRequest(t *testing.T) {
	expectedResponse := `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found."}}`
	mock := NewMockApi().SetResponse(fakeRequestMethod, expectedResponse)
	defer mock.Shutdown()
	req := fakeVkRequest{VkRequestBase: NewVkRequestBase(fakeRequestMethod, &fakeVkResponse{})}
	response, err := sendRequest(mock.Api, &req)

	require.Nil(t, err)
	require.Equal(t, "/testMethod", mock.LastRequest.RequestURI)
	require.Equal(t, expectedResponse, string(response))
}

func TestMockApi_SetResponseCanBeOverriden(t *testing.T) {
	expectedResponse := `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found."}}`

	mock := NewMockApi().SetResponse(fakeRequestMethod, "Some Random String")
	defer mock.Shutdown()
	req := fakeVkRequest{VkRequestBase: NewVkRequestBase(fakeRequestMethod, &fakeVkResponse{})}

	mock.SetResponse(fakeRequestMethod, expectedResponse)
	response, err := sendRequest(mock.Api, &req)

	require.Nil(t, err)
	require.Equal(t, expectedResponse, string(response))
}

func TestMockApi_SupportsDifferentResponsesForDifferentMethods(t *testing.T) {
	expectedResponseA := `A`
	methodA := "ma"
	expectedResponseB := `B`
	methodB := "mb"

	mock := NewMockApi().
		SetResponse(methodA, expectedResponseA).
		SetResponse(methodB, expectedResponseB)
	defer mock.Shutdown()
	reqA := fakeVkRequest{VkRequestBase: NewVkRequestBase(methodA, &fakeVkResponse{})}
	reqB := fakeVkRequest{VkRequestBase: NewVkRequestBase(methodB, &fakeVkResponse{})}

	responseA, _ := sendRequest(mock.Api, &reqA)
	responseB, _ := sendRequest(mock.Api, &reqB)

	require.Equal(t, expectedResponseA, string(responseA))
	require.Equal(t, expectedResponseB, string(responseB))
}
