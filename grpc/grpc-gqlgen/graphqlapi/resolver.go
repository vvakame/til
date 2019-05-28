//go:generate gqlgen
//go:generate wire .

package graphqlapi

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

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
	*todoServiceHandler
}

var _ MutationResolver = (*mutationResolver)(nil)

type mutationResolver struct {
	*todoServiceHandler
}
