package graphqlapi

import (
	"context"
	"net"
	"time"

	"github.com/akutz/memconn"
	"github.com/google/uuid"
	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

var _ todopb.TodoServiceClient = (*mockTodoServiceClient)(nil)

func ProvideTodoServiceClient(ctx context.Context) (todopb.TodoServiceClient, error) {
	return &mockTodoServiceClient{}, nil
}

func ProvideEchoClient(ctx context.Context) (echopb.EchoClient, error) {
	conn, err := grpc.Dial(
		"grpc",
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return memconn.Dial("memu", addr)
		}),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	)
	if err != nil {
		return nil, err
	}

	return echopb.NewEchoClient(conn), nil
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
