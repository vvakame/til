package idp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"go.mercari.io/datastore/dsmiddleware/dslog"

	cloudds "cloud.google.com/go/datastore"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/vvakame/til/go/oauth2idp-example/domains"
	"github.com/vvakame/til/go/oauth2idp-example/dsstorage"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
)

var baseURL string
var dsCli datastore.Client
var privateKey *rsa.PrivateKey

func init() {
	ctx := context.Background()

	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	log.Printf("DATASTORE_PROJECT_ID: %s", os.Getenv("DATASTORE_PROJECT_ID"))
	baseDsCli, err := cloudds.NewClient(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	dsCli, err = clouddatastore.FromClient(ctx, baseDsCli)
	if err != nil {
		log.Fatal(err)
	}
	dsCli.AppendMiddleware(dslog.NewLogger("datastore: ", func(ctx context.Context, format string, args ...interface{}) {
		log.Printf(format, args...)
	}))

	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
}

func ProvideDatastore() (datastore.Client, error) {
	return dsCli, nil
}

func ProvideStore(dsCli datastore.Client) (dsstorage.Storage, error) {
	store, err := dsstorage.NewStorage(&dsstorage.Config{
		DatastoreClient: func(ctx context.Context) (datastore.Client, error) {
			return dsCli, nil
		},
		AuthenticateUser: func(ctx context.Context, name, secret string) error {
			if name == "vvakame" && secret == "foobar" {
				return nil
			}
			return fosite.ErrNotFound
		},
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = store.CreateClient(ctx, &fosite.DefaultOpenIDConnectClient{
		DefaultClient: &fosite.DefaultClient{
			ID:            "my-client",
			Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
			RedirectURIs:  []string{baseURL + "/callback"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			ResponseTypes: []string{"id_token", "code", "token", "id_token token"}, // NOTE https://github.com/ory/fosite/issues/304
			Scopes:        []string{"fosite", "openid", "photos", "offline"},
		},
		// TokenEndpointAuthMethod: "client_secret_basic",
		TokenEndpointAuthMethod: "client_secret_post",
	})
	if err != nil {
		return nil, err
	}
	err = store.CreateClient(ctx, &fosite.DefaultClient{
		ID:            "encoded:client",
		Secret:        []byte(`$2a$10$A7M8b65dSSKGHF0H2sNkn.9Z0hT8U1Nv6OWPV3teUUaczXkVkxuDS`), // = "encoded&password"
		RedirectURIs:  []string{baseURL + "/callback"},
		ResponseTypes: []string{"id_token", "code", "token"},
		GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
		Scopes:        []string{"fosite", "openid", "photos", "offline"},
	})
	if err != nil {
		return nil, err
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

func ProvideOAuth2Provider(config *compose.Config, store dsstorage.Storage, strategy *compose.CommonStrategy) (fosite.OAuth2Provider, error) {
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
	// or oauth2.JWTSessionContainer
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
