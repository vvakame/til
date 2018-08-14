//go:generate gorunpkg github.com/99designs/gqlgen

package gqlgen_directive_parent

import (
	"context"
)

func NewResolver() *Resolver {
	return &Resolver{
		userMap: map[string]User{
			"User:123": {
				ID:   "User:123",
				Name: "123-san",
			},
			"User:345": {
				ID:   "User:345",
				Name: "123-345",
			},
		},
		todos: []Todo{
			{
				ID:   "Todo:abc",
				Text: "foobar",
				User: User{
					ID:   "User:123",
					Name: "123-san",
				},
			},
			{
				ID:   "Todo:efg",
				Text: "fizzbuzz",
				User: User{
					ID:   "User:345",
					Name: "123-345",
				},
			},
		},
	}
}

type Resolver struct {
	todos   []Todo
	userMap map[string]User
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]Todo, error) {
	return r.todos, nil
}
