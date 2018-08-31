package friends

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	methodBase string         = "friends."
	Hints      OrderOfFriends = "hints"
	Random     OrderOfFriends = "random"
	Mobile     OrderOfFriends = "mobile"
	Name       OrderOfFriends = "name"
)

type friendsGetRequest struct {
	values url.Values
}

type OrderOfFriends string

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

func (fg *friendsGetRequest) SetOrder(order OrderOfFriends) *friendsGetRequest {
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

func Get() *friendsGetRequest {
	return &friendsGetRequest{values: url.Values{}}
}
