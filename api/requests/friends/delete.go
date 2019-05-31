package friends

import (
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
)

var SuccessFriendDeleted = `{"response":{"success":1,"friend_deleted":1}}`
var SuccessInRequestDeclined = `{"response":{"success":1,"in_request_deleted":1}}`
var SuccessOutRequestCancelled = `{"response":{"success":1,"out_request_deleted":1}}`
var FriedDeleteFailureAccessDenied = `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found.","request_params":[{"key":"oauth","value":"1"},{"key":"method","value":"friends.delete"},{"key":"https","value":"1"},{"key":"user_id","value":"2916112"},{"key":"v","value":"5.92"}]}}`

type deleteRequest struct {
	*api.VkRequestBase
}

type DeleteResponse struct {
	Response deleteResponse `json:"response"`
	*vkErrors.Error
}

func (dr *DeleteResponse) GetError() *vkErrors.Error {
	return dr.Error
}

type deleteResponse struct {
	Success           int `json:"success"`
	FriendDeleted     int `json:"friend_deleted"`
	OutRequestDeleted int `json:"out_request_deleted"`
	InRequestDeleted  int `json:"in_request_deleted"`
	SuggestionDeleted int `json:"suggestion_deleted"`
}

func Delete(id int) *deleteRequest {
	rd := &deleteRequest{
		VkRequestBase: api.NewVkRequestBase(methodBase+"delete", &DeleteResponse{}),
	}

	rd.setUserId(id)
	return rd
}

// Returns error only if sending request or type conversion fails
func (dr *deleteRequest) Perform(api api.RequestSendRetrier) (response *DeleteResponse, err error) {
	err = api.SendVkRequestAndRetryOnCaptcha(dr)

	resp, ok := dr.ResponseStructPointer.(*DeleteResponse)
	if ok {
		return resp, err
	}

	return &DeleteResponse{}, err
}

func (dr *deleteRequest) setUserId(usrId int) *deleteRequest {
	dr.Values.Add("user_id", strconv.Itoa(usrId))
	return dr
}
