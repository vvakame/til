package idp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
	"go.mercari.io/datastore"
	"time"
)

var privateKey *rsa.PrivateKey

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
}

func ProvideDatastore() (datastore.Client, error) {
	panic("test")
}

type Storage interface {
	fosite.Storage
	// oauth2.CoreStorage
	oauth2.AuthorizeCodeStorage
	oauth2.AccessTokenStorage
	oauth2.RefreshTokenStorage
	// +oauth2.TokenRevocationStorage
	RevokeRefreshToken(ctx context.Context, requestID string) error
	RevokeAccessToken(ctx context.Context, requestID string) error
	// and...
	openid.OpenIDConnectRequestStorage
}

func ProvideStore() (Storage, error) {
	store := &storage.MemoryStore{
		IDSessions: make(map[string]fosite.Requester),
		Clients: map[string]fosite.Client{
			"my-client": &fosite.DefaultClient{
				ID:            "my-client",
				Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
				RedirectURIs:  []string{"http://localhost:8080/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scopes:        []string{"fosite", "openid", "photos", "offline"},
			},
			"encoded:client": &fosite.DefaultClient{
				ID:            "encoded:client",
				Secret:        []byte(`$2a$10$A7M8b65dSSKGHF0H2sNkn.9Z0hT8U1Nv6OWPV3teUUaczXkVkxuDS`), // = "encoded&password"
				RedirectURIs:  []string{"http://localhost:8080/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scopes:        []string{"fosite", "openid", "photos", "offline"},
			},
		},
		Users: map[string]storage.MemoryUserRelation{
			// see UseUserDI
			"100": {
				Username: "vvakame",
				Password: "foobar",
			},
		},
		AuthorizeCodes:         map[string]storage.StoreAuthorizeCode{},
		Implicit:               map[string]fosite.Requester{},
		AccessTokens:           map[string]fosite.Requester{},
		RefreshTokens:          map[string]fosite.Requester{},
		PKCES:                  map[string]fosite.Requester{},
		AccessTokenRequestIDs:  map[string]string{},
		RefreshTokenRequestIDs: map[string]string{},
	}

	return store, nil
}

func ProvideConfig() (*compose.Config, error) {
	return &compose.Config{}, nil
}

func ProvideRSAPrivateKey() (*rsa.PrivateKey, error) {
	return privateKey, nil
}

func ProvideStrategy(config *compose.Config, privateKey *rsa.PrivateKey) (*compose.CommonStrategy, error) {
	secret := []byte("test-test-test-test-test-test-test-test-test-test-test-test")
	strategy := &compose.CommonStrategy{
		CoreStrategy:               compose.NewOAuth2HMACStrategy(config, secret, nil),
		OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(config, privateKey),
	}
	return strategy, nil
}

func ProvideOAuth2Provider(config *compose.Config, store Storage, strategy *compose.CommonStrategy) (fosite.OAuth2Provider, error) {
	provider := compose.Compose(
		config,
		store,
		strategy,
		nil,

		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2ResourceOwnerPasswordCredentialsFactory,

		compose.OAuth2TokenRevocationFactory,
		compose.OAuth2TokenIntrospectionFactory,

		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectImplicitFactory,
		compose.OpenIDConnectHybridFactory,
		compose.OpenIDConnectRefreshFactory,
	)
	return provider, nil
}

type Session interface {
	openid.Session
}

func ProvideSession(ctx context.Context, user *domains.User) (Session, error) {
	subject := ""
	if user != nil {
		subject = fmt.Sprintf("%d", user.ID)
	}
	session := &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "https://fosite.my-application.com",
			Subject:     subject,
			Audience:    []string{"https://my-client.my-application.com"},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}

	return session, nil
}
