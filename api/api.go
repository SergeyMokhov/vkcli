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

func NewApiFromMock(mock *requests.MockRequestSender) *Api {
	return &Api{
		requestSender: mock,
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

/*userId - ID of the user whose friend request will be approved or to whom a friend request will be sent.

text - Text of the message (up to 500 characters) for the friend request, if any.

follow - true to pass an incoming request to followers list.

Returns one of the following values:
   1 — friend request sent;
   2 — friend request from the user approved;
   4 — request resending.*/
func (rd *Api) AddFriend(userId int, text string, follow bool) (response int, err error) {
	addAs := requests.AsFriend

	if follow {
		addAs = requests.AsFollower
	}

	request := requests.FriendsAdd(userId, text, addAs)
	err = rd.requestSender.SendVkRequestAndRetryOnCaptcha(request)

	if err != nil {
		return
	}

	resp, ok := request.ResponseStructPointer.(*requests.FriendsAddResponse)
	if !ok {
		return
	}

	if resp.Error != nil {
		err = fmt.Errorf("Vk.com returned an error: %v", resp.Error)
	}

	return resp.Response, err
}

func (rd *Api) DeleteFriend(userId int) (success, friendDeleted, inRequestDeleted, OutRequestDeleted, suggestionDeleted int, err error) {
	request := requests.FriendsDelete(userId)
	err = rd.requestSender.SendVkRequestAndRetryOnCaptcha(request)

	if err != nil {
		return
	}

	resp, ok := request.ResponseStructPointer.(*requests.FriendsDeleteResponse)
	if !ok {
		return
	}

	if resp.Error != nil {
		err = fmt.Errorf("Vk.com returned an error: %v", resp.Error)
	}
	//ToDo make this return enum of "FriendDeleted/InRequestRemoved/OutRequestRemoved
	return resp.Response.Success,
		resp.Response.FriendDeleted,
		resp.Response.InRequestDeleted,
		resp.Response.OutRequestDeleted,
		resp.Response.SuggestionDeleted,
		err
}
