package dsstorage

import (
	"context"
	"encoding/json"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/pkce"
	"github.com/ory/fosite/storage"
	"github.com/pkg/errors"
	"go.mercari.io/datastore"
	"golang.org/x/xerrors"
)

var _ Storage = (*datastoreStorage)(nil)
var _ oauth2.CoreStorage = (*datastoreStorage)(nil)
var _ oauth2.TokenRevocationStorage = (*datastoreStorage)(nil)
var _ oauth2.ResourceOwnerPasswordCredentialsGrantStorage = (*datastoreStorage)(nil)
var _ oauth2.CoreStorage = (*datastoreStorage)(nil)

// Storage provides fosite storage by Google Cloud Datastore or AppEngine Datastore.
type Storage interface {
	fosite.Storage
	// for oauth2.CoreStorage
	oauth2.AuthorizeCodeStorage
	oauth2.AccessTokenStorage
	oauth2.RefreshTokenStorage
	// for oauth2.TokenRevocationStorage
	RevokeRefreshToken(ctx context.Context, requestID string) error
	RevokeAccessToken(ctx context.Context, requestID string) error
	// for oauth2.ResourceOwnerPasswordCredentialsGrantStorage
	Authenticate(ctx context.Context, name string, secret string) error
	// and...
	openid.OpenIDConnectRequestStorage
	storage.Transactional
	pkce.PKCERequestStorage

	// original
	CreateClient(ctx context.Context, client fosite.Client) error
}

// Config provides some settings.
type Config struct {
	DatastoreClient func(context.Context) (datastore.Client, error)

	NewClientEntity func() fosite.Client
	NewRequester    func() fosite.Requester
	NewSession      func() fosite.Session

	AuthenticateUser func(ctx context.Context, name, secret string) error

	ClientKind        string
	AuthorizeCodeKind string
	IDSessionKind     string
	AccessTokenKind   string
	RefreshTokenKind  string
	PKCEKind          string
}

func NewStorage(config *Config) (Storage, error) {
	if config == nil {
		config = &Config{}
	}

	dsStorage := &datastoreStorage{}

	if config.DatastoreClient == nil {
		return nil, errors.New("DatastoreClient is required")
	}
	dsStorage.datastoreClient = config.DatastoreClient

	if config.NewClientEntity != nil {
		dsStorage.newClientEntity = config.NewClientEntity
	} else {
		dsStorage.newClientEntity = func() fosite.Client {
			return &DefaultClient{}
		}
	}
	if config.NewRequester != nil {
		dsStorage.newRequester = config.NewRequester
	} else {
		dsStorage.newRequester = func() fosite.Requester {
			return &DefaultRequester{}
		}
	}
	if config.NewSession != nil {
		dsStorage.newSession = config.NewSession
	} else {
		dsStorage.newSession = func() fosite.Session {
			return &openid.DefaultSession{}
		}
	}
	if config.AuthenticateUser != nil {
		dsStorage.authenticateUser = config.AuthenticateUser
	} else {
		dsStorage.authenticateUser = func(ctx context.Context, name, secret string) error {
			return errors.New("invalid credentials")
		}
	}

	if config.ClientKind != "" {
		dsStorage.ClientKind = config.ClientKind
	} else {
		dsStorage.ClientKind = "FositeClient"
	}
	if config.AuthorizeCodeKind != "" {
		dsStorage.AuthorizeCodeKind = config.AuthorizeCodeKind
	} else {
		dsStorage.AuthorizeCodeKind = "FositeAuthorizeCode"
	}
	if config.IDSessionKind != "" {
		dsStorage.IDSessionKind = config.IDSessionKind
	} else {
		dsStorage.IDSessionKind = "FositeIDSession"
	}
	if config.AccessTokenKind != "" {
		dsStorage.AccessTokenKind = config.AccessTokenKind
	} else {
		dsStorage.AccessTokenKind = "FositeAccessToken"
	}
	if config.RefreshTokenKind != "" {
		dsStorage.RefreshTokenKind = config.RefreshTokenKind
	} else {
		dsStorage.RefreshTokenKind = "FositeRefreshToken"
	}
	if config.PKCEKind != "" {
		dsStorage.PKCEKind = config.PKCEKind
	} else {
		dsStorage.PKCEKind = "FositePKCE"
	}

	return dsStorage, nil
}

type datastoreStorage struct {
	datastoreClient  func(context.Context) (datastore.Client, error)
	newClientEntity  func() fosite.Client
	newRequester     func() fosite.Requester
	newSession       func() fosite.Session
	authenticateUser func(ctx context.Context, name, secret string) error

	ClientKind        string
	AuthorizeCodeKind string
	IDSessionKind     string
	AccessTokenKind   string
	RefreshTokenKind  string
	PKCEKind          string
}

type contextTxKey struct{}

func (s *datastoreStorage) BeginTX(ctx context.Context) (context.Context, error) {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := dsCli.NewTransaction(ctx)
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, contextTxKey{}, tx)
	return ctx, nil
}

func (s *datastoreStorage) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if !ok {
		return errors.WithStack(ErrInvalidTxContext)
	}
	_, err := tx.Commit()
	return err
}

func (s *datastoreStorage) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if !ok {
		return errors.WithStack(ErrInvalidTxContext)
	}
	return tx.Rollback()
}

func (s *datastoreStorage) CreateClient(ctx context.Context, client fosite.Client) error {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return err
	}
	put := func(key datastore.Key, src interface{}) error {
		_, err := dsCli.Put(ctx, key, src)
		return err
	}
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if ok {
		put = func(key datastore.Key, src interface{}) error {
			_, err := tx.Put(key, src)
			return err
		}
	}

	switch v := client.(type) {
	case *fosite.DefaultClient:
		cliEntity := &DefaultClient{}

		cliEntity.ID = v.GetID()
		cliEntity.Secret = v.GetHashedSecret()
		cliEntity.RedirectURIs = v.GetRedirectURIs()
		cliEntity.GrantTypes = v.GetGrantTypes()
		cliEntity.ResponseTypes = v.GetResponseTypes()
		cliEntity.Scopes = v.GetScopes()
		cliEntity.Audience = v.GetAudience()
		cliEntity.Public = v.IsPublic()

		key := dsCli.NameKey(s.ClientKind, client.GetID(), nil)
		err := put(key, cliEntity)
		if err != nil {
			return err
		}

	case *fosite.DefaultOpenIDConnectClient:
		cliEntity := &DefaultClient{}

		cliEntity.ID = v.GetID()
		cliEntity.Secret = v.GetHashedSecret()
		cliEntity.RedirectURIs = v.GetRedirectURIs()
		cliEntity.GrantTypes = v.GetGrantTypes()
		cliEntity.ResponseTypes = v.GetResponseTypes()
		cliEntity.Scopes = v.GetScopes()
		cliEntity.Audience = v.GetAudience()
		cliEntity.Public = v.IsPublic()

		cliEntity.JSONWebKeysURI = v.GetJSONWebKeysURI()
		cliEntity.JSONWebKeys = v.GetJSONWebKeys()
		cliEntity.TokenEndpointAuthMethod = v.GetTokenEndpointAuthMethod()
		cliEntity.RequestURIs = v.GetRequestURIs()
		cliEntity.RequestObjectSigningAlgorithm = v.GetRequestObjectSigningAlgorithm()

		key := dsCli.NameKey(s.ClientKind, client.GetID(), nil)
		err := put(key, cliEntity)
		if err != nil {
			return err
		}

	case datastore.PropertyLoadSaver:
		key := dsCli.NameKey(s.ClientKind, client.GetID(), nil)
		err := put(key, client)
		if err != nil {
			return err
		}

	default:
		return ErrUnsupportedType
	}

	return nil
}

func (s *datastoreStorage) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return nil, err
	}
	get := func(key datastore.Key, src interface{}) error {
		return dsCli.Get(ctx, key, src)
	}
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if ok {
		get = func(key datastore.Key, src interface{}) error {
			return tx.Get(key, src)
		}
	}

	client := s.newClientEntity()

	switch v := client.(type) {
	case *fosite.DefaultClient:
		cliEntity := &DefaultClient{}
		key := dsCli.NameKey(s.ClientKind, id, nil)
		err := get(key, cliEntity)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		v.ID = cliEntity.GetID()
		v.Secret = cliEntity.GetHashedSecret()
		v.RedirectURIs = cliEntity.GetRedirectURIs()
		v.GrantTypes = cliEntity.GetGrantTypes()
		v.ResponseTypes = cliEntity.GetResponseTypes()
		v.Scopes = cliEntity.GetScopes()
		v.Audience = cliEntity.GetAudience()
		v.Public = cliEntity.IsPublic()

		return client, nil

	case *fosite.DefaultOpenIDConnectClient:
		cliEntity := &DefaultClient{}
		key := dsCli.NameKey(s.ClientKind, id, nil)
		err := get(key, cliEntity)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		v.ID = cliEntity.GetID()
		v.Secret = cliEntity.GetHashedSecret()
		v.RedirectURIs = cliEntity.GetRedirectURIs()
		v.GrantTypes = cliEntity.GetGrantTypes()
		v.ResponseTypes = cliEntity.GetResponseTypes()
		v.Scopes = cliEntity.GetScopes()
		v.Audience = cliEntity.GetAudience()
		v.Public = cliEntity.IsPublic()

		v.JSONWebKeysURI = cliEntity.GetJSONWebKeysURI()
		v.JSONWebKeys = cliEntity.GetJSONWebKeys()
		v.TokenEndpointAuthMethod = cliEntity.GetTokenEndpointAuthMethod()
		v.RequestURIs = cliEntity.GetRequestURIs()
		v.RequestObjectSigningAlgorithm = cliEntity.GetRequestObjectSigningAlgorithm()

		return client, nil

	case datastore.PropertyLoadSaver:
		key := dsCli.NameKey(s.ClientKind, id, nil)
		err := get(key, v)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		return client, nil

	default:
		return nil, ErrUnsupportedType
	}
}

func (s *datastoreStorage) putRequestEntity(ctx context.Context, kind string, id string, request fosite.Requester, prePut func(request fosite.Requester) error) error {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return err
	}
	put := func(key datastore.Key, src interface{}) error {
		_, err := dsCli.Put(ctx, key, src)
		return err
	}
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if ok {
		put = func(key datastore.Key, src interface{}) error {
			_, err := tx.Put(key, src)
			return err
		}
	}

	switch v := request.(type) {
	case *fosite.Request:
		reqEntity := &DefaultRequester{}

		reqEntity.ID = v.GetID()
		reqEntity.RequestedAt = v.GetRequestedAt()
		reqEntity.Client = v.GetClient()
		reqEntity.RequestedScope = v.GetRequestedScopes()
		reqEntity.GrantedScope = v.GetGrantedScopes()
		reqEntity.Form = v.GetRequestForm()
		reqEntity.Session = v.GetSession()
		reqEntity.RequestedAudience = v.GetRequestedAudience()
		reqEntity.GrantedAudience = v.GetGrantedAudience()

		key := dsCli.NameKey(kind, id, nil)
		if prePut != nil {
			err := prePut(reqEntity)
			if err != nil {
				return err
			}
		}
		err = put(key, reqEntity)
		if err != nil {
			return err
		}

	case *fosite.AccessRequest:
		reqEntity := &DefaultRequester{}

		reqEntity.ID = v.GetID()
		reqEntity.RequestedAt = v.GetRequestedAt()
		reqEntity.Client = v.GetClient()
		reqEntity.RequestedScope = v.GetRequestedScopes()
		reqEntity.GrantedScope = v.GetGrantedScopes()
		reqEntity.Form = v.GetRequestForm()
		reqEntity.Session = v.GetSession()
		reqEntity.RequestedAudience = v.GetRequestedAudience()
		reqEntity.GrantedAudience = v.GetGrantedAudience()

		reqEntity.GrantTypes = v.GetGrantTypes()
		reqEntity.HandledGrantType = v.HandledGrantType

		key := dsCli.NameKey(kind, id, nil)
		if prePut != nil {
			err := prePut(reqEntity)
			if err != nil {
				return err
			}
		}
		err = put(key, reqEntity)
		if err != nil {
			return err
		}

	case *fosite.AuthorizeRequest:
		reqEntity := &DefaultRequester{}

		reqEntity.ID = v.GetID()
		reqEntity.RequestedAt = v.GetRequestedAt()
		reqEntity.Client = v.GetClient()
		reqEntity.RequestedScope = v.GetRequestedScopes()
		reqEntity.GrantedScope = v.GetGrantedScopes()
		reqEntity.Form = v.GetRequestForm()
		reqEntity.Session = v.GetSession()
		reqEntity.RequestedAudience = v.GetRequestedAudience()
		reqEntity.GrantedAudience = v.GetGrantedAudience()

		reqEntity.ResponseTypes = v.GetResponseTypes()
		if v.GetRedirectURI() != nil {
			reqEntity.RedirectURI = v.GetRedirectURI().String()
		}
		reqEntity.State = v.GetState()
		reqEntity.HandledResponseTypes = v.HandledResponseTypes

		key := dsCli.NameKey(kind, id, nil)
		if prePut != nil {
			err := prePut(reqEntity)
			if err != nil {
				return err
			}
		}
		err = put(key, reqEntity)
		if err != nil {
			return err
		}

	case datastore.PropertyLoadSaver:
		key := dsCli.NameKey(kind, id, nil)
		if prePut != nil {
			err := prePut(request)
			if err != nil {
				return err
			}
		}
		err = put(key, request)
		if err != nil {
			return err
		}

	default:
		return ErrUnsupportedType
	}

	return nil
}

func (s *datastoreStorage) getRequestEntity(ctx context.Context, kind string, id string) (fosite.Requester, error) {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return nil, err
	}
	get := func(key datastore.Key, src interface{}) error {
		return dsCli.Get(ctx, key, src)
	}
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if ok {
		get = func(key datastore.Key, src interface{}) error {
			return tx.Get(key, src)
		}
	}

	request := s.newRequester()

	switch v := request.(type) {
	case *fosite.Request:
		reqEntity := &DefaultRequester{}
		key := dsCli.NameKey(kind, id, nil)
		err := get(key, reqEntity)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		client, err := s.GetClient(ctx, reqEntity.ClientID)
		if err != nil {
			return nil, err
		}
		reqEntity.Client = client

		if reqEntity.SessionJSON != "" {
			session := s.newSession()
			err = json.Unmarshal([]byte(reqEntity.SessionJSON), session)
			if err != nil {
				return nil, err
			}
			reqEntity.Session = session
		}

		v.ID = reqEntity.GetID()
		v.RequestedAt = reqEntity.GetRequestedAt()
		v.Client = reqEntity.GetClient()
		v.RequestedScope = reqEntity.GetRequestedScopes()
		v.GrantedScope = reqEntity.GetGrantedScopes()
		v.Form = reqEntity.GetRequestForm()
		v.Session = reqEntity.GetSession()
		v.RequestedAudience = reqEntity.GetRequestedAudience()
		v.GrantedAudience = reqEntity.GetGrantedAudience()

		if !reqEntity.Active {
			return v, fosite.ErrInvalidatedAuthorizeCode
		}
		return v, nil

	case *fosite.AccessRequest:
		reqEntity := &DefaultRequester{}
		key := dsCli.NameKey(kind, id, nil)
		err := get(key, reqEntity)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		client, err := s.GetClient(ctx, reqEntity.ClientID)
		if err != nil {
			return nil, err
		}
		reqEntity.Client = client

		if reqEntity.SessionJSON != "" {
			session := s.newSession()
			err = json.Unmarshal([]byte(reqEntity.SessionJSON), session)
			if err != nil {
				return nil, err
			}
			reqEntity.Session = session
		}

		v.ID = reqEntity.GetID()
		v.RequestedAt = reqEntity.GetRequestedAt()
		v.Client = reqEntity.GetClient()
		v.RequestedScope = reqEntity.GetRequestedScopes()
		v.GrantedScope = reqEntity.GetGrantedScopes()
		v.Form = reqEntity.GetRequestForm()
		v.Session = reqEntity.GetSession()
		v.RequestedAudience = reqEntity.GetRequestedAudience()
		v.GrantedAudience = reqEntity.GetGrantedAudience()

		v.GrantTypes = reqEntity.GetGrantTypes()
		v.HandledGrantType = reqEntity.HandledGrantType

		if !reqEntity.Active {
			return v, fosite.ErrInvalidatedAuthorizeCode
		}
		return v, nil

	case *fosite.AuthorizeRequest:
		reqEntity := &DefaultRequester{}
		key := dsCli.NameKey(kind, id, nil)
		err := tx.Get(key, reqEntity)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		client, err := s.GetClient(ctx, reqEntity.ClientID)
		if err != nil {
			return nil, err
		}
		reqEntity.Client = client

		if reqEntity.SessionJSON != "" {
			session := s.newSession()
			err = json.Unmarshal([]byte(reqEntity.SessionJSON), session)
			if err != nil {
				return nil, err
			}
			reqEntity.Session = session
		}

		v.ID = reqEntity.GetID()
		v.RequestedAt = reqEntity.GetRequestedAt()
		v.Client = reqEntity.GetClient()
		v.RequestedScope = reqEntity.GetRequestedScopes()
		v.GrantedScope = reqEntity.GetGrantedScopes()
		v.Form = reqEntity.GetRequestForm()
		v.Session = reqEntity.GetSession()
		v.RequestedAudience = reqEntity.GetRequestedAudience()
		v.GrantedAudience = reqEntity.GetGrantedAudience()

		v.ResponseTypes = reqEntity.GetResponseTypes()
		v.RedirectURI = reqEntity.GetRedirectURI()
		v.State = reqEntity.GetState()
		v.HandledResponseTypes = reqEntity.HandledResponseTypes

		if !reqEntity.Active {
			return v, fosite.ErrInvalidatedAuthorizeCode
		}
		return v, nil

	case datastore.PropertyLoadSaver:
		key := dsCli.NameKey(kind, id, nil)
		err := get(key, v)
		if xerrors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, ErrNoSuchEntity
		} else if err != nil {
			return nil, err
		}

		invalidator, ok := v.(ActiveStateModifier)
		if !ok {
			return nil, ErrUnsupportedType
		}
		if !invalidator.IsActive() {
			return request, fosite.ErrInvalidatedAuthorizeCode
		}

		if request.GetClient() == nil {
			clientLoader, ok := v.(ClientLoader)
			if !ok {
				return nil, ErrUnsupportedType
			}
			if clientLoader.GetClientID() != "" {
				client, err := s.GetClient(ctx, clientLoader.GetClientID())
				if err != nil {
					return nil, err
				}
				clientLoader.SetClient(client)
			}
		}
		if sessionLoader, ok := v.(SessionRestorer); ok {
			session := s.newSession()
			err := sessionLoader.RestoreSession(ctx, session)
			if err != nil {
				return nil, err
			}
		}

		return request, nil

	default:
		return nil, ErrUnsupportedType
	}
}

func (s *datastoreStorage) deleteRequestEntity(ctx context.Context, kind string, id string) error {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return err
	}
	del := func(key datastore.Key) error {
		return dsCli.Delete(ctx, key)
	}
	tx, ok := ctx.Value(contextTxKey{}).(datastore.Transaction)
	if ok {
		del = func(key datastore.Key) error {
			return tx.Delete(key)
		}
	}

	key := dsCli.NameKey(kind, id, nil)
	err = del(key)
	if xerrors.Is(err, datastore.ErrNoSuchEntity) {
		return ErrNoSuchEntity
	} else if err != nil {
		return err
	}
	return nil
}

func (s *datastoreStorage) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) error {
	return s.putRequestEntity(ctx, s.AuthorizeCodeKind, code, request, func(request fosite.Requester) error {
		invalidator, ok := request.(ActiveStateModifier)
		if !ok {
			return ErrUnsupportedType
		}
		invalidator.SetActive(true)
		return nil
	})
}

func (s *datastoreStorage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (fosite.Requester, error) {
	return s.getRequestEntity(ctx, s.AuthorizeCodeKind, code)
}

func (s *datastoreStorage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	request, err := s.getRequestEntity(ctx, s.AuthorizeCodeKind, code)
	if err != nil {
		return err
	}
	return s.putRequestEntity(ctx, s.AuthorizeCodeKind, code, request, func(request fosite.Requester) error {
		invalidator, ok := request.(ActiveStateModifier)
		if !ok {
			return ErrUnsupportedType
		}
		invalidator.SetActive(false)
		return nil
	})
}

func (s *datastoreStorage) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return s.putRequestEntity(ctx, s.AccessTokenKind, signature, request, func(request fosite.Requester) error {
		invalidator, ok := request.(ActiveStateModifier)
		if !ok {
			return ErrUnsupportedType
		}
		invalidator.SetActive(true)
		return nil
	})
}

func (s *datastoreStorage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return s.getRequestEntity(ctx, s.AccessTokenKind, signature)
}

func (s *datastoreStorage) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	return s.deleteRequestEntity(ctx, s.AccessTokenKind, signature)
}

func (s *datastoreStorage) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	return s.putRequestEntity(ctx, s.RefreshTokenKind, signature, request, func(request fosite.Requester) error {
		invalidator, ok := request.(ActiveStateModifier)
		if !ok {
			return ErrUnsupportedType
		}
		invalidator.SetActive(true)
		return nil
	})
}

func (s *datastoreStorage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	return s.getRequestEntity(ctx, s.RefreshTokenKind, signature)
}

func (s *datastoreStorage) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return s.deleteRequestEntity(ctx, s.RefreshTokenKind, signature)
}

func (s *datastoreStorage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return err
	}
	q := dsCli.NewQuery(s.RefreshTokenKind).Filter("ID =", requestID).KeysOnly().Limit(1)
	keys, err := dsCli.GetAll(ctx, q, nil)
	if err != nil {
		return err
	}
	if len(keys) != 0 {
		_ = s.DeleteRefreshTokenSession(ctx, keys[0].Name())
		_ = s.DeleteAccessTokenSession(ctx, keys[0].Name())
	}
	return nil
}

func (s *datastoreStorage) RevokeAccessToken(ctx context.Context, requestID string) error {
	dsCli, err := s.datastoreClient(ctx)
	if err != nil {
		return err
	}
	q := dsCli.NewQuery(s.AccessTokenKind).Filter("ID =", requestID).KeysOnly().Limit(1)
	keys, err := dsCli.GetAll(ctx, q, nil)
	if err != nil {
		return err
	}
	if len(keys) != 0 {
		_ = s.DeleteAccessTokenSession(ctx, keys[0].Name())
	}
	return nil
}

func (s *datastoreStorage) Authenticate(ctx context.Context, name string, secret string) error {
	return s.authenticateUser(ctx, name, secret)
}

func (s *datastoreStorage) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, request fosite.Requester) error {
	return s.putRequestEntity(ctx, s.IDSessionKind, authorizeCode, request, func(request fosite.Requester) error {
		invalidator, ok := request.(ActiveStateModifier)
		if !ok {
			return ErrUnsupportedType
		}
		invalidator.SetActive(true)
		return nil
	})

}

func (s *datastoreStorage) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	return s.getRequestEntity(ctx, s.RefreshTokenKind, authorizeCode)
}

func (s *datastoreStorage) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	return s.deleteRequestEntity(ctx, s.RefreshTokenKind, authorizeCode)
}

func (s *datastoreStorage) CreatePKCERequestSession(ctx context.Context, signature string, request fosite.Requester) error {
	return s.putRequestEntity(ctx, s.PKCEKind, signature, request, func(request fosite.Requester) error {
		invalidator, ok := request.(ActiveStateModifier)
		if !ok {
			return ErrUnsupportedType
		}
		invalidator.SetActive(true)
		return nil
	})
}

func (s *datastoreStorage) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.getRequestEntity(ctx, s.PKCEKind, signature)
}

func (s *datastoreStorage) DeletePKCERequestSession(ctx context.Context, signature string) error {
	return s.deleteRequestEntity(ctx, s.PKCEKind, signature)
}
