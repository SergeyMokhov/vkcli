package friends

import (
	"fmt"
	"net/url"
)

const methodBase string = "friends."

type friendsGetRequest struct {
	values url.Values
}

func (fg friendsGetRequest) UrlValues() url.Values {
	return fg.values
}

func (fg friendsGetRequest) Method() string {
	return fmt.Sprint(methodBase, "get")
}

func Get() friendsGetRequest {
	return friendsGetRequest{values: url.Values{}}
}
