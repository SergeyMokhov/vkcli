package requests

import (
	"bufio"
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

type RequestSendRetryer interface {
	SendVkRequestAndRetryOnCaptcha(request vkRequest) error
}

type VkRequestSender struct {
	client       *http.Client
	token        *oauth2.Token
	BaseUrl      *url.URL
	speedLimiter <-chan bool
}

func NewVkRequestSender(token *oauth2.Token) *VkRequestSender {
	defaultSpeedLimit := 3
	defaultSpeedUnit := time.Second

	apiUrl, err := url.Parse("https://api.vk.com/method/")
	if err != nil {
		log.Fatalf("cannot parse VK api URL:%v", err)
	}

	return &VkRequestSender{
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
		defaultParams.Add("v", "5.101")
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

func promptForCaptcha(vkErr *vkErrors.Error) (answer string) {
	fmt.Printf("Please, solve the captcha: %v\nCaptcha unswer is: ", vkErr.CaptchaImg)
	reader := bufio.NewReader(os.Stdin)
	_, scanErr := fmt.Fscanln(reader, &answer)
	if scanErr != nil {
		fmt.Printf("Error reading captcha unswer: %v", scanErr)
		//todo add proper logging
	}
	return answer
}

func sendRequest(rb *VkRequestSender, request vkRequest) (body []byte, err error) {
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
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			fmt.Printf("Cannot close request body: %v", closeErr)
			//Todo add proper logging
		}
	}()
	if errReadBody != nil {
		return body, fmt.Errorf("cannot read response:%v", errReadBody)
	}

	return body, nil
}

func (rb *VkRequestSender) SendVkRequestAndRetryOnCaptcha(request vkRequest) (err error) {
	response, err := sendRequest(rb, request)
	if err != nil {
		return err
	}
	responseType := request.ResponseType()
	err = Unmarshal(response, responseType)
	if err != nil {
		return err
	}
	vkErr := responseType.GetError()
	if vkErr != nil {
		//TODO make amount of retries configurable. User might enter incorrect captcha multiple times
		if vkErr.ErrorCode == vkErrors.CaptchaRequired {
			captcha := promptForCaptcha(vkErr)
			addSolvedCaptcha(request, vkErr, captcha)
			response, err := sendRequest(rb, request)
			if err != nil {
				return err
			}
			return Unmarshal(response, responseType)
		}
	}
	return nil
}
