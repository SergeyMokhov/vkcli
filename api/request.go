package api

import (
	"gitlab.com/g00g/vk-cli/api/obj/vkErrors"
	"net/url"
)

type vkRequest interface {
	UrlValues() url.Values
	Method() string
	ResponseType() vkResponse
}

type vkResponse interface {
	GetError() *vkErrors.Error
}

type VkRequestBase struct {
	Values                url.Values
	MethodStr             string
	ResponseStructPointer vkResponse
}

func (dvk *VkRequestBase) UrlValues() url.Values {
	return dvk.Values
}

func (dvk *VkRequestBase) Method() string {
	return dvk.MethodStr
}

func (dvk *VkRequestBase) ResponseType() vkResponse {
	return dvk.ResponseStructPointer
}

func NewVkRequestBase(method string, responseType vkResponse) *VkRequestBase {
	return &VkRequestBase{
		Values:                url.Values{},
		MethodStr:             method,
		ResponseStructPointer: responseType}
}
