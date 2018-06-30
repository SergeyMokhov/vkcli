package test

import (
	"gitlab.com/g00g/vkcli/auth"
	"strconv"
	"testing"
	"time"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

func TestParseString(t *testing.T) {
	expectedAccessToken := "533bacf01e11f55b536a565b57531ad114461ae8736d6506a3"
	expectedExpires := 86400
	url := "https://oauth.vk.com/blank.html#access_token=" + expectedAccessToken + "&expires_in" +
		"=" + strconv.Itoa(expectedExpires) + "&user_id=8492"

	token, err := auth.ParseUrlString(url)

	if err != nil {
		t.Fatalf("Got an error '%s', want nil", err)
	}

	if got, want := token.AccessToken, expectedAccessToken; got != want {
		t.Fatalf("Unexpected access token. Got %s, want %s", got, want)
	}

	if got, now := token.Expiry, time.Now(); !got.After(now) {
		t.Fatalf("Unexpected expiration. Token expiry is not in the future. Got: %s, now %s", got,
			now)
	}
}

func TestParseStringErrorUrl(t *testing.T) {
	url := "http://REDIRECT_URI?error=access_denied&error_description=The+user+or+authorization+" +
		"server+denied+the+request. "

	_, err := auth.ParseUrlString(url)

	if err == nil {
		t.Fatal("Got nil, want an error")
	}
}

func TestParseStringNoExpiry(t *testing.T) {
	expectedAccessToken := "533bacf01e11f55b536a565b57531ad114461ae8736d6506a3"
	url := "https://oauth.vk.com/blank.html#access_token=" + expectedAccessToken + "&user_id=8492"

	_, err := auth.ParseUrlString(url)

	if err == nil {
		t.Fatal("Got nil, want an error")
	}
}

func TestParseStringNonIntExpiry(t *testing.T) {
	expectedAccessToken := "533bacf01e11f55b536a565b57531ad114461ae8736d6506a3"
	expectedExpires := "str"
	url := "https://oauth.vk.com/blank.html#access_token=" + expectedAccessToken + "&expires_in" +
		"=" + expectedExpires + "&user_id=8492"

	_, err := auth.ParseUrlString(url)

	if err == nil {
		t.Fatal("Got nil, want an error")
	}
}

func TestAuthCodeURL(t *testing.T) {
	want := "https://oauth.vk.com/authorize?client_id=999&display=page&redirect_uri=https%3A%2F%2F" +
		"oauth.vk.com%2Fblank.html&response_type=token&scope=1026&state=randomState&v=5.80"
	c := oauth2.Config{
		ClientID:    "999",
		Endpoint:    vk.Endpoint,
		RedirectURL: "https://oauth.vk.com/blank.html",
		Scopes:      []string{"1026"},
	}
	conf := auth.Config{Config: c}
	opts := map[string]string{
		"display": "page",
		"v":       "5.80",
	}

	if got := conf.AuthCodeURL("randomState", opts); got != want {
		t.Fatalf("Unexpected URL.\nGot : '%s'\nWant: '%s'", got, want)
	}
}

func TestAuthCodeURLNoState(t *testing.T) {
	want := "https://oauth.vk.com/authorize?client_id=999&display=page&redirect_uri=https%3A%2F%2F" +
		"oauth.vk.com%2Fblank.html&response_type=token&scope=1026&state=&v=5.80"
	c := oauth2.Config{
		ClientID:    "999",
		Endpoint:    vk.Endpoint,
		RedirectURL: "https://oauth.vk.com/blank.html",
		Scopes:      []string{"1026"},
	}
	conf := auth.Config{Config: c}
	opts := map[string]string{
		"display": "page",
		"v":       "5.80",
	}

	if got := conf.AuthCodeURL("", opts); got == want {
		t.Fatalf("Unexpected URL. State is not generated automatically.\nUrl:%s", got)
	}
}
