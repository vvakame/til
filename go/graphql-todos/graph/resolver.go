package graph

import "github.com/vvakame/til/go/graphql-todos/domains"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoRepository domains.TodoRepository
}
