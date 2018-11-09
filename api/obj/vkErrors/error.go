package vkErrors

const CaptchaRequired = 14

type ErrorInfo struct {
	ErrorCode  int    `json:"error_code"`
	ErrorMsg   string `json:"error_msg"`
	CaptchaSid string `json:"captcha_sid"`
	CaptchaImg string `json:"captcha_img"`
}

type Error struct {
	ErrorInfo `json:"error"`
}
