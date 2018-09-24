package friends

import (
	"bufio"
	"fmt"
	"gitlab.com/g00g/vkcli/api"
	"gitlab.com/g00g/vkcli/api/obj"
	"os"
	"strconv"
	"strings"
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
	Response int             `json:"response"`
	Error    obj.VkErrorInfo `json:"error"`
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
	err = api.SendRequest(fa)

	resp, ok := fa.ResponseStructPointer.(*FriendsAddResponse)
	if ok {
		return resp, err
	}

	return &FriendsAddResponse{}, err
}

func (fa *friendsAddRequest) Do(api *api.Api) (response *FriendsAddResponse, err error) {
	try, err := fa.Perform(api)
	if try.Error.ErrorCode == obj.CaptchaRequired {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Please, solve the capture: %v\nCapture unswer is: ", try.Error.CaptchaImg)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimRight(answer, "\n")
		api.AddSolvedCapture(fa, try.Error, answer)
		return fa.Perform(api)
	} else {
		return try, err
	}
} //TODO  add tests for capture handling and retry.
