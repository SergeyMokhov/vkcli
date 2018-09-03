package friends

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	methodBase             string = "friends."
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

type friendsGetRequest struct {
	values url.Values
}

type fields string

type order string

func (fg *friendsGetRequest) UrlValues() url.Values {
	return fg.values
}

func (fg *friendsGetRequest) Method() string {
	return fmt.Sprint(methodBase, "get")
}

func (fg *friendsGetRequest) SetUserId(usrId int) *friendsGetRequest {
	fg.values.Add("user_id", strconv.Itoa(usrId))
	return fg
}

func (fg *friendsGetRequest) SetOrder(order order) *friendsGetRequest {
	fg.values.Add("order", string(order))
	return fg
}

func (fg *friendsGetRequest) SetListId(positive int) *friendsGetRequest {
	fg.values.Add("list_id", strconv.Itoa(positive))
	return fg
}

func (fg *friendsGetRequest) SetCount(positive int) *friendsGetRequest {
	fg.values.Add("count", strconv.Itoa(positive))
	return fg
}

func (fg *friendsGetRequest) SetOffset(positive int) *friendsGetRequest {
	fg.values.Add("offset", strconv.Itoa(positive))
	return fg
}

func (fg *friendsGetRequest) SetFields(fields ...fields) *friendsGetRequest {
	fieldsStrings := make([]string, len(fields))
	for i, f := range fields {
		fieldsStrings[i] = string(f)
	}

	fg.values.Add("fields", strings.Join(fieldsStrings, ","))

	return fg
}

func Get() *friendsGetRequest {
	return &friendsGetRequest{values: url.Values{}}
}
