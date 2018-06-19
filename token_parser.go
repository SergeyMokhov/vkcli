package vkcli

import (
	"golang.org/x/oauth2"
	"log"
	"net/url"
	"strconv"
	"time"
)

func ParseUrl(url *url.URL) *oauth2.Token {
	//TODO add tests
	values := url.Query()
	accessToken := values.Get("access_token")
	if accessToken == "" {
		log.Fatal("Url does not contain an access token.")
	}

	expiresIn := values.Get("expires_in")
	if expiresIn == "" {
		log.Fatal("Url does not contain expiration.")
	}
	expiresInSeconds, err := strconv.Atoi(expiresIn)
	if err != nil {
		log.Fatalf("Cannot parse token expiration time, is it integer? &v", err)
	}

	return &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Duration(expiresInSeconds) * time.Second),
	}
}

func ParseString(urlString string) *oauth2.Token {
	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatalf("Cannot parse URL from string '%s' %v", urlString, err)
	}
	return ParseUrl(u)
}
