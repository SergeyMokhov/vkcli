package api

import (
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"time"
)

type VkRequestSender interface {
	SendVkRequestAndRetryOnCaptcha(request vkRequest) (err error)
}

type Api struct {
	client       *http.Client
	token        *oauth2.Token
	BaseUrl      *url.URL
	speedLimiter <-chan bool
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

func (rb *Api) SendVkRequestAndRetryOnCaptcha(request vkRequest) (err error) {
	response, err := sendRequest(rb, request)
	if err != nil {
		return err
	}
	responseType := request.ResponseType()
	err = unmarshal(response, responseType)
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
			return unmarshal(response, responseType)
		}
	}
	return nil
}
