package api

import (
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Api struct {
	client       *http.Client
	token        *oauth2.Token
	BaseUrl      *url.URL
	speedLimiter <-chan bool
}

type vkRequest interface {
	UrlValues() url.Values
	Method() string
}

func NewInstance(token *oauth2.Token) *Api {
	defaultSpeedLimit := 3
	defaultSpeedUnit := time.Second

	apiUrl, err := url.Parse("https://api.vk.com/method/")
	if err != nil {
		log.Fatalf("cannot parse VK api URL:%v", err)
	}

	return &Api{
		client:       &http.Client{},
		token:        token,
		BaseUrl:      apiUrl,
		speedLimiter: NewSpeedLimiter(defaultSpeedLimit, defaultSpeedUnit).Channel(),
	}
}

func (rb *Api) addDefaultParams(request vkRequest) {
	defaultParams := request.UrlValues()
	defaultParams.Add("https", "1")
	defaultParams.Add("v", "5.84")
	defaultParams.Add("access_token", rb.token.AccessToken)
}

func (rb *Api) Perform(request vkRequest) (responseBody []byte, err error) {
	<-rb.speedLimiter
	rb.addDefaultParams(request)
	method, errUrl := rb.BaseUrl.Parse(request.Method())
	if errUrl != nil {
		return []byte{}, fmt.Errorf("cannot parse method URL:%v", errUrl)
	}

	req, err := http.NewRequest("POST", method.String(), strings.NewReader(request.UrlValues().Encode()))
	if err != nil {
		return []byte{}, fmt.Errorf("error creating request:%v", err)
	}

	resp, errResp := rb.client.Do(req)
	if errResp != nil {
		return []byte{}, fmt.Errorf("error performing request:%v", errResp)
	}

	body, errReadBody := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if errReadBody != nil {
		return []byte{}, fmt.Errorf("cannot read response:%v", errReadBody)
	}
	return body, nil
}
