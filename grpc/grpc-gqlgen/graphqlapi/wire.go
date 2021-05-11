//+build wireinject

package graphqlapi

import (
	"context"

	"github.com/google/wire"
	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
)

var grpcClientSet = wire.NewSet(
	ProvideTodoServiceClient,
	ProvideEchoClient,
)

var gqlHandlerSet = wire.NewSet(
	echopb.NewEchoHandler,
	todopb.NewTodoServiceHandler,
)

func InitializeGraphQLConfig(ctx context.Context) (Config, error) {
	wire.Build(
		initializeResolvers,
		wire.Value(DirectiveRoot{}),
		wire.Value(ComplexityRoot{}),
		wire.Struct(new(Config), "*"),
	)

	return Config{}, nil
}

func initializeResolvers(ctx context.Context) (ResolverRoot, error) {
	wire.Build(
		grpcClientSet,
		gqlHandlerSet,

		wire.Struct(new(queryResolver), "*"),
		wire.Struct(new(mutationResolver), "*"),
		wire.Struct(new(resolver), "*"),
		wire.Bind(new(ResolverRoot), new(*resolver)),
	)

	return nil, nil
}
