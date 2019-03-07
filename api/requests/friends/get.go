package friends

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api"
	"gitlab.com/g00g/vk-cli/api/obj"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"strconv"
	"strings"
)

const (
	Hints                  order  = "hints"
	Random                 order  = "random"
	Mobile                 order  = "mobile"
	Name                   order  = "name"
	Nickname               fields = "nickname"
	Domain                 fields = "domain"
	Sex                    fields = "sex"
	BDate                  fields = "bdate"
	City                   fields = "city"
	Country                fields = "country"
	Timezone               fields = "timezone"
	Photo50                fields = "photo_50"
	Photo100               fields = "photo_100"
	Photo200Orig           fields = "photo_200_orig"
	HasMobile              fields = "has_mobile"
	Contacts               fields = "contacts"
	Education              fields = "education"
	Online                 fields = "online"
	Relation               fields = "relation"
	LastSeen               fields = "last_seen"
	Status                 fields = "status"
	CanWritePrivateMessage fields = "can_write_private_message"
	CanSeeAllPhotos        fields = "can_see_all_posts"
	CanPost                fields = "can_post"
	Universities           fields = "universities"
)

type GetResponse struct {
	Response getResponseValue `json:"response"`
	*vkErrors.Error
}

func (gr *GetResponse) GetError() *vkErrors.Error {
	return gr.Error
}

type getResponseValue struct {
	Count int        `json:"count"`
	Items []obj.User `json:"items"`
}

type getRequest struct {
	*api.VkRequestBase
}

type fields string

type order string

// Returns error only if sending request or type conversion fails
func (fg *getRequest) Perform(api *api.Api) (response *GetResponse, err error) {
	err = api.SendVkRequestAndRetryOnCaptcha(fg)

	resp, ok := fg.ResponseStructPointer.(*GetResponse)
	if ok {
		return resp, err
	}

	return &GetResponse{}, err
}

func (fg *getRequest) SetUserId(usrId int) *getRequest {
	fg.Values.Add("user_id", strconv.Itoa(usrId))
	return fg
}

func (fg *getRequest) SetOrder(order order) *getRequest {
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

func (fg *getRequest) SetFields(fields ...fields) *getRequest {
	fieldsStrings := make([]string, len(fields))
	for i, f := range fields {
		fieldsStrings[i] = string(f)
	}

	fg.Values.Add("fields", strings.Join(fieldsStrings, ","))

	return fg
}

func Get() *getRequest {
	return &getRequest{
		VkRequestBase: api.NewVkRequestBase(fmt.Sprint(methodBase, "get"), &GetResponse{}),
	}
}
