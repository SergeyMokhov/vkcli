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

func TestFriendsGetRequest_Method(t *testing.T) {
	fgr := Get()
	actual := fgr.Method()
	require.EqualValues(t, "friends.get", actual)
}

func TestFriendsGetRequest_UrlValues(t *testing.T) {
	fgr := Get()
	actual := fgr.UrlValues()
	require.NotNil(t, actual)
}

func TestFriendsGetRequest_SetUserId(t *testing.T) {
	fgr := Get().SetUserId(9000)
	actual := fgr.UrlValues().Get("user_id")
	require.EqualValues(t, "9000", actual)
}

func TestFriendsGetRequest_SetOrder(t *testing.T) {
	fgr := Get().SetOrder(Name)
	actual := fgr.UrlValues().Get("order")
	require.EqualValues(t, "name", actual)
}

func TestFriendsGetRequest_SetListId(t *testing.T) {
	fgr := Get().SetListId(7)
	actual := fgr.UrlValues().Get("list_id")
	require.EqualValues(t, "7", actual)
}

func TestFriendsGetRequest_SetCount(t *testing.T) {
	fgr := Get().SetCount(10000)
	actual := fgr.UrlValues().Get("count")
	require.EqualValues(t, "10000", actual)
}

func TestFriendsGetRequest_SetOffset(t *testing.T) {
	fgr := Get().SetOffset(4)
	actual := fgr.UrlValues().Get("offset")
	require.EqualValues(t, "4", actual)
}

func TestFriendsGetRequest_SetFields(t *testing.T) {
	fgr := Get().SetFields(Nickname, Sex, Domain)
	actual := fgr.UrlValues().Get("fields")
	require.EqualValues(t, "nickname,sex,domain", actual)
}

func TestFriendsGetRequest_Perform(t *testing.T) {
	fakeResponse := `{"response": {"count": 1,"items": [{"id": 12345,"first_name": "Alexander","last_name": "Ivanov",
"bdate": "20.2.1985","online": 0}]}}`
	actualRequest := &http.Request{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualRequest = r
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fakeResponse)
	}))
	defer ts.Close()

	requestSender := api.NewInstance(&oauth2.Token{AccessToken: "123"})
	baseUrl, urlParseErr := url.Parse(ts.URL)
	require.Nil(t, urlParseErr)
	requestSender.BaseUrl = baseUrl

	userlist, err := Get().SetOrder(Name).Perform(requestSender)
	require.Nil(t, err)
	assert.EqualValues(t, "/friends.get", actualRequest.RequestURI)
	assert.EqualValues(t, 1, userlist.Value.Count)
	require.EqualValues(t, 1, len(userlist.Value.Items))
	require.EqualValues(t, 12345, userlist.Value.Items[0].Id)
}
