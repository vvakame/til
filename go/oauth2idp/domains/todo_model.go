package domains

import (
	"context"
	"go.mercari.io/datastore"
	"time"
)

var _ datastore.PropertyLoadSaver = (*Todo)(nil)

type Todo struct {
	ID        int64  `datastore:"-" boom:"id"`
	Text      string `validate:"req"`
	Complete  bool   ``
	UserID    int64  ``
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (todo *Todo) Load(ctx context.Context, ps []datastore.Property) error {
	return datastore.LoadStruct(ctx, todo, ps)
}

func (todo *Todo) Save(ctx context.Context) ([]datastore.Property, error) {
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = time.Now()
	}
	todo.UpdatedAt = time.Now()

	err := Validate(todo)
	if err != nil {
		return nil, err
	}

	return datastore.SaveStruct(ctx, todo)
}
