package auth

import (
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/url"
	"strconv"
	"time"
)

func ParseUrlString(urlStr string) (token *oauth2.Token, err error) {
	pastedUrl, errU := url.Parse(urlStr)
	if errU != nil {
		err = errors.New(fmt.Sprintf("Cannot parse URL from string '%s' %s", urlStr, errU))
		return nil, err
	}

	values := pastedUrl.Query()
	if errDesc := values.Get("error_description"); errDesc != "" {
		err = errors.New(fmt.Sprintf("Cannot parse token from the URL: %s", errDesc))
		return nil, err
	}

	if errParam := values.Get("error"); errParam != "" {
		err = errors.New(fmt.Sprintf("Cannot parse token from the URL: %s", errParam))
		return nil, err
	}

	urlWithToken, errUrl := pastedUrl.Parse(pastedUrl.Scheme + "://" + pastedUrl.Host +
		"?" + pastedUrl.Fragment)
	if errUrl != nil {
		err = errors.New(fmt.Sprintf("Cannot parse fragment to Url: %s", errUrl))
		return nil, err
	}

	values = urlWithToken.Query()
	accessToken := values.Get("access_token")
	if accessToken == "" {
		err = errors.New("Url does not contain an access token.")
		return nil, err
	}

	expiresIn := values.Get("expires_in")
	if expiresIn == "" {
		err = errors.New("Url does not contain expiration.")
		return nil, err
	}

	expiresInSeconds, errConv := strconv.Atoi(expiresIn)
	if errConv != nil {
		err = errors.New(fmt.Sprintf("Cannot parse token expiration time. %s", errConv))
		return nil, err
	}

	token = &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Duration(expiresInSeconds) * time.Second),
	}

	return token, err
}
