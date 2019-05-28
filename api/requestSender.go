package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func addDefaultParams(request vkRequest, accessToken string) {
	defaultParams := request.UrlValues()
	if len(defaultParams["https"]) == 0 {
		defaultParams.Add("https", "1")
	}
	if len(defaultParams["v"]) == 0 {
		defaultParams.Add("v", "5.95")
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

func unmarshal(what []byte, to interface{}) (err error) {
	err = json.Unmarshal(what, to)
	if err != nil {
		err = fmt.Errorf("error parsing json to struct:%v", err)
	}
	return err
}
