package friends

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/g00g/vkcli/api"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const friendsAddresponse = `{
  "response": 1
}`

func TestAdd(t *testing.T) {
	req := Add(123, "add me please", AsFollower)
	require.EqualValues(t, "123", req.Values.Get("user_id"))
	require.EqualValues(t, "add me please", req.Values.Get("text"))
	require.EqualValues(t, "1", req.Values.Get("follow"))
	require.EqualValues(t, "friends.add", req.MethodStr)
}

func TestFriendsAddRequest_Perform(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, friendsAddresponse)
	}))
	defer ts.Close()

	requestSender := api.NewInstance(&oauth2.Token{AccessToken: "000"})
	baseUrl, urlParseErr := url.Parse(ts.URL)
	require.Nil(t, urlParseErr)
	requestSender.BaseUrl = baseUrl

	response, err := Add(505, "tst", AsFriend).Perform(requestSender)
	require.Nil(t, err)
	assert.EqualValues(t, 0, response.Error.ErrorCode)
	assert.EqualValues(t, 1, response.Response)
}
