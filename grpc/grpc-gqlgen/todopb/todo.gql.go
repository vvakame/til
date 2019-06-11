package todopb

import (
	"context"
)

var _ TodoServiceGraphQLInterface = (*todoServiceHandler)(nil)

type TodoServiceGraphQLInterface interface {
	CreateTodo(ctx context.Context, input *CreateRequest) (*CreateResponse, error)
	TodosA(ctx context.Context, input *ListARequest) (*ListAResponse, error)
	TodosB(ctx context.Context, input *ListBRequest) (*ListBResponse, error)
	UpdateTodo(ctx context.Context, input *UpdateRequest) (*UpdateResponse, error)
}

type todoServiceHandler struct {
	todoService TodoServiceClient
}

func (h *todoServiceHandler) CreateTodo(ctx context.Context, input *CreateRequest) (*CreateResponse, error) {

	resp, err := h.todoService.Create(ctx, input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *todoServiceHandler) TodosA(ctx context.Context, input *ListARequest) (*ListAResponse, error) {

	resp, err := h.todoService.ListA(ctx, input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *todoServiceHandler) TodosB(ctx context.Context, input *ListBRequest) (*ListBResponse, error) {

	resp, err := h.todoService.ListB(ctx, input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *todoServiceHandler) UpdateTodo(ctx context.Context, input *UpdateRequest) (*UpdateResponse, error) {

	resp, err := h.todoService.Update(ctx, input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}
