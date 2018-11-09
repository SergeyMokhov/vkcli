package client

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/requests/friends"
	"golang.org/x/oauth2"
	"log"
)

type Vk struct {
	api *api.Api
}

func NewVk(token *oauth2.Token) *Vk {
	return &Vk{api: api.NewInstance(token)}
}

func (vk *Vk) ListFriends() {
	v, err := friends.Get().SetFields(friends.BDate, friends.HasMobile, friends.Nickname).Perform(vk.api)
	if err != nil {
		log.Fatalf("Failed to list friends:%v", err)
	}
	if v.Error.ErrorCode != 0 {
		fmt.Printf("Vk returned an error: %v", v.Error)
		return
	}
	fmt.Printf("Count:%v, Lenth:%v\n", v.Value.Count, len(v.Value.Items))
	for i, val := range v.Value.Items {
		fmt.Printf("%v.	ID: %v,	Name: %v %v %v,	Deactivated: %v	BDate: %v,	HasMobile: %v\n",
			i, val.Id, val.FirstName, val.Nickname, val.LastName, val.Deactivated, val.BDate, val.HasMobile)
	}
}

func (vk *Vk) AddFriend(id int) {
	resp, err := friends.Add(id, "lol", friends.AsFollower).Perform(vk.api)
	fmt.Printf("Response: %v, Error: %v, VkError: %v", resp.Response, err, resp.Error)
}
