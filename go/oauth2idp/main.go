package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/ory/fosite"
	"github.com/vvakame/til/go/oauth2idp-example/app"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
	"github.com/vvakame/til/go/oauth2idp-example/idp"
)

var _ ucon.HTTPErrorResponse = (*fositeError)(nil)

type fositeError struct {
	Base *fosite.RFC6749Error
}

func (fe *fositeError) StatusCode() int {
	return fe.Base.Code
}

func (fe *fositeError) ErrorMessage() interface{} {
	return fe.Base
}

func main() {
	log.Println("main: 👀")

	ucon.Middleware(UseUserDI)
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

	ucon.Middleware(func(b *ucon.Bubble) error {
		log.Printf("request url: %s %s", b.R.Method, b.R.URL.String())
		return b.Next()
	})

	swPlugin := swagger.NewPlugin(&swagger.Options{
		Object: &swagger.Object{
			Info: &swagger.Info{
				Title:   "OAuth2 IDP",
				Version: "v1",
			},
		},
		DefinitionNameModifier: func(refT reflect.Type, defName string) string {
			if strings.HasSuffix(defName, "JSON") {
				return defName[:len(defName)-4]
			}
			return defName
		},
	})
	ucon.Plugin(swPlugin)

	idp.SetupIDP(swPlugin)
	app.SetupAppAPI(swPlugin)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("listen: %s", addr)

	if err := ucon.DefaultMux.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}

var userType = reflect.TypeOf((*domains.User)(nil))

func UseUserDI(b *ucon.Bubble) error {
	// 真面目にログイン処理実装するのがめんどくさすぎるので固定ユーザでログインしていることにする

	for idx, argT := range b.ArgumentTypes {
		if argT == userType {
			user := &domains.User{
				ID:          100,
				Name:        "vvakame",
				NewPassword: "foobar",
			}
			err := user.EncryptIfNeeded()
			if err != nil {
				return err
			}

			b.Arguments[idx] = reflect.ValueOf(user)
			continue
		}
	}

	return b.Next()
}
