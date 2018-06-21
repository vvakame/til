//go:generate gqlgen -schema ../schema.graphql -typemap types.json

package graph

import (
	"context"
	"fmt"
	"math/rand"
	"github.com/vvakame/til/graphql/try-go-gqlgen/models"
)

type MyApp struct {
	todos   []models.Todo
	userMap map[string]models.UserImpl
}

func NewMyApp() *MyApp {
	return &MyApp{
		userMap: make(map[string]models.UserImpl),
	}
}

func (a *MyApp) Query_todos(ctx context.Context) ([]models.Todo, error) {
	return a.todos, nil
}

func (a *MyApp) Mutation_createTodo(ctx context.Context, text string) (models.Todo, error) {
	user := models.UserImpl{
		ID:   fmt.Sprintf("U%d", rand.Int()),
		Name: fmt.Sprintf("Name of U%d", rand.Int()),
	}
	todo := models.Todo{
		Text:   text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: user.ID,
	}
	a.userMap[user.ID] = user
	a.todos = append(a.todos, todo)
	return todo, nil
}

func (a *MyApp) Todo_user(ctx context.Context, obj *models.Todo) (models.UserImpl, error) {
	user := a.userMap[obj.UserID]
	return user, nil
}
