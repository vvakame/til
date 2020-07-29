package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vvakame/til/go/graphql-todos/domains"
	"github.com/vvakame/til/go/graphql-todos/graph/generated"
	"github.com/vvakame/til/go/graphql-todos/graph/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*domains.Todo, error) {
	return r.TodoRepository.Create(ctx, &domains.Todo{
		Text: input.Text,
	})
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input model.UpdateTodoInput) (*domains.Todo, error) {
	todo, err := r.TodoRepository.Get(ctx, domains.TodoID(input.ID))
	if err != nil {
		return nil, err
	}

	if input.Text != nil {
		todo.Text = *input.Text
	}
	if input.Done != nil {
		todo.Done = *input.Done
	}

	return r.TodoRepository.Update(ctx, todo)
}

func (r *queryResolver) Todos(ctx context.Context) ([]*domains.Todo, error) {
	return r.TodoRepository.GetAll(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
