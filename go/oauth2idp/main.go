package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/favclip/ucon"
	"github.com/vvakame/til/go/oauth2idp-example/app"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
	"github.com/vvakame/til/go/oauth2idp-example/idp"
)

func main() {
	log.Println("main: 👀")

	ucon.Middleware(UseUserDI)
	ucon.Orthodox()

	ucon.Middleware(func(b *ucon.Bubble) error {
		log.Printf("request url: %s %s", b.R.Method, b.R.URL.String())
		return b.Next()
	})

	idp.SetupIDP(ucon.DefaultMux)
	app.SetupAppAPI(ucon.DefaultMux)

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
