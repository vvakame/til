//go:generate gqlgen

package graphqlapi

import "context"

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
	todoServiceGraphQLInterface
	echoGraphQLInterface
}

func (r *queryResolver) Tmp(ctx context.Context) (*string, error) {
	return nil, nil
}

var _ MutationResolver = (*mutationResolver)(nil)

type mutationResolver struct {
	todoServiceGraphQLInterface
	echoGraphQLInterface
}

func (r *mutationResolver) Tmp(ctx context.Context) (*string, error) {
	return nil, nil
}
