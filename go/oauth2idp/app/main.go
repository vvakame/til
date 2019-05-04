package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// OAuth2のクライアント側になるアプリの世界…

var baseURL string

// アプリのClient アプリに対して固定値
var clientConf *oauth2.Config

var appClientConf *clientcredentials.Config

func init() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	clientConf = &oauth2.Config{
		ClientID:     "my-client",
		ClientSecret: "foobar",
		RedirectURL:  baseURL + "/callback",
		Scopes:       []string{"photos", "openid", "offline"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   baseURL + "/oauth2/auth",
			TokenURL:  baseURL + "/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams, // client_secret_post
		},
	}

	appClientConf = &clientcredentials.Config{
		ClientID:     "my-client",
		ClientSecret: "foobar",
		Scopes:       []string{"fosite"},
		TokenURL:     baseURL + "/oauth2/token",
	}
}

func SetupAppAPI(swPlugin *swagger.Plugin) {
	ucon.HandleFunc("GET", "/", indexHandler)
	ucon.HandleFunc("GET", "/client", clientHandler)
	ucon.HandleFunc("GET,POST", "/owner", ownerHandler)
	ucon.HandleFunc("GET", "/callback", callbackHandler)
	ucon.HandleFunc("GET", "/protected", protectedHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	data, err := ioutil.ReadFile("./public/app/index.html.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("index").Parse(string(data))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"authLinkURL":      clientConf.AuthCodeURL("some-random-state-foobar") + "&nonce=some-random-nonce",
		"implicitGrantURL": "/oauth2/auth?client_id=my-client&redirect_uri=" + url.QueryEscape(baseURL+"/callback") + "&response_type=token%20id_token&scope=fosite%20openid&state=some-random-state-foobar&nonce=some-random-nonce",
		"refreshGrantURL":  clientConf.AuthCodeURL("some-random-state-foobar") + "&nonce=some-random-nonce",
		"invalidAccessURL": "/oauth2/auth?client_id=my-client&scope=fosite&response_type=123&redirect_uri=" + url.QueryEscape(baseURL+"/callback"),
	})
	if err != nil {
		return err
	}

	return nil
}

func clientHandler(w http.ResponseWriter, r *http.Request, req *CallbackRequest) error {
	data, err := ioutil.ReadFile("./public/app/client.html.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("client").Parse(string(data))
	if err != nil {
		return err
	}

	token, err := appClientConf.Token(r.Context())

	p := map[string]interface{}{
		"error": "",
	}
	if err != nil {
		p["error"] = err.Error()
	} else {
		p["accessToken"] = token.AccessToken
		p["token"] = fmt.Sprintf("%#v", token)
	}

	err = tmpl.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}

type OwnerRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func ownerHandler(w http.ResponseWriter, r *http.Request, req *OwnerRequest) error {
	data, err := ioutil.ReadFile("./public/app/owner.html.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("owner").Parse(string(data))
	if err != nil {
		return err
	}

	p := map[string]interface{}{
		"error":       "",
		"accessToken": "",
	}

	if req.UserName != "" {
		token, err := clientConf.PasswordCredentialsToken(r.Context(), req.UserName, req.Password)
		if err != nil {
			p["error"] = err.Error()
		} else {
			p["accessToken"] = token.AccessToken
			p["token"] = fmt.Sprintf("%#v", token)
		}
	}

	err = tmpl.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}

type CallbackRequest struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Revoke           string `json:"revoke"`
	Refresh          string `json:"refresh"`
	AccessToken      string `json:"access_token"`
	Code             string `json:"code"`
}

func callbackHandler(w http.ResponseWriter, r *http.Request, req *CallbackRequest) error {
	data, err := ioutil.ReadFile("./public/app/callback.html.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("callback").Parse(string(data))
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"error":            req.Error,
		"errorDescription": req.ErrorDescription,
		"code":             req.Code,
		"revoke":           "",
		"refresh":          "",
	}

	if req.Revoke != "" {
		params["revoke"] = req.Revoke

		revokeURL := strings.Replace(clientConf.Endpoint.TokenURL, "token", "revoke", 1)
		vs := url.Values{
			"client_id":       {clientConf.ClientID},
			"client_secret":   {clientConf.ClientSecret},
			"token_type_hint": {"refresh_token"},
			"token":           {req.Revoke},
		}

		hr, err := http.NewRequest("POST", revokeURL, strings.NewReader(vs.Encode()))
		if err != nil {
			return err
		}
		hr.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(hr)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		params["revokeStatusCode"] = resp.StatusCode
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		params["revokeBody"] = string(b)
		params["revokeRefreshTokenAfterProtectURL"] = "?refresh=" + url.QueryEscape(req.Revoke)
		params["revokeAccessTokenAfterProtectURL"] = "/protected?token=" + url.QueryEscape(req.AccessToken)
	}

	if req.Refresh != "" {

		vs := url.Values{
			"client_id":     {clientConf.ClientID},
			"client_secret": {clientConf.ClientSecret},
			"grant_type":    {"refresh_token"},
			"refresh_token": {req.Refresh},
			"scope":         {"fosite"},
		}

		hr, err := http.NewRequest("POST", clientConf.Endpoint.TokenURL, strings.NewReader(vs.Encode()))
		if err != nil {
			return err
		}
		hr.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(hr)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		params["refreshBody"] = string(b)

		v := make(map[string]interface{})
		err = json.Unmarshal(b, &v)
		if err != nil {
			return err
		}
		accessToken, ok := v["access_token"].(string)
		if !ok {
			accessToken = "😿"
		}
		refreshToken, ok := v["refresh_token"].(string)
		if !ok {
			refreshToken = "😿"
		}

		params["refresh"] = refreshToken
		params["protectedURL"] = "/protected?token=" + accessToken
		params["accessToken"] = accessToken
		params["useRefreshTokenURL"] = "?refresh=" + url.QueryEscape(refreshToken)
		params["revokeURL"] = "?revoke=" + url.QueryEscape(refreshToken) + "&access_token=" + url.QueryEscape(accessToken)
		params["refreshToken"] = refreshToken
		params["extraInfo"] = string(b)

	} else if req.Code != "" {
		token, err := clientConf.Exchange(r.Context(), req.Code)
		if err != nil {
			return err
		}
		params["protectedURL"] = "/protected?token=" + token.AccessToken
		params["accessToken"] = token.AccessToken
		params["useRefreshTokenURL"] = "?refresh=" + url.QueryEscape(token.RefreshToken)
		params["revokeURL"] = "?revoke=" + url.QueryEscape(token.RefreshToken) + "&access_token=" + url.QueryEscape(token.AccessToken)
		params["refreshToken"] = token.RefreshToken
		params["extraInfo"] = fmt.Sprintf("%#v", token)
	}

	err = tmpl.Execute(w, params)
	if err != nil {
		return err
	}

	return nil
}

func protectedHandler(w http.ResponseWriter, r *http.Request, req *CallbackRequest) error {
	data, err := ioutil.ReadFile("./public/app/protected.html.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("protected").Parse(string(data))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, map[string]interface{}{})
	if err != nil {
		return err
	}

	return nil
}
