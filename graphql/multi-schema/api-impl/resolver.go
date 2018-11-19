package api_impl

import (
	"context"
	"fmt"
)

func NewResolver() *Resolver {
	r := &Resolver{}

	r.users = []User{
		{
			ID:    "User:1",
			Name:  "foobar",
			Staff: true,
		},
		{
			ID:    "User:2",
			Name:  "fizzbuzz",
			Staff: false,
		},
	}

	r.todos = []Todo{
		{
			ID:     "Todo:1",
			Text:   "Buy a milk",
			Done:   true,
			UserID: "User:1",
		},
		{
			ID:     "Todo:2",
			Text:   "Buy a cat food",
			Done:   false,
			UserID: "User:1",
		},
		{
			ID:     "Todo:1",
			Text:   "Check post box",
			Done:   false,
			UserID: "User:2",
		},
	}

	return r
}

func NewMutationResolver(r *Resolver) *mutationResolver {
	return &mutationResolver{r}
}

func NewQueryResolver(r *Resolver) *queryResolver {
	return &queryResolver{r}
}

func NewTodoResolver(r *Resolver) *todoResolver {
	return &todoResolver{r}
}

type Resolver struct {
	todos []Todo
	users []User
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (User, error) {
	user := User{
		ID:    fmt.Sprintf("User:%d", len(r.users)+1),
		Name:  input.Name,
		Staff: input.Staff,
	}
	r.users = append(r.users, user)
	return user, nil
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (Todo, error) {
	todo := Todo{
		ID:     fmt.Sprintf("Todo:%d", len(r.todos)+1),
		Text:   input.Text,
		Done:   false,
		UserID: input.UserID,
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]Todo, error) {
	return r.todos, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *Todo) (User, error) {
	for _, user := range r.users {
		if user.ID == obj.UserID {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("UserID: %s is not exists", obj.UserID)
}
