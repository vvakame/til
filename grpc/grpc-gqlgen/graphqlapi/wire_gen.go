// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package graphqlapi

import (
	"context"
	"github.com/google/wire"
	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
)

// Injectors from wire.go:

func InitializeGraphQLConfig(ctx context.Context) (Config, error) {
	resolverRoot, err := initializeResolvers(ctx)
	if err != nil {
		return Config{}, err
	}
	directiveRoot := _wireDirectiveRootValue
	complexityRoot := _wireComplexityRootValue
	config := Config{
		Resolvers:  resolverRoot,
		Directives: directiveRoot,
		Complexity: complexityRoot,
	}
	return config, nil
}

var (
	_wireDirectiveRootValue  = DirectiveRoot{}
	_wireComplexityRootValue = ComplexityRoot{}
)

func initializeResolvers(ctx context.Context) (ResolverRoot, error) {
	todoServiceClient, err := ProvideTodoServiceClient(ctx)
	if err != nil {
		return nil, err
	}
	todoServiceGraphQLInterface := todopb.NewTodoServiceHandler(todoServiceClient)
	echoClient, err := ProvideEchoClient(ctx)
	if err != nil {
		return nil, err
	}
	echoGraphQLInterface := echopb.NewEchoHandler(echoClient)
	graphqlapiQueryResolver := &queryResolver{
		TodoServiceGraphQLInterface: todoServiceGraphQLInterface,
		EchoGraphQLInterface:        echoGraphQLInterface,
	}
	graphqlapiMutationResolver := &mutationResolver{
		TodoServiceGraphQLInterface: todoServiceGraphQLInterface,
		EchoGraphQLInterface:        echoGraphQLInterface,
	}
	graphqlapiResolver := &resolver{
		queryResolver:    graphqlapiQueryResolver,
		mutationResolver: graphqlapiMutationResolver,
	}
	return graphqlapiResolver, nil
}

// wire.go:

var grpcClientSet = wire.NewSet(
	ProvideTodoServiceClient,
	ProvideEchoClient,
)

var gqlHandlerSet = wire.NewSet(echopb.NewEchoHandler, todopb.NewTodoServiceHandler)
