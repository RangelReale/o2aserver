package o2aserver

import (
)

type AppConfig interface {
	ProcessAuthorizeResponse(parameters AuthorizeParameters) interface{}
	ProcessAccessTokenResponse(parameters AccessTokenParameters) interface{}
	ProcessInfoResponse(parameters InfoParameters) interface{}
}

type AppConfigDefault struct {

}

func (c *AppConfigDefault) ProcessAuthorizeResponse(parameters AuthorizeParameters) interface{} {
	return parameters
}

func (c *AppConfigDefault) ProcessAccessTokenResponse(parameters AccessTokenParameters) interface{} {
	return parameters
}

func (c *AppConfigDefault) ProcessInfoResponse(parameters InfoParameters) interface{} {
	return parameters
}
