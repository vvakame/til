//go:generate gqlgen

package graphqlapi

import (
	"context"

	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
)

var _ ResolverRoot = (*resolver)(nil)

type resolver struct {
	queryResolver    *queryResolver
	mutationResolver *mutationResolver
}

func (r *resolver) Query() QueryResolver {
	return r.queryResolver
}

func (r *resolver) Mutation() MutationResolver {
	return r.mutationResolver
}

var _ QueryResolver = (*queryResolver)(nil)

type queryResolver struct {
	todopb.TodoServiceGraphQLInterface
	echopb.EchoGraphQLInterface
}

func (r *queryResolver) Tmp(ctx context.Context) (*string, error) {
	return nil, nil
}

var _ MutationResolver = (*mutationResolver)(nil)

type mutationResolver struct {
	todopb.TodoServiceGraphQLInterface
	echopb.EchoGraphQLInterface
}

func (r *mutationResolver) Tmp(ctx context.Context) (*string, error) {
	return nil, nil
}
