package todopb

import (
	"context"
	fmt "fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

var _ TodoServiceGraphQLInterface = (*todoServiceHandler)(nil)

func NewTodoServiceHandler(cli TodoServiceClient) TodoServiceGraphQLInterface {
	return &todoServiceHandler{cli}
}

type TodoServiceGraphQLInterface interface {
	CreateTodo(ctx context.Context, input CreateRequest) (*CreateResponse, error)
	TodosA(ctx context.Context, input ListARequest) (*ListAResponse, error)
	TodosB(ctx context.Context, input ListBRequest) (*ListBResponse, error)
	UpdateTodo(ctx context.Context, input UpdateRequest) (*UpdateResponse, error)
}

type todoServiceHandler struct {
	todoService TodoServiceClient
}

func (h *todoServiceHandler) CreateTodo(ctx context.Context, input CreateRequest) (*CreateResponse, error) {

	resp, err := h.todoService.Create(ctx, &input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *todoServiceHandler) TodosA(ctx context.Context, input ListARequest) (*ListAResponse, error) {

	resp, err := h.todoService.ListA(ctx, &input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *todoServiceHandler) TodosB(ctx context.Context, input ListBRequest) (*ListBResponse, error) {

	resp, err := h.todoService.ListB(ctx, &input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *todoServiceHandler) UpdateTodo(ctx context.Context, input UpdateRequest) (*UpdateResponse, error) {

	resp, err := h.todoService.Update(ctx, &input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func MarshalListADoneFilter(v ListARequest_DoneFilter) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(v.String()))
	})
}

func UnmarshalListADoneFilter(v interface{}) (ListARequest_DoneFilter, error) {
	if tmpStr, ok := v.(string); ok {
		v, ok := ListARequest_DoneFilter_value[tmpStr]
		if !ok {
			return 0, fmt.Errorf("invalid value format: %s", tmpStr)
		}
		return ListARequest_DoneFilter(v), nil
	}
	return 0, fmt.Errorf("unexpected value type: %T", v)
}

func MarshalListBDoneFilter(v ListBRequest_DoneFilter) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(v.String()))
	})
}

func UnmarshalListBDoneFilter(v interface{}) (ListBRequest_DoneFilter, error) {
	if tmpStr, ok := v.(string); ok {
		v, ok := ListBRequest_DoneFilter_value[tmpStr]
		if !ok {
			return 0, fmt.Errorf("invalid value format: %s", tmpStr)
		}
		return ListBRequest_DoneFilter(v), nil
	}
	return 0, fmt.Errorf("unexpected value type: %T", v)
}
