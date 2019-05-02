package main

import (
	"context"
	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
	"github.com/vvakame/til/go/oauth2idp-example/idp"
	"github.com/vvakame/til/go/oauth2idp-example/restapi"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
	"reflect"
	"strings"
)

var dsCli datastore.Client

func main() {
	var err error
	dsCli, err = clouddatastore.FromContext(context.Background())
	if err != nil {
		panic(err)
	}

	ucon.Middleware(UseUserDI)
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

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

	restapi.SetupRestAPI(swPlugin)
	idp.SetupIDP(swPlugin)

	err = ucon.DefaultMux.ListenAndServe(":8080")
	if err != nil {
		panic(err)
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
