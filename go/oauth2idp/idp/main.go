//go:generate wire .

package idp

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/ory/fosite"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
)

func SetupIDP(swPlugin *swagger.Plugin) {

	idpProvider, err := InitializeProvider()
	if err != nil {
		log.Fatal(err)
	}

	h := &handlers{
		Provider: idpProvider,
	}

	ucon.HandleFunc("GET", "/oauth2/auth", h.AuthHTML)
	ucon.HandleFunc("POST", "/oauth2/auth", h.AuthEndpoint)
	ucon.HandleFunc("POST", "/oauth2/token", h.TokenEndpoint)
	ucon.HandleFunc("POST", "/oauth2/revoke", h.RevokeEndpoint)
	ucon.HandleFunc("POST", "/oauth2/introspect", h.IntrospectEndpoint)
}

type handlers struct {
	Provider fosite.OAuth2Provider
}

type AuthHTMLRequest struct {
	ClientID     string `json:"client_id" swagger:"in=query"`
	RedirectURL  string `json:"redirect_uri" swagger:"in=query"`
	ResponseType string `json:"response_type" swagger:"in=query"`
	Scope        string `json:"scope" swagger:"in=query"`
	State        string `json:"state" swagger:"in=query"`
	Nonce        string `json:"nonce" swagger:"in=query"`
}

func (h *handlers) AuthHTML(w http.ResponseWriter, req *AuthHTMLRequest) error {
	data, err := ioutil.ReadFile("./public/idp/auth.html.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("auth").Parse(string(data))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"scopes": strings.Split(req.Scope, " "),
	})
	if err != nil {
		return err
	}

	return nil
}

type AuthEndpointRequest struct {
	Scopes []string `json:"scopes"`
}

func (h *handlers) AuthEndpoint(ctx context.Context, r *http.Request, w http.ResponseWriter, req *AuthEndpointRequest, user *domains.User) {

	ar, err := h.Provider.NewAuthorizeRequest(ctx, r)
	if err != nil {
		h.Provider.WriteAuthorizeError(w, ar, err)
		return
	}

	if user == nil {
		h.Provider.WriteAuthorizeError(w, ar, fmt.Errorf("login required"))
		return
	}

	for _, scope := range req.Scopes {
		ar.GrantScope(scope)
	}

	sessionData, err := InitializeSession(ctx, user)
	if err != nil {
		h.Provider.WriteAuthorizeError(w, ar, err)
		return
	}

	response, err := h.Provider.NewAuthorizeResponse(ctx, ar, sessionData)
	if err != nil {
		h.Provider.WriteAuthorizeError(w, ar, err)
		return
	}

	h.Provider.WriteAuthorizeResponse(w, ar, response)
}

func (h *handlers) TokenEndpoint(ctx context.Context, r *http.Request, w http.ResponseWriter, req *AuthEndpointRequest, user *domains.User) {
	sessionData, err := ProvideSession(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	accessRequest, err := h.Provider.NewAccessRequest(ctx, r, sessionData)
	if err != nil {
		h.Provider.WriteAccessError(w, accessRequest, err)
		return
	}

	if accessRequest.GetGrantTypes().Exact("client_credentials") {
		for _, scope := range accessRequest.GetRequestedScopes() {
			if fosite.HierarchicScopeStrategy(accessRequest.GetClient().GetScopes(), scope) {
				accessRequest.GrantScope(scope)
			}
		}
	}

	response, err := h.Provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		h.Provider.WriteAccessError(w, accessRequest, err)
		return
	}

	h.Provider.WriteAccessResponse(w, accessRequest, response)
}

func (h *handlers) RevokeEndpoint(ctx context.Context, r *http.Request, w http.ResponseWriter) {
	err := h.Provider.NewRevocationRequest(ctx, r)
	h.Provider.WriteRevocationResponse(w, err)
}

func (h *handlers) IntrospectEndpoint(ctx context.Context, r *http.Request, w http.ResponseWriter) {
	sessionData, err := ProvideSession(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	ir, err := h.Provider.NewIntrospectionRequest(ctx, r, sessionData)
	if err != nil {
		h.Provider.WriteIntrospectionError(w, err)
		return
	}

	h.Provider.WriteIntrospectionResponse(w, ir)
}
