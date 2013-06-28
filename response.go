package o2aserver

import (
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"reflect"
	"errors"
	"time"
	"fmt"
)

type Response struct {
	statusCode int
	statusText string
	parameters interface{}
	httpHeaders map[string]string
}

type ErrorParameters struct {
	Error string					`json:"error"`
	Description string				`json:"error_description"`
	Uri string						`json:"error_uri,omitempty"`
	State string					`json:"state,omitempty"`
}

type RedirectParameters struct {
	State string			`json:"state,omitempty"`
}

type AuthorizeParameters struct {
	RedirectParameters
	Code string				`json:"code"`
}

type AccessTokenParameters struct {
	AccessToken string				`json:"access_token"`
	TokenType string				`json:"token_type"`
	ExpiresIn int32					`json:"expires_in,omitempty"`
	RefreshToken string				`json:"refresh_token,omitempty"`
	Scope string					`json:"scope,omitempty"`
	State string			`json:"state,omitempty"`
}

type InfoParameters struct {
	ClientId string				`json:"client_id"`
	ExpiresIn int64				`json:"expires_in"`
	Scope string				`json:"scope"`
	UserId string				`json:"user_id"`
	CreatedAt time.Time			`json:"created_at"`
}


func NewResponse() *Response {
	ret := Response{}
	ret.Initialize()
	return &ret
}

func (r *Response) Initialize() {
	r.statusCode = 200
	r.httpHeaders = make(map[string]string)
}

func (r *Response) SetStatusCode(statuscode int) {
	r.statusCode = statuscode
}

func (r *Response) SetError(statuscode int, parameters ErrorParameters ) error {
	r.parameters = parameters
	r.httpHeaders["Cache-Control"] = "no-store"
	r.statusCode = statuscode

	return nil
}

func (r *Response) SetRedirect(statuscode int, redirectUrl string, parameters interface{}) error {
	gurl, err := url.Parse(redirectUrl)
	if err != nil {
		return err
	}


	v := gurl.Query()
	MarshalUrl(parameters, &v)
	//v.Add("code", parameters.State)
	gurl.RawQuery = v.Encode()

	//r.parameters = parameters
	r.parameters = struct{Uri string}{Uri: gurl.String()}
	r.statusCode = statuscode
	//r.httpHeaders["Location"] = gurl.String()

	return nil
}

func (r *Response) SetParameters(parameters interface{}) error {
	r.parameters = parameters
	return nil
}

func (r *Response) AddHttpHeader(header string, value string) {
	r.httpHeaders[header] = value
}

func (r *Response) AddHttpHeaders(headers map[string]string) {
	for i, k := range headers {
		r.httpHeaders[i] = k
	}
}

func (r* Response) Send(w http.ResponseWriter) error {
	for i, k := range r.httpHeaders {
		w.Header().Add(i, k)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	data, err := json.Marshal(r.parameters)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

func MarshalUrl(v interface{}, vurl *url.Values) (error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return errors.New("Only structs are supported")
	}

	val := reflect.ValueOf(v)
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		v := val.Field(i)
		if !p.Anonymous {
			tag := p.Tag.Get("json")
			if tag == "-" {
				continue
			}
			name, opts := parseTag(tag)
			if name == "" {
				name = strings.ToLower(p.Name)
			}
			if opts.Contains("omitempty") && isEmptyValue(v) {
				continue
			}
			vurl.Add(name, fmt.Sprint(v.Interface()))
		} else {
			err := MarshalUrl(v.Interface(), vurl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
