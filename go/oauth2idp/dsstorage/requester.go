package dsstorage

import (
	"context"
	"encoding/json"
	"github.com/ory/fosite"
	"go.mercari.io/datastore"
	"net/url"
	"time"
)

var _ fosite.Requester = (*DefaultRequester)(nil)
var _ fosite.AccessRequester = (*DefaultRequester)(nil)
var _ fosite.AuthorizeRequester = (*DefaultRequester)(nil)

var _ datastore.PropertyLoadSaver = (*DefaultRequester)(nil)

var _ RequestInvalidator = (*DefaultRequester)(nil)
var _ ClientLoader = (*DefaultRequester)(nil)
var _ SessionLoader = (*DefaultRequester)(nil)

type RequestInvalidator interface {
	IsActive() bool
	SetActive(active bool)
}

type ClientLoader interface {
	GetClientID() string
	SetClient(client fosite.Client)
}

type SessionLoader interface {
	LoadSession(ctx context.Context, session fosite.Session) error
}

type DefaultRequester struct {
	// for fosite.Request
	ID                string         `` // Not Datastore Key
	RequestedAt       time.Time      ``
	ClientID          string         `json:"-"`
	Client            fosite.Client  `datastore:"-"`
	RequestedScope    []string       ``
	GrantedScope      []string       ``
	EncodedForm       string         `json:"-"`
	Form              url.Values     `datastore:"-"`
	SessionJSON       string         `json:"-"`
	Session           fosite.Session `datastore:"-"`
	RequestedAudience []string       ``
	GrantedAudience   []string       ``
	// for fosite.AccessRequest
	GrantTypes       []string ``
	HandledGrantType []string ``
	// for fosite.AuthorizeRequest
	ResponseTypes        []string ``
	RedirectURI          string   ``
	State                string   ``
	HandledResponseTypes []string ``
	// others...
	Active    bool      ``
	UpdatedAt time.Time ``
	CreatedAt time.Time ``
}

func (r *DefaultRequester) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, r, ps)
	if err != nil {
		return err
	}

	r.Form, err = url.ParseQuery(r.EncodedForm)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultRequester) Save(ctx context.Context) ([]datastore.Property, error) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	r.UpdatedAt = time.Now()

	if r.Client != nil {
		r.ClientID = r.Client.GetID()
	} else {
		r.ClientID = ""
	}

	if r.Session != nil {
		b, err := json.Marshal(r.Session)
		if err != nil {
			return nil, err
		}
		r.SessionJSON = string(b)
	}

	r.EncodedForm = r.Form.Encode()

	return datastore.SaveStruct(ctx, r)
}

func (r *DefaultRequester) IsActive() bool {
	return r.Active
}

func (r *DefaultRequester) SetActive(active bool) {
	r.Active = active
}

func (r *DefaultRequester) GetClientID() string {
	return r.ClientID
}

func (r *DefaultRequester) SetClient(client fosite.Client) {
	r.Client = client
}

func (r *DefaultRequester) LoadSession(ctx context.Context, session fosite.Session) error {
	if r.SessionJSON != "" {
		err := json.Unmarshal([]byte(r.SessionJSON), session)
		if err != nil {
			return err
		}
		r.Session = session
	} else {
		r.Session = nil
	}

	return nil
}

func (r *DefaultRequester) SetID(id string) {
	r.ID = id
}

func (r *DefaultRequester) GetID() string {
	return r.ID
}

func (r *DefaultRequester) GetRequestedAt() time.Time {
	return r.RequestedAt
}

func (r *DefaultRequester) GetClient() fosite.Client {
	return r.Client
}

func (r *DefaultRequester) GetRequestedScopes() fosite.Arguments {
	return r.RequestedScope
}

func (r *DefaultRequester) GetRequestedAudience() fosite.Arguments {
	return r.RequestedAudience
}

func (r *DefaultRequester) SetRequestedScopes(scopes fosite.Arguments) {
	r.RequestedScope = nil
	for _, scope := range scopes {
		r.AppendRequestedScope(scope)
	}
}

func (r *DefaultRequester) SetRequestedAudience(audience fosite.Arguments) {
	r.RequestedAudience = nil
	for _, a := range audience {
		r.AppendRequestedAudience(a)
	}
}

func (r *DefaultRequester) AppendRequestedScope(scope string) {
	for _, old := range r.RequestedScope {
		if scope == old {
			return
		}
	}
	r.RequestedScope = append(r.RequestedScope, scope)
}

func (r *DefaultRequester) AppendRequestedAudience(s string) {
	for _, old := range r.RequestedAudience {
		if s == old {
			return
		}
	}
	r.RequestedAudience = append(r.RequestedAudience, s)
}

func (r *DefaultRequester) GetGrantedScopes() fosite.Arguments {
	return r.GrantedScope
}

func (r *DefaultRequester) GetGrantedAudience() fosite.Arguments {
	return r.GrantedAudience
}

func (r *DefaultRequester) GrantScope(scope string) {
	for _, old := range r.GrantedScope {
		if scope == old {
			return
		}
	}
	r.GrantedScope = append(r.GrantedScope, scope)
}

func (r *DefaultRequester) GrantAudience(audience string) {
	for _, old := range r.GrantedAudience {
		if audience == old {
			return
		}
	}
	r.GrantedAudience = append(r.GrantedAudience, audience)
}

func (r *DefaultRequester) GetSession() (session fosite.Session) {
	return r.Session
}

func (r *DefaultRequester) SetSession(session fosite.Session) {
	r.Session = session
}

func (r *DefaultRequester) GetRequestForm() url.Values {
	return r.Form
}

func (r *DefaultRequester) Merge(requester fosite.Requester) {
	for _, scope := range requester.GetRequestedScopes() {
		r.AppendRequestedScope(scope)
	}
	for _, scope := range requester.GetGrantedScopes() {
		r.GrantScope(scope)
	}

	for _, aud := range requester.GetRequestedAudience() {
		r.AppendRequestedAudience(aud)
	}
	for _, aud := range requester.GetGrantedAudience() {
		r.GrantAudience(aud)
	}

	r.RequestedAt = requester.GetRequestedAt()
	r.Client = requester.GetClient()
	r.Session = requester.GetSession()

	for k, v := range requester.GetRequestForm() {
		r.Form[k] = v
	}
}

func (r *DefaultRequester) Sanitize(allowedParameters []string) fosite.Requester {
	n := &DefaultRequester{}
	allowed := make(map[string]bool)
	for _, v := range allowedParameters {
		allowed[v] = true
	}

	*n = *r
	n.ID = r.GetID()
	n.Form = url.Values{}
	for k := range r.Form {
		if _, ok := allowed[k]; ok {
			n.Form.Add(k, r.Form.Get(k))
		}
	}

	return n
}

func (r *DefaultRequester) GetGrantTypes() (grantTypes fosite.Arguments) {
	return r.GrantTypes
}

func (r *DefaultRequester) GetResponseTypes() (responseTypes fosite.Arguments) {
	return r.ResponseTypes
}

func (r *DefaultRequester) SetResponseTypeHandled(responseType string) {
	r.HandledResponseTypes = append(r.HandledResponseTypes, responseType)
}

func (r *DefaultRequester) DidHandleAllResponseTypes() bool {
	for _, rt := range r.ResponseTypes {
		for _, handle := range r.HandledResponseTypes {
			if rt == handle {
				return false
			}
		}
	}

	return len(r.ResponseTypes) > 0

}

func (r *DefaultRequester) GetRedirectURI() *url.URL {
	if r.RedirectURI == "" {
		return nil
	}
	redirectURL, err := url.Parse(r.RedirectURI)
	if err != nil {
		return nil
	}
	return redirectURL
}

func (r *DefaultRequester) IsRedirectURIValid() (isValid bool) {
	if r.GetRedirectURI() == nil {
		return false
	}

	raw := r.GetRedirectURI().String()
	if r.GetClient() == nil {
		return false
	}

	redirectURI, err := fosite.MatchRedirectURIWithClientRedirectURIs(raw, r.GetClient())
	if err != nil {
		return false
	}
	return fosite.IsValidRedirectURI(redirectURI)

}

func (r *DefaultRequester) GetState() (state string) {
	return r.State
}
