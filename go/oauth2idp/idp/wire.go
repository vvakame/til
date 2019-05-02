//+build wireinject

package idp

import (
	"context"
	"github.com/google/wire"
	"github.com/ory/fosite"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
)

func InitializeProvider() (fosite.OAuth2Provider, error) {
	wire.Build(ProvideStore, ProvideConfig, ProvideRSAPrivateKey, ProvideStrategy, ProvideOAuth2Provider)
	return nil, nil
}

func InitializeSession(ctx context.Context, user *domains.User) (Session, error) {
	wire.Build(ProvideSession)
	return nil, nil
}
