//go:generate gqlgen gen

package try_go_gqlgen

import (
	"context"
	"fmt"

	"github.com/vvakame/til/graphql/try-go-gqlgen/models"
)

func NewResolver() *Resolver {
	return &Resolver{
		latestTodoID: 2,
		UserMap: map[string]models.UserImpl{
			"User:123": {
				ID:   "User:123",
				Name: "123-san",
			},
			"User:245": {
				ID:   "User:245",
				Name: "245-san",
			},
		},
		todos: []models.Todo{
			{
				ID:     "Todo:1",
				Text:   "go to office",
				Done:   true,
				UserID: "User:123",
			},
			{
				ID:     "Todo:2",
				Text:   "make slide of mercari.go #2",
				Done:   false,
				UserID: "User:123",
			},
		},
	}
}

type Resolver struct {
	latestTodoID int
	UserMap      map[string]models.UserImpl
	todos        []models.Todo
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (models.Todo, error) {
	r.latestTodoID++
	todo := models.Todo{
		ID:     fmt.Sprintf("TODO:%d", r.latestTodoID),
		Text:   input.Text,
		Done:   false,
		UserID: input.UserID,
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]models.Todo, error) {
	return r.todos, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *models.Todo) (models.UserImpl, error) {
	user, err := ctx.Value(models.UserLoaderKey).(*models.UserImplLoader).Load(obj.UserID)
	if err != nil {
		return models.UserImpl{}, err
	}
	return *user, nil
}
