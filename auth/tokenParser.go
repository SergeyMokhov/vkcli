package auth

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.com/g00g/vkcli/tools"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParseUrlString(urlStr string) (token *oauth2.Token, err error) {
	pastedUrl, errU := url.Parse(urlStr)
	if errU != nil {
		err = fmt.Errorf("Cannot parse URL from string '%s' %s", urlStr, errU)
		return nil, err
	}

	if urlIsErr, errDesc := isErr(pastedUrl); urlIsErr {
		return nil, fmt.Errorf("Cannot parse token from the URL: %s", errDesc)
	}

	urlWithToken, errUrl := pastedUrl.Parse(pastedUrl.Scheme + "://" + pastedUrl.Host +
		"?" + pastedUrl.Fragment)
	if errUrl != nil {
		err = fmt.Errorf("Cannot parse fragment to Url: %s", errUrl)
		return nil, err
	}

	values := urlWithToken.Query()
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
		err = fmt.Errorf("Cannot parse token expiration time. %s", errConv)
		return nil, err
	}

	token = &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Duration(expiresInSeconds) * time.Second),
	}

	return token, err
}

func (c *config) AuthCodeURL(state string, opts map[string]string) string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"token"},
		"client_id":     {c.ClientID},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	if state != "" {
		v.Set("state", state)
	} else {
		v.Set("state", tools.PseudoUuid())
	}
	for key, val := range opts {
		v.Set(key, val)
	}
	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

func (c *config) DefaultOptions() map[string]string {
	return map[string]string{
		"display": "page",
		"v":       "5.80",
	}
}

type config struct {
	oauth2.Config
}

func NewConfig(clientId string) config {
	c := oauth2.Config{
		ClientID:    clientId,
		Endpoint:    vk.Endpoint,
		RedirectURL: "https://oauth.vk.com/blank.html",
		Scopes:      []string{strconv.FormatInt(int64(FullUserScope()), 10)},
	}
	return config{Config: c}
}

func isErr(urlToTest *url.URL) (isErr bool, err string) {
	values := urlToTest.Query()
	if errDesc := values.Get("error_description"); errDesc != "" {
		isErr = true
		err = errDesc
	}

	if errParam := values.Get("error"); errParam != "" {
		isErr = true
		err = err + ". " + errParam
	}
	return
}
