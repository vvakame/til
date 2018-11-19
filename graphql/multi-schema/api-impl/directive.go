package api_impl

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, requires *Role) (res interface{}, err error) {
	if requires == nil {
		return nil, errors.New("@hasRole required 'requires' parameter value")
	}

	switch *requires {
	case RolePublic:
		// OK
	case RoleStaff:
		// めんどいからOKにするけど本来はちゃんとアレコレしてよね
		log.Printf("hasRole: %s", *requires)
	default:
		return nil, fmt.Errorf("unexpected role value: %s", *requires)
	}

	return next(ctx)
}
