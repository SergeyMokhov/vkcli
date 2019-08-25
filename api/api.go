package api

import (
	"fmt"
	"gitlab.com/g00g/vk-cli/api/obj"
	"gitlab.com/g00g/vk-cli/api/requests"
	"golang.org/x/oauth2"
)

//TODO Add notes on how to use. Mention that you supposed to use single instance, otherwise you'll get banned by VK for
// sending too many requests. And  that you can use more than one API if they use different tokens.
type Api struct {
	requestSender requests.RequestSendRetryer
}

func NewApi(token *oauth2.Token) *Api {
	return &Api{
		requestSender: requests.NewVkRequestSender(token),
	}
}

func (rd *Api) GetAllFriends(userFields ...requests.FriendsGetFields) (users []obj.User, err error) {
	request := requests.FriendsGet().SetFields(userFields)

	err = rd.requestSender.SendVkRequestAndRetryOnCaptcha(request)

	if err != nil {
		return
	}

	response, ok := request.ResponseStructPointer.(*requests.FriendsGetResponse)
	if !ok {
		return
	}

	if response.Error != nil {
		err = fmt.Errorf("Vk.com returned an error: %v", response.Error)
	}

	return response.Response.Items, nil
}
