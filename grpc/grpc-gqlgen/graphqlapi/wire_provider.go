package graphqlapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
	"google.golang.org/grpc"
)

var _ todopb.TodoServiceClient = (*mockTodoServiceClient)(nil)

func ProvideTodoServiceClient(ctx context.Context) (todopb.TodoServiceClient, error) {
	return &mockTodoServiceClient{}, nil
}

type mockTodoServiceClient struct {
	todos []*todopb.Todo
}

func (mock *mockTodoServiceClient) Create(ctx context.Context, in *todopb.CreateRequest, opts ...grpc.CallOption) (*todopb.CreateResponse, error) {
	todo := &todopb.Todo{
		Id:   uuid.New().String(),
		Text: in.Text,
	}
	mock.todos = append(mock.todos, todo)

	return &todopb.CreateResponse{Todo: todo}, nil
}

func (mock *mockTodoServiceClient) Update(ctx context.Context, in *todopb.UpdateRequest, opts ...grpc.CallOption) (*todopb.UpdateResponse, error) {
	for _, todo := range mock.todos {
		if in.Id == todo.Id {
			if in.Text != "" {
				todo.Text = in.Text
			}
			todo.Done = in.Done
			return &todopb.UpdateResponse{Todo: todo}, nil
		}
	}

	return nil, nil
}

func (mock *mockTodoServiceClient) ListA(ctx context.Context, in *todopb.ListARequest, opts ...grpc.CallOption) (*todopb.ListAResponse, error) {
	// TODO ちゃんと
	return &todopb.ListAResponse{Todos: mock.todos}, nil
}

func (mock *mockTodoServiceClient) ListB(ctx context.Context, in *todopb.ListBRequest, opts ...grpc.CallOption) (*todopb.ListBResponse, error) {
	// TODO ちゃんと
	return &todopb.ListBResponse{Todos: mock.todos}, nil
}
