package dsstorage

import (
	"errors"

	"github.com/ory/fosite"
)

var ErrNoSuchEntity = fosite.ErrNotFound

var ErrUnsupportedType = errors.New("unsupported type")

var ErrInvalidTxContext = errors.New("context doesn't in tx context")
