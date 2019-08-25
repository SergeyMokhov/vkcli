package requests

import (
	"gitlab.com/g00g/vk-cli/api/obj"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
	"strings"
)

const (
	Hints  FriendsGetOrder = "hints"
	Random FriendsGetOrder = "random"
	Mobile FriendsGetOrder = "mobile"
	Name   FriendsGetOrder = "name"

	Nickname               FriendsGetFields = "nickname"
	Domain                 FriendsGetFields = "domain"
	Sex                    FriendsGetFields = "sex"
	BDate                  FriendsGetFields = "bdate"
	City                   FriendsGetFields = "city"
	Country                FriendsGetFields = "country"
	Timezone               FriendsGetFields = "timezone"
	Photo50                FriendsGetFields = "photo_50"
	Photo100               FriendsGetFields = "photo_100"
	Photo200Orig           FriendsGetFields = "photo_200_orig"
	HasMobile              FriendsGetFields = "has_mobile"
	Contacts               FriendsGetFields = "contacts"
	Education              FriendsGetFields = "education"
	Online                 FriendsGetFields = "online"
	Relation               FriendsGetFields = "relation"
	LastSeen               FriendsGetFields = "last_seen"
	Status                 FriendsGetFields = "status"
	CanWritePrivateMessage FriendsGetFields = "can_write_private_message"
	CanSeeAllPhotos        FriendsGetFields = "can_see_all_posts"
	CanPost                FriendsGetFields = "can_post"
	Universities           FriendsGetFields = "universities"
)

type FriendsGetResponse struct {
	Response FriendsGetResponseValue `json:"response"`
	*vkErrors.Error
}

func (gr *FriendsGetResponse) GetError() *vkErrors.Error {
	return gr.Error
}

type FriendsGetResponseValue struct {
	Count int        `json:"count"`
	Items []obj.User `json:"items"`
}

type friendsGetRequest struct {
	*VkRequestBase
}

type FriendsGetFields string

type FriendsGetOrder string

func (fg *friendsGetRequest) SetUserId(usrId int) *friendsGetRequest {
	fg.Values.Add("user_id", strconv.Itoa(usrId))
	return fg
}

func (fg *friendsGetRequest) SetOrder(order FriendsGetOrder) *friendsGetRequest {
	fg.Values.Add("order", string(order))
	return fg
}

func (fg *friendsGetRequest) SetListId(positive int) *friendsGetRequest {
	fg.Values.Add("list_id", strconv.Itoa(positive))
	return fg
}

func (fg *friendsGetRequest) SetCount(positive int) *friendsGetRequest {
	fg.Values.Add("count", strconv.Itoa(positive))
	return fg
}

func (fg *friendsGetRequest) SetOffset(positive int) *friendsGetRequest {
	fg.Values.Add("offset", strconv.Itoa(positive))
	return fg
}

func (fg *friendsGetRequest) SetFields(fields []FriendsGetFields) *friendsGetRequest {
	fieldsStrings := make([]string, len(fields))
	for i, f := range fields {
		fieldsStrings[i] = string(f)
	}

	fg.Values.Add("fields", strings.Join(fieldsStrings, ","))

	return fg
}

func FriendsGet() (request *friendsGetRequest) {
	return &friendsGetRequest{
		VkRequestBase: NewVkRequestBase("friends.get", &FriendsGetResponse{}),
	}
}
