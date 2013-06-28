package o2aserver

type Storage interface {
	GetClient(clientId string) *Client
	SaveClient(client *Client) error

	GetAuthorize(code string) (*AuthorizationData, error)
	SaveAuthorize(parameters AuthorizationData) error
	RemoveAuthorize(code string) error

	GetAccessToken(code string) (*AccessTokenData, error)
	SaveAccessToken(parameters AccessTokenData) error
}
