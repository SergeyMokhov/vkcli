package api

import "net/url"

type RequestBuilder struct {
	vk *Vk
}

type vkRequest interface {
	UrlValues() url.Values
	Method() string
}

func NewRequestBuilder(vk *Vk) *RequestBuilder {
	return &RequestBuilder{vk: vk}
}

func (rb *RequestBuilder) NewRequest(request vkRequest) (urlToMethod *url.URL, defaultParams url.Values, err error) {
	urlToMethod, err = rb.vk.baseUrl.Parse(request.Method())
	defaultParams = request.UrlValues()
	defaultParams.Add("https", "1")
	defaultParams.Add("v", "5.84")
	defaultParams.Add("access_token", rb.vk.token.AccessToken)

	return urlToMethod, defaultParams, err
}
