package requests

import (
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
)

var SuccessFriendDeleted = `{"response":{"success":1,"friend_deleted":1}}`
var SuccessInRequestDeclined = `{"response":{"success":1,"in_request_deleted":1}}`
var SuccessOutRequestCancelled = `{"response":{"success":1,"out_request_deleted":1}}`
var FriedDeleteFailureAccessDenied = `{"error":{"error_code":15,"error_msg":"Access denied: No friend or friend request found.","request_params":[{"key":"oauth","value":"1"},{"key":"method","value":"friends.delete"},{"key":"https","value":"1"},{"key":"user_id","value":"2916112"},{"key":"v","value":"5.92"}]}}`

type friendsDeleteRequest struct {
	*VkRequestBase
}

type FriendsDeleteResponse struct {
	Response friendsDeleteResponse `json:"response"`
	*vkErrors.Error
}

func (dr *FriendsDeleteResponse) GetError() *vkErrors.Error {
	return dr.Error
}

type friendsDeleteResponse struct {
	Success           int `json:"success"`
	FriendDeleted     int `json:"friend_deleted"`
	OutRequestDeleted int `json:"out_request_deleted"`
	InRequestDeleted  int `json:"in_request_deleted"`
	SuggestionDeleted int `json:"suggestion_deleted"`
}

func FriendsDelete(userId int) *friendsDeleteRequest {
	rd := &friendsDeleteRequest{
		VkRequestBase: NewVkRequestBase("friends.delete", &FriendsDeleteResponse{}),
	}
	rd.Values.Add("user_id", strconv.Itoa(userId))
	return rd
}
