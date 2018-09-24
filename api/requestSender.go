package api

import (
	"encoding/json"
	"fmt"
	"gitlab.com/g00g/vkcli/api/obj"
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
	ResponseType() interface{}
}

type DummyVkRequest struct {
	Values                url.Values
	MethodStr             string
	ResponseStructPointer interface{}
}

func (dvk *DummyVkRequest) UrlValues() url.Values {
	return dvk.Values
}

func (dvk *DummyVkRequest) Method() string {
	return dvk.MethodStr
}

func (dvk *DummyVkRequest) ResponseType() interface{} {
	return dvk.ResponseStructPointer
}

func NewDummyVkRequest(method string, responseType interface{}) *DummyVkRequest {
	return &DummyVkRequest{
		Values:                url.Values{},
		MethodStr:             method,
		ResponseStructPointer: responseType}
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
	defaultParams.Add("v", "5.85")
	defaultParams.Add("access_token", rb.token.AccessToken)
}

func (rb *Api) AddSolvedCapture(request vkRequest, capture obj.VkErrorInfo, captureAnswer string) {
	p := request.UrlValues()
	p.Add("captcha_sid", capture.CaptchaSid)
	p.Add("captcha_key", captureAnswer)
}

func (rb *Api) SendRequest(request vkRequest) (err error) {
	<-rb.speedLimiter
	rb.addDefaultParams(request)
	method, errUrl := rb.BaseUrl.Parse(request.Method())
	if errUrl != nil {
		return fmt.Errorf("cannot parse method URL:%v", errUrl)
	}

	req, err := http.NewRequest("POST", method.String(), strings.NewReader(request.UrlValues().Encode()))
	if err != nil {
		return fmt.Errorf("error creating request:%v", err)
	}

	resp, errResp := rb.client.Do(req)
	if errResp != nil {
		return fmt.Errorf("error performing request:%v", errResp)
	}

	body, errReadBody := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if errReadBody != nil {
		return fmt.Errorf("cannot read response:%v", errReadBody)
	}

	return unmarshal(body, request.ResponseType())
}

func unmarshal(what []byte, to interface{}) (err error) {
	err = json.Unmarshal(what, to)
	if err != nil {
		err = fmt.Errorf("error parsing json to struct:%v", err)
	}
	return err
}
