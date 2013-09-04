package o2aserver

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
)

type TokenGenAuthorization interface {
	GenerateAuthorizationToken(data *AuthorizationData) error
	ParseAuthorizationToken(data string) (interface{}, error)
}

type TokenGenAccess interface {
	GenerateAccessToken(data *AccessTokenData) error
	ParseAccessToken(data string) (interface{}, error)
}

type TokenGenAuthorizationDefault struct {
}

func (c *TokenGenAuthorizationDefault) GenerateAuthorizationToken(data *AuthorizationData) error {
	// generate authorization token
	token := uuid.New()
	data.Code = base64.StdEncoding.EncodeToString([]byte(token))
	return nil
}

func (c *TokenGenAuthorizationDefault) ParseAuthorizationToken(data string) (interface{}, error) {
	return nil, nil
}

type TokenGenAccessDefault struct {
}

func (c *TokenGenAccessDefault) GenerateAccessToken(data *AccessTokenData) error {
	data.AccessToken = uuid.New()
	data.AccessToken = base64.StdEncoding.EncodeToString([]byte(data.AccessToken))

	data.RefreshToken = uuid.New()
	data.RefreshToken = base64.StdEncoding.EncodeToString([]byte(data.RefreshToken))

	return nil
}

func (c *TokenGenAccessDefault) ParseAccessToken(data string) (interface{}, error) {
	return nil, nil
}
