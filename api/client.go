package api

import (
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Vk struct {
	paramLock sync.Mutex
	client    *http.Client
	token     *oauth2.Token
	baseUrl   *url.URL
	//TODO add auto scheduler that will queue and limit speed to 3 request per second (VK limitation)
	// see https://vk.com/dev/api_requests
}

func NewVk(token *oauth2.Token) *Vk {
	apiUrl, err := url.Parse("https://api.vk.com/method/")
	if err != nil {
		log.Fatalf("Cannot parse VK api URL:%v", err)
	}

	return &Vk{
		client:  &http.Client{},
		token:   token,
		baseUrl: apiUrl,
	}
}

func (vk *Vk) newRequest(method string) (urlToMethod *url.URL, defaultParams url.Values, err error) {
	vk.paramLock.Lock()
	defer vk.paramLock.Unlock()

	urlToMethod, err = vk.baseUrl.Parse(method)
	defaultParams = url.Values{}
	defaultParams.Add("https", "1")
	defaultParams.Add("v", "5.80")
	defaultParams.Add("access_token", vk.token.AccessToken)

	return urlToMethod, defaultParams, err
}

func (vk *Vk) DoIt() {
	method, params, err := vk.newRequest("friends.get")
	if err != nil {
		log.Fatalf("error preparing request:%v", err)
	}
	//TODO return objects from request method.
	//TODO Make each request configurable with optional parameters. Something like:
	// vk.Do(friends.Get{Order: "name"})
	params.Add("fields", "nickname, domain, sex, bdate, city, country, timezone, photo_50, photo_100, "+
		"photo_200_orig, has_mobile, contacts, education, online, relation, last_seen, status, can_write_private_message,"+
		" can_see_all_posts, can_post, universities")
	params.Add("order", "name")

	req, err := http.NewRequest("POST", method.String(), strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatalf("error creating request:%v", err)
	}
	resp, errResp := vk.client.Do(req)
	if errResp != nil {
		fmt.Printf("Error:%v", errResp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		//handle read response error
	}
	fmt.Printf("Error reading response:%v\nResponse:%v\nBody:%s", err, resp, body)

}
