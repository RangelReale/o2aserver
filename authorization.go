package o2aserver

import (
	"net/http"
	//"net/url"
	//"errors"
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
	"time"
)

type AuthorizationData struct {
	ClientId    string
	Code        string
	ExpiresIn   int64
	Scope       string
	RedirectUri string
	UserId      string
	CreatedAt   time.Time
}

type Authorization struct {
	storage   Storage
	appconfig AppConfig

	Data   AuthorizationData
	State  string
	Client *Client
}

func NewAuthorization(storage Storage, appconfig AppConfig) *Authorization {
	return &Authorization{
		storage:   storage,
		appconfig: appconfig,
		Data:      AuthorizationData{},
	}
}

// https://oauth2server.com/auth?response_type=code&client_id=CLIENT_ID&redirect_uri=REDIRECT_URI&scope=SCOPE&state=STATE
func (a *Authorization) HandleAuthorizeRequest(w *Response, r *http.Request) bool {
	r.ParseForm()

	a.State = r.Form.Get("state")
	a.Data.Scope = r.Form.Get("scope")

	// must have "response_type=code" parameter
	if r.Form.Get("response_type") != "code" {
		w.SetError(400, ErrorParameters{Error: "invalid_request", Description: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.", State: a.State})
		return false
	}

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

	a.Data.ExpiresIn = 600
	a.Data.CreatedAt = time.Now()

	// check redirect URI
	if !ValidateUri(a.Client.RedirectUri, a.Data.RedirectUri) {
		w.SetError(400, ErrorParameters{Error: "invalid_grant", Description: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.", State: a.State})
		return false
	}

	return true
}

func (a *Authorization) FinishAuthorizeRequest(w *Response, r *http.Request, authorized bool, userId string) {
	if authorized {
		// generate authorization token
		token := uuid.New()
		token = base64.StdEncoding.EncodeToString([]byte(token))

		// build response
		a.Data.Code = token
		a.Data.UserId = userId

		// save authorization token in storage
		a.storage.SaveAuthorize(a.Data)

		// build response
		var ret AuthorizeParameters
		ret.State = a.State
		ret.Code = token

		pret := a.appconfig.ProcessAuthorizeResponse(ret)

		w.SetRedirect(302, a.Data.RedirectUri, pret)
	} else {
		var ret ErrorParameters
		ret.Error = "access_denied"
		ret.Description = "The resource owner or authorization server denied the request."
		ret.State = a.State

		w.SetRedirect(302, a.Data.RedirectUri, ret)
	}
}
