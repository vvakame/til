package graphqlapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
)

var _ todoServiceGraphQLInterface = (*todoServiceHandler)(nil)

type todoServiceGraphQLInterface interface {
	TodosA(ctx context.Context, first *int, after *string, input TodoListAInput) (*TodoConnection, error)
	TodosB(ctx context.Context, first *int, after *string, input TodoListBInput) (*TodoConnection, error)
	CreateTodo(ctx context.Context, input CreateTodoInput) (*CreateTodoPayload, error)
	UpdateTodo(ctx context.Context, input UpdateTodoInput) (*UpdateTodoPayload, error)
}

// TODO ID変換レイヤーが必要

type todoServiceHandler struct {
	todoService todopb.TodoServiceClient
}

func (h *todoServiceHandler) TodosA(ctx context.Context, first *int, after *string, input TodoListAInput) (*TodoConnection, error) {
	in := &todopb.ListARequest{}
	if first != nil {
		// TODO エラーチェック
		in.First = uint32(*first)
	}
	if after != nil {
		in.After = *after
	}
	// TODO ここ厳しい…
	if input.NotDone == nil {
		in.Done = todopb.ListARequest_NONE
	} else if *input.NotDone {
		in.Done = todopb.ListARequest_NOT_DONE
	} else {
		in.Done = todopb.ListARequest_DONE
	}

	resp, err := h.todoService.ListA(ctx, in)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}
	type Connection = TodoConnection
	type Edge = TodoEdge
	listItems := resp.GetTodos

	conn := &Connection{
		PageInfo: &PageInfo{},
		Edges:    nil,
		Nodes:    nil,
	}
	if resp.GetCursor() != "" {
		conn.PageInfo.EndCursor = proto.String(resp.GetCursor())
		conn.PageInfo.HasNextPage = true
	}
	if in.GetAfter() != "" {
		conn.PageInfo.HasPreviousPage = true
	}
	for _, item := range listItems() {
		conn.Nodes = append(conn.Nodes, item)
		conn.Edges = append(conn.Edges, &Edge{
			// TODO Cursor
			Node: item,
		})
	}

	return conn, nil
}

func (h *todoServiceHandler) TodosB(ctx context.Context, first *int, after *string, input TodoListBInput) (*TodoConnection, error) {
	in := &todopb.ListBRequest{}
	if first != nil {
		// TODO エラーチェック
		in.Limit = uint32(*first)
	} else {
		// TODO なんかセットしないとカーソルの有無がわからん
	}
	if after != nil {
		vs, err := url.ParseQuery(*after)
		if err != nil {
			// TODO もっとよいエラー
			return nil, err
		}
		offset, err := strconv.Atoi(vs.Get("offset"))
		if err != nil {
			// TODO もっとよいエラー
			return nil, err
		}
		// TODO エラーチェック
		in.Offset = uint32(offset)
	}
	// TODO ここ厳しい…
	if input.NotDone == nil {
		in.Done = todopb.ListBRequest_NONE
	} else if *input.NotDone {
		in.Done = todopb.ListBRequest_NOT_DONE
	} else {
		in.Done = todopb.ListBRequest_DONE
	}

	resp, err := h.todoService.ListB(ctx, in)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}
	type Connection = TodoConnection
	type Edge = TodoEdge
	listItems := resp.GetTodos

	conn := &Connection{
		PageInfo: &PageInfo{},
		Edges:    nil,
		Nodes:    nil,
	}
	if in.Limit == uint32(len(listItems())) {
		var vs url.Values
		vs.Set("offset", fmt.Sprintf("%d", in.Limit+in.Offset))
		conn.PageInfo.EndCursor = proto.String(vs.Encode())
		conn.PageInfo.HasNextPage = true
	}
	if in.Offset != 0 {
		conn.PageInfo.HasPreviousPage = true
	}
	for _, item := range listItems() {
		conn.Nodes = append(conn.Nodes, item)
		conn.Edges = append(conn.Edges, &Edge{
			// TODO Cursor
			Node: item,
		})
	}

	return conn, nil
}

func (h *todoServiceHandler) CreateTodo(ctx context.Context, input CreateTodoInput) (*CreateTodoPayload, error) {
	in := &todopb.CreateRequest{}
	in.Text = input.Text

	resp, err := h.todoService.Create(ctx, in)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	payload := &CreateTodoPayload{}
	payload.Todo = resp.GetTodo()

	return payload, nil
}

func (h *todoServiceHandler) UpdateTodo(ctx context.Context, input UpdateTodoInput) (*UpdateTodoPayload, error) {
	in := &todopb.UpdateRequest{}
	// TODO ID変換のレイヤーが必要では？
	in.Id = input.ID
	if input.Text != nil {
		// TODO これだと 値の有無は表現できていない…
		in.Text = *input.Text
	}
	if input.Done != nil {
		// TODO これだと 3状態は表現できていない…
		in.Done = *input.Done
	}

	resp, err := h.todoService.Update(ctx, in)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	payload := &UpdateTodoPayload{}
	payload.Todo = resp.GetTodo()

	return payload, nil
}
