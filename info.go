package o2aserver

import (
	"net/http"
)

type Info struct {
	storage Storage
	appconfig AppConfig

	AccessToken string
	State string
}

func NewInfo(storage Storage, appconfig AppConfig) *Info {
	return &Info{
		storage: storage,
		appconfig: appconfig,
	}
}

func (a *Info) HandleInfoRequest(w *Response, r *http.Request) bool {
	r.ParseForm()

	a.AccessToken = r.Form.Get("access_token")
	a.State = r.Form.Get("state")

	token, err := a.storage.GetAccessToken(a.AccessToken)
	if err != nil {
		w.SetError(400, ErrorParameters{Error:"invalid_grant", Description:"The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.", State: a.State})
		return false
	}

	ret := InfoParameters{}
	ret.ClientId = token.ClientId
	ret.ExpiresIn = token.ExpiresIn
	ret.Scope = token.Scope
	ret.UserId = token.UserId
	ret.CreatedAt = token.CreatedAt

	pret := a.appconfig.ProcessInfoResponse(ret)

	w.SetParameters(pret)

	return true
}
