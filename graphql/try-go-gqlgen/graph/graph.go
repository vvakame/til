//go:generate gqlgen -schema ../schema.graphql -typemap types.json

package graph

import (
	"context"
	"fmt"
	"math/rand"
)

type MyApp struct {
	todos   []Todo
	userMap map[string]UserImpl
}

type Todo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID string `json:"-"`
}

type UserImpl struct {
	ID   string
	Name string
}

func NewMyApp() *MyApp {
	return &MyApp{
		userMap: make(map[string]UserImpl),
	}
}

func (a *MyApp) Query_todos(ctx context.Context) ([]Todo, error) {
	return a.todos, nil
}

func (a *MyApp) Mutation_createTodo(ctx context.Context, text string) (Todo, error) {
	user := UserImpl{
		ID:   fmt.Sprintf("U%d", rand.Int()),
		Name: fmt.Sprintf("Name of U%d", rand.Int()),
	}
	todo := Todo{
		Text:   text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: user.ID,
	}
	a.userMap[user.ID] = user
	a.todos = append(a.todos, todo)
	return todo, nil
}

func (a *MyApp) Todo_user(ctx context.Context, obj *Todo) (UserImpl, error) {
	user := a.userMap[obj.UserID]
	return user, nil
}
