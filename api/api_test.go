package api

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vk-cli/api/requests"
	"testing"
)

const fakeResponse = `{
  "response": {
    "count": 94,
    "items": [
      {
        "id": 988385,
        "first_name": "name1",
        "last_name": "lastName1",
        "is_closed": false,
        "can_access_closed": true,
        "nickname": "",
        "bdate": "20.2.1985",
        "online": 0
      },
      {
        "id": 2032170,
        "first_name": "name2",
        "last_name": "lastName2",
        "is_closed": false,
        "can_access_closed": true,
        "nickname": "★",
        "bdate": "1.11.1989",
        "online": 0
      }
    ]
  }
}
`

//TODO Add tests for all the API functions

func TestApi_GetAllFriends_ReturnsAllFriends(t *testing.T) {
	mock := requests.NewMockRequestSender()
	mock.SetResponse("friends.get", fakeResponse)
	api := Api{requestSender: mock}

	users, err := api.GetAllFriends()
	require.Nil(t, err)
	require.EqualValues(t, 2, len(users))
}
