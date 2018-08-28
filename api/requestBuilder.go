package api

import "net/url"

type RequestBuilder struct {
	vk *Vk
}

func NewRequestBuilder(vk *Vk) *RequestBuilder {
	return &RequestBuilder{vk: vk}
}

func (rb *RequestBuilder) NewRequest(method string) (urlToMethod *url.URL,
	defaultParams url.Values, err error) {
	urlToMethod, err = rb.vk.baseUrl.Parse(method)
	defaultParams = url.Values{}
	defaultParams.Add("https", "1")
	defaultParams.Add("v", "5.80")
	defaultParams.Add("access_token", rb.vk.token.AccessToken)

	return urlToMethod, defaultParams, err
}
