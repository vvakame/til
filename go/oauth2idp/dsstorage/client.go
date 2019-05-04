package dsstorage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ory/fosite"
	"go.mercari.io/datastore"
	"gopkg.in/square/go-jose.v2"
)

var _ fosite.Client = (*DefaultClient)(nil)
var _ fosite.OpenIDConnectClient = (*DefaultClient)(nil)
var _ datastore.KeyLoader = (*DefaultClient)(nil)
var _ datastore.PropertyLoadSaver = (*DefaultClient)(nil)

type DefaultClient struct {
	// for fosite.Client
	ID            string   `datastore:"-" boom:"id"`
	Secret        []byte   `datastore:",noindex"`
	RedirectURIs  []string ``
	GrantTypes    []string ``
	ResponseTypes []string ``
	Scopes        []string ``
	Audience      []string ``
	Public        bool     ``
	// for fosite.OpenIDConnectClient
	JSONWebKeysURI                string              ``
	JSONWebKeysJSON               string              `json:"-" datastore:",noindex"`
	JSONWebKeys                   *jose.JSONWebKeySet `datastore:"-"`
	TokenEndpointAuthMethod       string              ``
	RequestURIs                   []string            ``
	RequestObjectSigningAlgorithm string              ``
	// others...
	UpdatedAt time.Time ``
	CreatedAt time.Time ``
}

func (cli *DefaultClient) LoadKey(ctx context.Context, key datastore.Key) error {
	cli.ID = key.Name()
	return nil
}

func (cli *DefaultClient) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, cli, ps)
	if err != nil {
		return err
	}

	if cli.JSONWebKeysJSON != "" {
		var jwks jose.JSONWebKeySet
		err = json.Unmarshal([]byte(cli.JSONWebKeysJSON), &jwks)
		if err != nil {
			return err
		}
		cli.JSONWebKeys = &jwks
	}

	return nil
}

func (cli *DefaultClient) Save(ctx context.Context) ([]datastore.Property, error) {
	if cli.CreatedAt.IsZero() {
		cli.CreatedAt = time.Now()
	}
	cli.UpdatedAt = time.Now()

	if cli.JSONWebKeys != nil {
		b, err := json.Marshal(cli.JSONWebKeys)
		if err != nil {
			return nil, err
		}
		cli.JSONWebKeysJSON = string(b)
	} else {
		cli.JSONWebKeysJSON = ""
	}

	return datastore.SaveStruct(ctx, cli)
}

func (cli *DefaultClient) GetID() string {
	return cli.ID
}

func (cli *DefaultClient) GetHashedSecret() []byte {
	return cli.Secret
}

func (cli *DefaultClient) GetRedirectURIs() []string {
	return cli.RedirectURIs
}

func (cli *DefaultClient) GetGrantTypes() fosite.Arguments {
	if len(cli.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}
	return cli.GrantTypes
}

func (cli *DefaultClient) GetResponseTypes() fosite.Arguments {
	if len(cli.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}
	return cli.ResponseTypes
}

func (cli *DefaultClient) GetScopes() fosite.Arguments {
	return cli.Scopes
}

func (cli *DefaultClient) IsPublic() bool {
	return cli.Public
}

func (cli *DefaultClient) GetAudience() fosite.Arguments {
	return cli.Audience
}

func (cli *DefaultClient) GetRequestURIs() []string {
	return cli.RequestURIs
}

func (cli *DefaultClient) GetJSONWebKeys() *jose.JSONWebKeySet {
	return nil
}

func (cli *DefaultClient) GetJSONWebKeysURI() string {
	return cli.JSONWebKeysURI
}

func (cli *DefaultClient) GetRequestObjectSigningAlgorithm() string {
	return cli.RequestObjectSigningAlgorithm
}

func (cli *DefaultClient) GetTokenEndpointAuthMethod() string {
	return cli.TokenEndpointAuthMethod
}

func (cli *DefaultClient) GetTokenEndpointAuthSigningAlgorithm() string {
	return "RS256"
}
