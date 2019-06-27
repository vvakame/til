package echopb

import (
	"context"
)

var _ EchoGraphQLInterface = (*echoHandler)(nil)

func NewEchoHandler(cli EchoClient) EchoGraphQLInterface {
	return &echoHandler{cli}
}

type EchoGraphQLInterface interface {
	Say(ctx context.Context, input SayRequest) (*SayResponse, error)
	SaySecond(ctx context.Context, input SayRequest) (*SayResponse, error)
}

type echoHandler struct {
	echo EchoClient
}

func (h *echoHandler) Say(ctx context.Context, input SayRequest) (*SayResponse, error) {

	resp, err := h.echo.Say(ctx, &input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}

func (h *echoHandler) SaySecond(ctx context.Context, input SayRequest) (*SayResponse, error) {

	resp, err := h.echo.SaySecond(ctx, &input)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	return resp, nil
}
