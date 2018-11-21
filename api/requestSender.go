package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
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

type VkRequestBase struct {
	Values                url.Values
	MethodStr             string
	ResponseStructPointer interface{}
}

func (dvk *VkRequestBase) UrlValues() url.Values {
	return dvk.Values
}

func (dvk *VkRequestBase) Method() string {
	return dvk.MethodStr
}

func (dvk *VkRequestBase) ResponseType() interface{} {
	return dvk.ResponseStructPointer
}

func NewVkRequestBase(method string, responseType interface{}) *VkRequestBase {
	return &VkRequestBase{
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
	if len(defaultParams["https"]) == 0 {
		defaultParams.Add("https", "1")
	}
	if len(defaultParams["v"]) == 0 {
		defaultParams.Add("v", "5.85")
	}
	if len(defaultParams["access_token"]) == 0 {
		defaultParams.Add("access_token", accessToken)
	}
}

func addSolvedCaptcha(request vkRequest, captcha *vkErrors.Error, captchaAnswer string) {
	p := request.UrlValues()
	p.Add("captcha_sid", captcha.CaptchaSid)
	p.Add("captcha_key", captchaAnswer)
}

func (rb *Api) SendRequestAndRetyOnCaptcha(request vkRequest) (err error) {
	return sendVkRequestAndRetyOnCaptcha(rb, request)
}

func sendVkRequestAndRetyOnCaptcha(rb *Api, request vkRequest) (err error) {
	response, err := sendRequest(rb, request)
	if err != nil {
		return err
	}

	vkErr := &vkErrors.Error{}
	err = unmarshal(response, vkErr)
	if err != nil {
		return err
	}
	//TODO make amount of retries configurable. User might enter incorrect captcha multiple times
	if vkErr.ErrorCode == vkErrors.CaptchaRequired {
		captcha := promptForCaptcha(vkErr)
		addSolvedCaptcha(request, vkErr, captcha)
		response, err := sendRequest(rb, request)
		if err != nil {
			return err
		}
		return unmarshal(response, request.ResponseType())
	}

	return unmarshal(response, request.ResponseType())
}

func promptForCaptcha(vkErr *vkErrors.Error) (answer string) {
	fmt.Printf("Please, solve the captcha: %v\nCaptcha unswer is: ", vkErr.CaptchaImg)
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscanln(reader, &answer)
	return answer
}

func sendVkRequest(rb *Api, request vkRequest) (err error) {
	response, err := sendRequest(rb, request)
	if err != nil {
		return err
	}

	return unmarshal(response, request.ResponseType())
}

func sendRequest(rb *Api, request vkRequest) (body []byte, err error) {
	addDefaultParams(request, rb.token.AccessToken)
	method, errUrl := rb.BaseUrl.Parse(request.Method())
	if errUrl != nil {
		return body, fmt.Errorf("cannot parse method URL:%v", errUrl)
	}

	req, err := http.NewRequest("POST", method.String(), strings.NewReader(request.UrlValues().Encode()))
	if err != nil {
		return body, fmt.Errorf("error creating request:%v", err)
	}

	<-rb.speedLimiter
	resp, errResp := rb.client.Do(req)
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
