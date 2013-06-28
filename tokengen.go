package o2aserver

import (
	"encoding/base64"
	"code.google.com/p/go-uuid/uuid"
)

type TokenGen interface {
	GenerateAccessToken(data *AccessTokenData) error
	ParseAccessToken(data string) (interface{}, error)
}

type TokenGenDefault struct {

}

func (c *TokenGenDefault) GenerateAccessToken(data *AccessTokenData) error {
	data.AccessToken = uuid.New()
	data.AccessToken = base64.StdEncoding.EncodeToString([]byte(data.AccessToken))

	data.RefreshToken = uuid.New()
	data.RefreshToken = base64.StdEncoding.EncodeToString([]byte(data.RefreshToken))

	return nil
}

func (c *TokenGenDefault) ParseAccessToken(data string) (interface{}, error) {
	return nil, nil
}
