//go:generate gorunpkg github.com/99designs/gqlgen

package api_private

import "github.com/vvakame/til/graphql/multi-schema/api-impl"

var _ ResolverRoot = (*Resolver)(nil)

type Resolver api_impl.Resolver

func NewResolver() ResolverRoot {
	return (*Resolver)(api_impl.NewResolver())
}

func (r *Resolver) Mutation() MutationResolver {
	return api_impl.NewMutationResolver((*api_impl.Resolver)(r))
}

func (r *Resolver) Query() QueryResolver {
	return api_impl.NewQueryResolver((*api_impl.Resolver)(r))
}

func (r *Resolver) Todo() TodoResolver {
	return api_impl.NewTodoResolver((*api_impl.Resolver)(r))
}
