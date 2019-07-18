package requests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var fakeRequestMethod = "testMethod"

func TestNewMockRequestSender_ShouldUseKeepLatestRequest(t *testing.T) {
	expectedResponse := `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found."}}`
	mock := NewMockRequestSender().SetResponse(fakeRequestMethod, expectedResponse)
	defer mock.Shutdown()
	req := FakeVkRequest{VkRequestBase: NewVkRequestBase(fakeRequestMethod, &FakeVkResponse{})}
	response, err := sendRequest(mock.VkRequestSender, &req)

	require.Nil(t, err)
	require.Equal(t, "/testMethod", mock.LastRequest.RequestURI)
	require.Equal(t, expectedResponse, string(response))
}

func TestMockApi_SetResponseCanBeOverriden(t *testing.T) {
	expectedResponse := `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found."}}`

	mock := NewMockRequestSender().SetResponse(fakeRequestMethod, "Some Random String")
	defer mock.Shutdown()
	req := FakeVkRequest{VkRequestBase: NewVkRequestBase(fakeRequestMethod, &FakeVkResponse{})}

	mock.SetResponse(fakeRequestMethod, expectedResponse)
	response, err := sendRequest(mock.VkRequestSender, &req)

	require.Nil(t, err)
	require.Equal(t, expectedResponse, string(response))
}

func TestMockApi_SupportsDifferentResponsesForDifferentMethods(t *testing.T) {
	expectedResponseA := `A`
	methodA := "ma"
	expectedResponseB := `B`
	methodB := "mb"

	mock := NewMockRequestSender().
		SetResponse(methodA, expectedResponseA).
		SetResponse(methodB, expectedResponseB)
	defer mock.Shutdown()
	reqA := FakeVkRequest{VkRequestBase: NewVkRequestBase(methodA, &FakeVkResponse{})}
	reqB := FakeVkRequest{VkRequestBase: NewVkRequestBase(methodB, &FakeVkResponse{})}

	responseA, _ := sendRequest(mock.VkRequestSender, &reqA)
	responseB, _ := sendRequest(mock.VkRequestSender, &reqB)

	require.Equal(t, expectedResponseA, string(responseA))
	require.Equal(t, expectedResponseB, string(responseB))
}

func TestMockApi_NumberOfRequestsReceivedIncreasesIndependentluyForEachMethod(t *testing.T) {
	methodA := "ma"
	methodB := "mb"

	mock := NewMockRequestSender().
		SetResponse(methodA, "").
		SetResponse(methodB, "")
	defer mock.Shutdown()
	reqA := FakeVkRequest{VkRequestBase: NewVkRequestBase(methodA, &FakeVkResponse{})}
	reqB := FakeVkRequest{VkRequestBase: NewVkRequestBase(methodB, &FakeVkResponse{})}

	sendRequest(mock.VkRequestSender, &reqB)
	sendRequest(mock.VkRequestSender, &reqA)
	sendRequest(mock.VkRequestSender, &reqB)
	sendRequest(mock.VkRequestSender, &reqB)

	assert.Equal(t, 0, mock.NumberOfRequestsReceived("zero"))
	assert.Equal(t, 1, mock.NumberOfRequestsReceived(methodA))
	assert.Equal(t, 3, mock.NumberOfRequestsReceived(methodB))
}
