package api

import (
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Vk struct {
	client  *http.Client
	token   *oauth2.Token
	baseUrl *url.URL
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

func (vk *Vk) DoIt() {
	//TODO return objects from request method.
	method, errU := vk.baseUrl.Parse("friends.get")
	if errU != nil {
		log.Fatalf("Cannot friends.get URL:%v", errU)
	}
	params := url.Values{}
	//TODO Extract default parameters including token into a separate method.
	params.Add("https", "1")
	params.Add("v", "5.80")
	//TODO Make each request configurable with optional parameters. Something like:
	// vk.Do(friends.Get{Order: "name"})
	params.Add("fields", "nickname, domain, sex, bdate, city, country, timezone, photo_50, photo_100, "+
		"photo_200_orig, has_mobile, contacts, education, online, relation, last_seen, status, can_write_private_message,"+
		" can_see_all_posts, can_post, universities")
	params.Add("order", "name")
	params.Add("access_token", vk.token.AccessToken)

	req, err := http.NewRequest("POST", method.String(), strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatalf("error creating request:%v", err)
	}
	resp, errResp := vk.client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//handle read response error
	}
	fmt.Printf("Error:%v\nResponse:%v\nBody:%s", errResp, resp, body)
}
