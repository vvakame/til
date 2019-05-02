package restapi

import (
	"fmt"
	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"golang.org/x/oauth2"
	"html/template"
	"io/ioutil"
	"net/http"
)

// OAuth2のクライアント側になるアプリの世界…

// アプリのClient アプリに対して固定値
var clientConf = oauth2.Config{
	ClientID:     "my-client",
	ClientSecret: "foobar",
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{"photos", "openid", "offline"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://localhost:8080/oauth2/auth",
		TokenURL: "http://localhost:8080/oauth2/token",
	},
}

func SetupRestAPI(swPlugin *swagger.Plugin) {
	ucon.HandleFunc("GET", "/", indexHandler)
	ucon.HandleFunc("GET", "/callback", callbackHandler)
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
		"authLinkURL": clientConf.AuthCodeURL("some-random-state-foobar") + "&nonce=some-random-nonce",
	})
	if err != nil {
		return err
	}

	return nil
}

type CallbackRequest struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
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

	token, err := clientConf.Exchange(r.Context(), req.Code)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"error":            req.Error,
		"errorDescription": req.ErrorDescription,
		"code":             req.Code,
		"protectedURL":     "/protected?token=" + token.AccessToken,
		"accessToken":      token.AccessToken,
		"extraInfo":        fmt.Sprintf("%#v", token),
	})
	if err != nil {
		return err
	}

	return nil
}
