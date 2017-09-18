package app

import (
	"gopkg.in/zabawaba99/firego.v1"
	"context"
	"io/ioutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/urlfetch"
	"github.com/favclip/ucon/swagger"
	"github.com/favclip/ucon"
	"google.golang.org/appengine/log"
)

func exampleSetup(swPlugin *swagger.Plugin) {
	s := &exampleService{}

	tag := swPlugin.AddTag(&swagger.Tag{Name: "Example"})
	var info *swagger.HandlerInfo

	info = swagger.NewHandlerInfo(s.Push)
	ucon.Handle("POST", "/api/example/push", info)
	info.Description, info.Tags = "exec Push", []string{tag.Name}

	info = swagger.NewHandlerInfo(s.Set)
	ucon.Handle("POST", "/api/example/set", info)
	info.Description, info.Tags = "exec Set", []string{tag.Name}
}

type exampleService struct{}

type Example struct {
	Name  string `json:"name"` // いわゆるID
	Price int `json:"price"`
}

func (s *exampleService) client(c context.Context) (*firego.Firebase, error) {
	d, err := ioutil.ReadFile("firebase-service-account.json")
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		return nil, err
	}

	client := conf.Client(c)
	t := client.Transport.(*oauth2.Transport)
	t.Base = urlfetch.Client(c).Transport

	f := firego.New("https://vv-firebase-playground.firebaseio.com", client)

	return f, nil
}

func (s *exampleService) Push(c context.Context) (*Noop, error) {
	f, err := s.client(c)
	if err != nil {
		log.Infof(c, "on s.Client: %#v", err)
		return nil, err
	}

	exampleRef, err := f.Ref("firebase-gaego")
	if err != nil {
		log.Infof(c, "on f.Ref: %#v", err)
		return nil, err
	}

	example := &Example{
		Name:  "foobar",
		Price: 111,
	}

	_, err = exampleRef.Push(example)
	if err != nil {
		log.Infof(c, "on exampleRef.Push: %#v", err)
		return nil, err
	}

	return &Noop{}, nil
}

func (s *exampleService) Set(c context.Context) (*Noop, error) {
	f, err := s.client(c)
	if err != nil {
		log.Infof(c, "on s.Client: %#v", err)
		return nil, err
	}

	exampleRef, err := f.Ref("firebase-gaego")
	if err != nil {
		log.Infof(c, "on f.Ref: %#v", err)
		return nil, err
	}

	example := &Example{
		Name:  "foobar",
		Price: 111,
	}

	// 相対パスではないらしい…
	entityRef, err := exampleRef.Ref("firebase-gaego/__id111111__")
	if err != nil {
		log.Infof(c, "on exampleRef.Ref: %#v", err)
		return nil, err
	}

	err = entityRef.Set(example)
	if err != nil {
		log.Infof(c, "on entityRef.Set: %#v", err)
		return nil, err
	}

	return &Noop{}, nil
}
