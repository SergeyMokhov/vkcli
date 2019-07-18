package requests

import (
	"encoding/json"
	"fmt"
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

func Unmarshal(what []byte, to interface{}) (err error) {
	err = json.Unmarshal(what, to)
	if err != nil {
		err = fmt.Errorf("error parsing json to struct:%v", err)
	}
	return err
}

type FakeVkRequest struct {
	*VkRequestBase
}

func (fg *FakeVkRequest) UrlValues() url.Values {
	return fg.Values
}

func (fg *FakeVkRequest) Method() string {
	return fg.MethodStr
}

func (fg *FakeVkRequest) ResponseType() vkResponse {
	return fg.ResponseStructPointer
}

type FakeVkResponse struct {
	*vkErrors.Error
}

func (fr *FakeVkResponse) GetError() *vkErrors.Error {
	return fr.Error
}
