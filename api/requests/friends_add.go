package requests

import (
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
)

const (
	AsFollower friendsAddFollowerFlag = "1"
	AsFriend   friendsAddFollowerFlag = "0"
)

type friendsAddFollowerFlag string

type friendsAddRequest struct {
	*VkRequestBase
}

type FriendsAddResponse struct {
	Response int `json:"response"`
	*vkErrors.Error
}

func (ar *FriendsAddResponse) GetError() *vkErrors.Error {
	return ar.Error
}

// userId - ID of the user whose friend request will be approved or to whom a friend request will be sent.
// text	- Text of the message (up to 500 characters) for the friend request, if any
// follow - 1 to pass an incoming request to followers list.
func FriendsAdd(userId int, text string, follow friendsAddFollowerFlag) *friendsAddRequest {
	req := &friendsAddRequest{
		VkRequestBase: NewVkRequestBase("friends.add", &FriendsAddResponse{})}

	req.Values.Add("user_id", strconv.Itoa(userId))
	req.Values.Add("text", text)
	req.Values.Add("follow", string(follow))

	return req
}
