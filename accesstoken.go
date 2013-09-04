package o2aserver

import (
	"log"
	"net/http"
	"time"
)

type AccessTokenData struct {
	ClientId     string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	Scope        string
	RedirectUri  string
	UserId       string
	CreatedAt    time.Time
}

type AccessToken struct {
	storage   Storage
	tokengen  TokenGenAccess
	appconfig AppConfig

	Data          AccessTokenData
	State         string
	Client        *Client
	Authorization *AuthorizationData
}

func NewAccessToken(storage Storage, tokengen TokenGenAccess, appconfig AppConfig) *AccessToken {
	return &AccessToken{
		storage:   storage,
		tokengen:  tokengen,
		appconfig: appconfig,
		Data:      AccessTokenData{},
	}
}

func (a *AccessToken) HandleAccessTokenRequest(w *Response, r *http.Request) bool {
	r.ParseForm()

	grantType := r.Form.Get("grant_type")
	if grantType == "authorization_code" {
		return a.handleAuthorizationCode(w, r)
	}

	w.SetError(400, ErrorParameters{Error: "unsupported_grant_type", Description: "The authorization grant type is not supported by the authorization server.", State: a.State})
	return false
}

func (a *AccessToken) handleAuthorizationCode(w *Response, r *http.Request) bool {
	a.State = r.Form.Get("state")
	code := r.Form.Get("code")

	// must have a valid client
	a.Client = a.storage.GetClient(r.Form.Get("client_id"))
	if a.Client == nil {
		w.SetError(400, ErrorParameters{Error: "unauthorized_client", Description: "The client is not authorized to request an authorization code using this method", State: a.State})
		return false
	}
	a.Data.ClientId = a.Client.Id

	a.Data.RedirectUri = r.Form.Get("redirect_uri")
	if a.Data.RedirectUri == "" {
		a.Data.RedirectUri = a.Client.RedirectUri
	}

	// check redirect URI
	if !ValidateUri(a.Client.RedirectUri, a.Data.RedirectUri) {
		w.SetError(400, ErrorParameters{Error: "invalid_grant", Description: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.", State: a.State})
		return false
	}

	// load authorization code
	var err error
	a.Authorization, err = a.storage.GetAuthorize(code)
	if err != nil {
		w.SetError(400, ErrorParameters{Error: "invalid_grant", Description: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired or revoked.", State: a.State})
		return false
	}

	// redirect uri must match
	if a.Client.Id != a.Authorization.ClientId || a.Data.RedirectUri != a.Authorization.RedirectUri {
		w.SetError(400, ErrorParameters{Error: "invalid_grant", Description: "The provided authorization grant does not match the redirection URI used in the authorization request, or was issued to another client.", State: a.State})
		return false
	}

	// generate access token
	a.Data.CreatedAt = time.Now()
	a.Data.UserId = a.Authorization.UserId
	a.Data.ExpiresIn = 3600
	a.Data.Scope = a.Authorization.Scope

	if err := a.tokengen.GenerateAccessToken(&a.Data); err != nil {
		w.SetError(400, ErrorParameters{Error: "invalid_request", Description: "Server error.", State: a.State})
		return false
	}

	// save access token
	if err := a.storage.SaveAccessToken(a.Data); err != nil {
		w.SetError(400, ErrorParameters{Error: "invalid_request", Description: "Server error.", State: a.State})
		return false
	}

	// remove authorization token
	if err := a.storage.RemoveAuthorize(a.Authorization.Code); err != nil {
		log.Printf("Error: %s\n", err)
	}

	// return data
	ret := AccessTokenParameters{}
	ret.AccessToken = a.Data.AccessToken
	ret.TokenType = "bearer"
	ret.ExpiresIn = 3600
	ret.RefreshToken = a.Data.RefreshToken
	ret.Scope = a.Authorization.Scope
	ret.State = a.State

	pret := a.appconfig.ProcessAccessTokenResponse(ret)

	w.SetParameters(pret)
	return true
}
