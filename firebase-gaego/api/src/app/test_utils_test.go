package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/favclip/ucon/swagger"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/user"
)

func setupHTTPRequest(t *testing.T, inst aetest.Instance, method, path string, req interface{}) *http.Request {
	if req == nil {
		r, err := inst.NewRequest(method, path, nil)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	refV := reflect.ValueOf(req)
	if refV.Kind() == reflect.Ptr {
		refV = refV.Elem()
	}

	if refV.Kind() != reflect.Struct {
		t.Fatalf("req should be struct. or do not use this helper function.")
	}

	reqURL, err := url.Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	vs := url.Values{}
	var processQueryValue func(refV reflect.Value)
	processQueryValue = func(refV reflect.Value) {
		for i, fLen := 0, refV.NumField(); i < fLen; i++ {
			sf := refV.Type().Field(i)
			if sf.Anonymous {
				processQueryValue(refV.Field(i))
				continue
			}
			swaggerTag := swagger.NewTagSwagger(sf.Tag)
			if swaggerTag.Private() {
				continue
			}
			if swaggerTag.In() == "query" {
				v := fmt.Sprintf("%v", refV.Field(i).Interface())
				if v != "" {
					vs.Add(swaggerTag.Name(), v)
				}
			}
		}
	}
	processQueryValue(refV)
	reqURL.RawQuery = vs.Encode()

	var body io.Reader
	if method == "POST" || method == "PUT" {
		b, err := json.MarshalIndent(req, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		body = bytes.NewReader(b)
	}

	r, err := inst.NewRequest(method, reqURL.String(), body)
	if err != nil {
		t.Fatal(err)
	}

	if body != nil {
		r.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	return r
}

func makeTestDataUser(t *testing.T, usr *user.User) *user.User {
	if usr == nil {
		usr = &user.User{}
	}

	if usr.ID == "" {
		usr.ID = randomdata.SillyName()
	}
	if usr.Email == "" {
		usr.Email = randomdata.Email()
	}

	return usr
}
