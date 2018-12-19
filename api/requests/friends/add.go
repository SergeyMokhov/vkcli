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

type addRequest struct {
	*api.VkRequestBase
}

type AddResponse struct {
	Response int `json:"response"`
	vkErrors.Error
}

func Add(userId int, text string, follow followerFlag) *addRequest {
	req := &addRequest{
		VkRequestBase: api.NewVkRequestBase(fmt.Sprint(methodBase, "add"), &AddResponse{})}

	req.Values.Add("user_id", strconv.Itoa(userId))
	req.Values.Add("text", text)
	req.Values.Add("follow", string(follow))

	return req
}

func (fa *addRequest) Perform(api *api.Api) (response *AddResponse, err error) {
	err = api.SendVkRequestAndRetryOnCaptcha(fa)

	resp, ok := fa.ResponseStructPointer.(*AddResponse)
	if ok {
		return resp, err
	}

	return &AddResponse{}, err
}
