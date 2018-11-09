package friends

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
)

const (
	AsFollower followerFlag = "1"
	AsFriend   followerFlag = "0"
)

type followerFlag string

type friendsAddRequest struct {
	*api.DummyVkRequest
}

type FriendsAddResponse struct {
	Response int `json:"response"`
	vkErrors.Error
}

func Add(userId int, text string, follow followerFlag) *friendsAddRequest {
	req := &friendsAddRequest{
		DummyVkRequest: api.NewDummyVkRequest(fmt.Sprint(methodBase, "add"), &FriendsAddResponse{})}

	req.Values.Add("user_id", strconv.Itoa(userId))
	req.Values.Add("text", text)
	req.Values.Add("follow", string(follow))

	return req
}

func (fa *friendsAddRequest) Perform(api *api.Api) (response *FriendsAddResponse, err error) {
	err = api.SendRequestAndRetyOnCaptcha(fa)

	resp, ok := fa.ResponseStructPointer.(*FriendsAddResponse)
	if ok {
		return resp, err
	}

	return &FriendsAddResponse{}, err
}
