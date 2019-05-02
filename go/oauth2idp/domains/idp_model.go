package domains

import (
	"context"
	"go.mercari.io/datastore"
	"time"
)

var DefaultAppInfo = &ApplicationInfo{
	Scopes: []string{
		"todo:read",
		"todo:write",
	},
}

type ApplicationInfo struct {
	Scopes []string
}

var _ datastore.PropertyLoadSaver = (*ClientInfo)(nil)

// ClientInfo provides OAuth2 client configurations.
// related: golang.org/x/oauth2.Config
//          golang.org/x/oauth2/clientcredentials.Config
type ClientInfo struct {
	ClientID     string `datastore:"-" boom:"id"`
	ClientSecret string `validate:"req"`
	Endpoint     string ``
	AuthURL      string ``
	TokenURL     string ``
	// AuthStyle      oauth2.AuthStyle ``
	// EndpointParams url.Values       ``
	RedirectURL string    ``
	Scopes      []string  ``
	UpdatedAt   time.Time ``
	CreatedAt   time.Time ``
}

func (cliInfo *ClientInfo) Load(ctx context.Context, ps []datastore.Property) error {
	return datastore.LoadStruct(ctx, cliInfo, ps)
}

func (cliInfo *ClientInfo) Save(ctx context.Context) ([]datastore.Property, error) {
	if cliInfo.CreatedAt.IsZero() {
		cliInfo.CreatedAt = time.Now()
	}
	cliInfo.UpdatedAt = time.Now()

	err := Validate(cliInfo)
	if err != nil {
		return nil, err
	}

	return datastore.SaveStruct(ctx, cliInfo)
}
