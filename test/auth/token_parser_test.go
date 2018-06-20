package auth

import (
	"github.com/SergeyMokhov/vkcli/auth"
	"strconv"
	"testing"
	"time"
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
