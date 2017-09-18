package app

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"

	"github.com/favclip/ucon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// UseAppengineContext replace Bubble.Context by appengine's context.Context.
// This middleware must append before the NetContextDI (it contains in ucon.Orthodox()).
func UseAppengineContext(b *ucon.Bubble) error {
	if b.Context == nil {
		b.Context = appengine.NewContext(b.R)
	} else {
		b.Context = appengine.WithContext(b.Context, b.R)
	}

	return b.Next()
}

var httpReqType = reflect.TypeOf((*http.Request)(nil))
var httpRespType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
var netContextType = reflect.TypeOf((*context.Context)(nil)).Elem()

// UseReqRespLogger logging request & response value.
func UseReqRespLogger(b *ucon.Bubble) error {
	re := regexp.MustCompile("[pP]assword\"")
	for idx, v := range b.Arguments {
		if !v.IsValid() {
			continue
		}
		if t := v.Type(); t == httpReqType || t == httpRespType {
			continue
		} else if t.AssignableTo(netContextType) {
			continue
		}

		entity, err := json.MarshalIndent(v.Interface(), "", "  ")
		if err != nil {
			return err
		}
		if re.FindString(string(entity)) != "" {
			log.Infof(b.Context, "req %d - **censored**", idx+1)
			continue
		}

		log.Infof(b.Context, "req %d - %s", idx+1, string(entity))
	}

	err := b.Next()
	if err != nil {
		log.Infof(b.Context, "resp go level error: %#v", err)
		return err
	}

	for idx, v := range b.Returns {
		if !v.IsValid() {
			continue
		}

		entity, err := json.MarshalIndent(v.Interface(), "", "  ")
		if err != nil {
			return err
		}
		log.Infof(b.Context, "resp %d - %s", idx+1, string(entity))
	}

	return nil
}
