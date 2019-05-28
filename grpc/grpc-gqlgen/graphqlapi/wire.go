//+build wireinject

package graphqlapi

import (
	"context"

	"github.com/google/wire"
)

var grpcClientSet = wire.NewSet(
	ProvideTodoServiceClient,
	ProvideEchoClient,
)

func InitializeGraphQLConfig(ctx context.Context) (Config, error) {
	wire.Build(
		initializeResolvers,
		wire.Value(DirectiveRoot{}),
		wire.Value(ComplexityRoot{}),
		Config{},
	)

	return Config{}, nil
}

func initializeResolvers(ctx context.Context) (ResolverRoot, error) {
	wire.Build(
		resolver{},
		queryResolver{},
		mutationResolver{},

		grpcClientSet,
		todoServiceHandler{},
		echoHandler{},

		wire.Bind(new(ResolverRoot), new(resolver)),
	)

	return nil, nil
}
