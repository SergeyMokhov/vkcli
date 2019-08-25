package client

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/requests"
	"golang.org/x/oauth2"
	"strconv"
)

//type VkClient interface {
//	RequestSender() api.RequestSendRetryer
//}

type Vk struct {
	api *api.Api
}

//func (vk *Vk) RequestSender() api.RequestSendRetryer {
//	return vk.api
//}

func NewVk(token *oauth2.Token) *Vk {
	return &Vk{api: api.NewApi(token)}
}

//func newVkFromMockApi(mock *requests.MockRequestSender) *Vk {
//	return &Vk{api: mock.Api}
//}

func (vk *Vk) ListFriends() {
	format := "%6s%20s%20s%20s%20s%13s%7s%12s\n"

	allFriends, err := vk.api.GetAllFriends(requests.Online, requests.BDate, requests.Nickname)
	if err != nil {
		fmt.Printf("Failed to list friends:%v", err)
		return
	}

	fmt.Printf(format, "#", "ID", "First name", "Last name", "Nickname", "Birthday", "Online", "Deactivated")
	for i, val := range allFriends {
		fmt.Printf(format, strconv.Itoa(i), strconv.Itoa(val.Id), val.FirstName, val.LastName, val.Nickname, val.BDate,
			strconv.Itoa(val.Online), val.Deactivated)
	}
}

func (vk *Vk) AddFriend(id int) {
	//errorAddingFriendFormat := "Error adding friend with Id %v: %v\n"
	//
	//resp, err := friends.Add(id, "lol", friends.AsFollower).Perform(vk.RequestSender())
	//if err != nil {
	//	fmt.Printf(errorAddingFriendFormat, id, err)
	//}
	//
	//var successStringFormat string
	//switch {
	//case resp.Response == 1:
	//	successStringFormat = "Friend request sent to %v"
	//case resp.Response == 2:
	//	successStringFormat = "Friend request from %v has been approved"
	//case resp.Response == 4:
	//	successStringFormat = "Request resending"
	//default:
	//	successStringFormat = "Unknown response code: " + strconv.Itoa(resp.Response)
	//}
	//
	//switch resp.Error {
	//case nil:
	//	fmt.Printf(successStringFormat+"\n", id)
	//default:
	//	fmt.Printf(errorAddingFriendFormat, id, resp.Error)
	//}
}

func (vk *Vk) DeleteFriend(id int) {
	//resp, err := friends.Delete(id).Perform(vk.RequestSender())
	//if err != nil {
	//	fmt.Printf("Error deleting %v from friend list: %v\n", id, err)
	//}
	//
	////Vk uses same api for deleting from friend list and declining friend requests,create method body
	//// etc. So, what did just happened?
	//var successStringFormat string
	//switch {
	//case resp.Response.FriendDeleted == 1:
	//	successStringFormat = "%v successfully deleted form friend list"
	//case resp.Response.InRequestDeleted == 1:
	//	successStringFormat = "Successfully declined friend request from %v"
	//case resp.Response.OutRequestDeleted == 1:
	//	successStringFormat = "Successfully cancelled friend request to %v"
	//default:
	//	successStringFormat = "Successfully deleted suggestion of %v"
	//}
	//
	//switch resp.Response.Success {
	//case 1:
	//	fmt.Printf(successStringFormat+"\n", id)
	//default:
	//	fmt.Printf("Error removing/declining request for %v: %v\n", id, resp.GetError())
	//}
}
