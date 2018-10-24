package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitlab.com/g00g/vkcli/api/obj"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

func addDefaultParams(request vkRequest, accessToken string) {
	defaultParams := request.UrlValues()
	defaultParams.Add("https", "1")
	defaultParams.Add("v", "5.85")
	defaultParams.Add("access_token", accessToken)
}

func addSolvedCaptcha(request vkRequest, capture *obj.Error, captureAnswer string) {
	p := request.UrlValues()
	p.Add("captcha_sid", capture.CaptchaSid)
	p.Add("captcha_key", captureAnswer)
}

func (rb *Api) SendRequestAndRetyOnCaptcha(request vkRequest) (err error) {
	return sendVkRequestAndRetyOnCaptcha(rb, request)
}

func sendVkRequestAndRetyOnCaptcha(rb *Api, request vkRequest) (err error) {
	<-rb.speedLimiter
	response, err := sendRequest(request, rb.BaseUrl, rb.token.AccessToken, rb.client)
	if err != nil {
		return err
	}

	vkErr := &obj.Error{}
	err = unmarshal(response, vkErr)
	if err != nil {
		return err
	}

	if vkErr.ErrorCode == obj.CaptchaRequired {
		captcha := promptForCaptcha(vkErr)
		addSolvedCaptcha(request, vkErr, captcha)
		response, err := sendRequest(request, rb.BaseUrl, rb.token.AccessToken, rb.client)
		if err != nil {
			return err
		}
		return unmarshal(response, request.ResponseType())
	}

	return unmarshal(response, request.ResponseType())
}

//TODO use monkeypatch to test this part.
func promptForCaptcha(vkErr *obj.Error) (answer string) {
	fmt.Printf("Please, solve the captcha: %v\nCaptcha unswer is: ", vkErr.CaptchaImg)
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscanln(reader, &answer)
	return answer
}

func (rb *Api) SendRequest(request vkRequest) (err error) {
	return sendVkRequest(rb, request)
}

func sendVkRequest(rb *Api, request vkRequest) (err error) {
	<-rb.speedLimiter
	response, err := sendRequest(request, rb.BaseUrl, rb.token.AccessToken, rb.client)
	if err != nil {
		return err
	}

	return unmarshal(response, request.ResponseType())
}

func sendRequest(request vkRequest, baseUrl *url.URL, accessToken string, client *http.Client) (body []byte, err error) {
	addDefaultParams(request, accessToken)
	method, errUrl := baseUrl.Parse(request.Method())
	if errUrl != nil {
		return body, fmt.Errorf("cannot parse method URL:%v", errUrl)
	}

	req, err := http.NewRequest("POST", method.String(), strings.NewReader(request.UrlValues().Encode()))
	if err != nil {
		return body, fmt.Errorf("error creating request:%v", err)
	}

	resp, errResp := client.Do(req)
	if errResp != nil {
		return body, fmt.Errorf("error performing request:%v", errResp)
	}

	body, errReadBody := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if errReadBody != nil {
		return body, fmt.Errorf("cannot read response:%v", errReadBody)
	}

	return body, nil
}

func unmarshal(what []byte, to interface{}) (err error) {
	err = json.Unmarshal(what, to)
	if err != nil {
		err = fmt.Errorf("error parsing json to struct:%v", err)
	}
	return err
}
