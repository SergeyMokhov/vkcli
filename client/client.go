package client

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/requests/friends"
	"golang.org/x/oauth2"
	"log"
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
	v, err := friends.Get().SetFields(friends.BDate, friends.HasMobile, friends.Nickname).Perform(vk.RequestSender())
	if err != nil {
		log.Fatalf("Failed to list friends:%v", err)
	}
	if v.Error != nil {
		fmt.Printf("Vk returned an error: %v", v.Error)
		return
	}
	fmt.Printf("Count:%v, Lenth:%v\n", v.Response.Count, len(v.Response.Items))
	for i, val := range v.Response.Items {
		fmt.Printf("%v.	ID: %v,	Name: %v %v %v,	Deactivated: %v	BDate: %v,	HasMobile: %v\n",
			i, val.Id, val.FirstName, val.Nickname, val.LastName, val.Deactivated, val.BDate, val.HasMobile)
	}
}

func (vk *Vk) AddFriend(id int) {
	resp, err := friends.Add(id, "lol", friends.AsFollower).Perform(vk.RequestSender())
	fmt.Printf("Response: %v, Error: %v, VkError: %v", resp.Response, err, resp.Error)
}

func (vk *Vk) DeleteFriend(id int) {
	//Todo implement client tests. Including this function.
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
