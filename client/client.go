package client

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/requests/friends"
	"golang.org/x/oauth2"
	"strconv"
)

type VkClient interface {
	RequestSender() api.VkRequestSender
}

type Vk struct {
	api api.VkRequestSender
}

func (vk *Vk) RequestSender() api.VkRequestSender {
	return vk.api
}

func NewVk(token *oauth2.Token) *Vk {
	return &Vk{api: api.NewInstance(token)}
}

func newVkFromMockApi(mock *api.MockApi) *Vk {
	return &Vk{api: mock.Api}
}

func (vk *Vk) ListFriends() {
	v, err := friends.Get().SetFields(friends.Online, friends.BDate, friends.Nickname).Perform(vk.RequestSender())
	if err != nil {
		fmt.Printf("Failed to list friends:%v", err)
		return
	}
	if v.Error != nil {
		fmt.Printf("Vk returned an error: %v", v.Error)
		return
	}

	format := "%6s%20s%20s%20s%20s%13s%6s\n"
	fmt.Printf(format, "#", "ID", "First name", "Last name", "Nickname", "Birthday", "Online")
	for i, val := range v.Response.Items {
		fmt.Printf(format, strconv.Itoa(i), strconv.Itoa(val.Id), val.FirstName, val.LastName, val.Nickname, val.BDate, strconv.Itoa(val.Online))
	}
}

func (vk *Vk) AddFriend(id int) {
	resp, err := friends.Add(id, "lol", friends.AsFollower).Perform(vk.RequestSender())
	if err != nil {
		fmt.Printf("Error Adding friend with Id: %v. %v", id, err)
	}

	var successStringFormat string
	switch {
	case resp.Response == 1:
		successStringFormat = "friend request sent to %v"
	case resp.Response == 2:
		successStringFormat = "friend request from %v approved"
	case resp.Response == 4:
		successStringFormat = "request resending"
	default:
		successStringFormat = "unknown response code: " + strconv.Itoa(resp.Response)
	}

	switch resp.ErrorCode {
	case 0:
		fmt.Printf(successStringFormat, id)
	default:
		fmt.Printf("Error adding friend %v: %v", id, resp.Error)
	}
	fmt.Printf("Response: %v, Error: %v, VkError: %v", resp.Response, err, resp.Error)
}

func (vk *Vk) DeleteFriend(id int) {
	resp, err := friends.Delete(id).Perform(vk.RequestSender())
	if err != nil {
		fmt.Printf("Error deleting %v from friend list: %v\n", id, err)
	}

	//Vk uses same api for deleting from friend list and declining friend requests,create method body
	// etc. So, what did just happened?
	var successStringFormat string
	switch {
	case resp.Response.FriendDeleted == 1:
		successStringFormat = "%v successfully deleted form friend list"
	case resp.Response.InRequestDeleted == 1:
		successStringFormat = "successfully declined friend request from %v"
	case resp.Response.OutRequestDeleted == 1:
		successStringFormat = "successfully cancelled friend request to %v"
	default:
		successStringFormat = "successfully deleted suggestion of %v"
	}

	switch resp.Response.Success {
	case 1:
		fmt.Printf(successStringFormat+"\n", id)
	default:
		fmt.Printf("Error removing/declining request for %v: %v\n", id, resp.GetError())
	}
}
