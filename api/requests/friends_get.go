package requests

import (
	"gitlab.com/g00g/vk-cli/api/obj"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
	"strings"
)

const (
	Hints                  Order  = "hints"
	Random                 Order  = "random"
	Mobile                 Order  = "mobile"
	Name                   Order  = "name"
	Nickname               Fields = "nickname"
	Domain                 Fields = "domain"
	Sex                    Fields = "sex"
	BDate                  Fields = "bdate"
	City                   Fields = "city"
	Country                Fields = "country"
	Timezone               Fields = "timezone"
	Photo50                Fields = "photo_50"
	Photo100               Fields = "photo_100"
	Photo200Orig           Fields = "photo_200_orig"
	HasMobile              Fields = "has_mobile"
	Contacts               Fields = "contacts"
	Education              Fields = "education"
	Online                 Fields = "online"
	Relation               Fields = "relation"
	LastSeen               Fields = "last_seen"
	Status                 Fields = "status"
	CanWritePrivateMessage Fields = "can_write_private_message"
	CanSeeAllPhotos        Fields = "can_see_all_posts"
	CanPost                Fields = "can_post"
	Universities           Fields = "universities"
)

type GetResponse struct {
	Response GetResponseValue `json:"response"`
	*vkErrors.Error
}

func (gr *GetResponse) GetError() *vkErrors.Error {
	return gr.Error
}

type GetResponseValue struct {
	Count int        `json:"count"`
	Items []obj.User `json:"items"`
}

type getRequest struct {
	*VkRequestBase
}

type Fields string

type Order string

func (fg *getRequest) SetUserId(usrId int) *getRequest {
	fg.Values.Add("user_id", strconv.Itoa(usrId))
	return fg
}

func (fg *getRequest) SetOrder(order Order) *getRequest {
	fg.Values.Add("order", string(order))
	return fg
}

func (fg *getRequest) SetListId(positive int) *getRequest {
	fg.Values.Add("list_id", strconv.Itoa(positive))
	return fg
}

func (fg *getRequest) SetCount(positive int) *getRequest {
	fg.Values.Add("count", strconv.Itoa(positive))
	return fg
}

func (fg *getRequest) SetOffset(positive int) *getRequest {
	fg.Values.Add("offset", strconv.Itoa(positive))
	return fg
}

func (fg *getRequest) SetFields(fields []Fields) *getRequest {
	fieldsStrings := make([]string, len(fields))
	for i, f := range fields {
		fieldsStrings[i] = string(f)
	}

	fg.Values.Add("fields", strings.Join(fieldsStrings, ","))

	return fg
}

func Get() (request *getRequest) {
	return &getRequest{
		VkRequestBase: NewVkRequestBase("friends.get", &GetResponse{}),
	}
}
