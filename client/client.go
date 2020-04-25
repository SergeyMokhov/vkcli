package client

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/obj"
	"gitlab.com/g00g/vk-cli/api/requests"
	"golang.org/x/oauth2"
	"strconv"
	"time"
)

type Vk struct {
	api *api.Api
}

func NewVk(token *oauth2.Token) *Vk {
	return &Vk{api: api.NewApi(token)}
}

func NewVkFromMock(mock *requests.MockRequestSender) *Vk {
	return &Vk{api: api.NewApiFromMock(mock)}
}

func (vk *Vk) ListFriends() {
	format := "%6s%20s%20s%20s%20s%13s%7s%12s%30s\n"

	allFriends, err := vk.api.GetAllFriends(requests.Online, requests.BDate, requests.Nickname, requests.LastSeen)
	if err != nil {
		fmt.Printf("Failed to list friends:%v", err)
		return
	}

	fmt.Printf(format, "#", "ID", "First name", "Last name", "Nickname", "Birthday", "Online", "Deactivated", "Last Seen")
	for i, val := range allFriends {
		fmt.Printf(format, strconv.Itoa(i), strconv.Itoa(val.Id), val.FirstName, val.LastName, val.Nickname, val.BDate,
			strconv.Itoa(val.Online), val.Deactivated, time.Unix(val.LastSeen.Time, 0))
	}
}

func (vk *Vk) AddFriend(id int) {
	errorAddingFriendFormat := "Error adding friend with Id %v: %v\n"

	response, err := vk.api.AddFriend(id, "it is me, VK bot.", false)
	if err != nil {
		fmt.Printf(errorAddingFriendFormat, id, err)
	}

	var successStringFormat string
	switch {
	case response == 1:
		successStringFormat = "Friend request sent to %v"
	case response == 2:
		successStringFormat = "Friend request from %v has been approved"
	case response == 4:
		successStringFormat = "Resending request to %v"
	default:
		successStringFormat = "Unknown response code: " + strconv.Itoa(response) + " for friends.add request for userid: %v"
	}

	fmt.Printf(successStringFormat+"\n", id)
}

func (vk *Vk) DeleteFriend(userId int) {
	success, fr, in, out, _, err := vk.api.DeleteFriend(userId)

	//Vk uses same api for deleting from friend list and declining friend requests,create method body
	// etc. So, what did just happened?
	var successStringFormat string
	switch {
	case fr == 1:
		successStringFormat = "%v successfully deleted form friend list"
	case in == 1:
		successStringFormat = "Successfully declined friend request from %v"
	case out == 1:
		successStringFormat = "Successfully cancelled friend request to %v"
	default:
		successStringFormat = "Successfully deleted suggestion of %v"
	}

	switch success {
	case 1:
		fmt.Printf(successStringFormat+"\n", userId)
	default:
		fmt.Printf("Error removing/declining request for %v: %v\n", userId, err.Error())
	}
}

func (vk *Vk) RemoveDeletedFriends() {
	allFriends, err := vk.api.GetAllFriends(requests.LastSeen)
	if err != nil {
		fmt.Printf("Failed to list friends:%v", err)
		return
	}

	for _, friend := range allFriends {
		if friend.Deactivated == obj.UserDeleted {
			fmt.Printf("Deleting friend %d, that goes by name '%s %s'\n", friend.Id, friend.FirstName, friend.LastName)
			vk.DeleteFriend(friend.Id)
		}
	}
}
