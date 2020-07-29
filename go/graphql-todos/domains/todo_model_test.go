package domains

import (
	"context"
	"testing"
)

func Test_todoRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewTodoRepository()

	todoA, err := repo.Create(ctx, &Todo{
		Text: "test A",
	})
	if err != nil {
		t.Fatal(err)
	}
	if v := todoA.Done; v {
		t.Errorf("unexpected: %#v", v)
	}
	if v := todoA.DoneAt; !v.IsZero() {
		t.Errorf("unexpected: %#v", v)
	}
	if v := todoA.CreatedAt; v.IsZero() {
		t.Errorf("unexpected: %#v", v)
	}

	todoANew, err := repo.Update(ctx, &Todo{
		ID:   todoA.ID,
		Text: "test A!",
		Done: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if v := todoANew.Done; !v {
		t.Errorf("unexpected: %#v", v)
	}
	if v := todoANew.DoneAt; v.IsZero() {
		t.Errorf("unexpected: %#v", v)
	}

	todoA2, err := repo.Get(ctx, todoA.ID)
	if err != nil {
		t.Fatal(err)
	}
	if v1, v2 := todoA.ID, todoA2.ID; v1 != v2 {
		t.Errorf("unexpected: %#v, %#v", v1, v2)
	}

	todoB, err := repo.Create(ctx, &Todo{
		Text: "test B",
	})
	if err != nil {
		t.Fatal(err)
	}

	list, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v := len(list); v != 2 {
		t.Fatalf("unexpected: %#v", v)
	}
	if v := list[0]; v.ID != todoB.ID {
		t.Errorf("unexpected: %#v", v)
	}
	if v := list[1]; v.ID != todoA.ID {
		t.Errorf("unexpected: %#v", v)
	}
}
