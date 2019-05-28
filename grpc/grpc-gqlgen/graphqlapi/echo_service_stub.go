package graphqlapi

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
)

var _ echoGraphQLInterface = (*echoHandler)(nil)

type echoGraphQLInterface interface {
	Say(ctx context.Context, input SayInput) (*SayPayload, error)
}

// TODO ID変換レイヤーが必要

type echoHandler struct {
	echo echopb.EchoClient
}

func (h *echoHandler) Say(ctx context.Context, input SayInput) (*SayPayload, error) {
	in := &echopb.SayRequest{}
	if input.ClientMutationID != nil {
		in.MessageId = *input.ClientMutationID
	}
	in.MessageBody = input.MessageBody

	resp, err := h.echo.Say(ctx, in)
	if err != nil {
		// TODO なんらかのエラーハンドラが必要なはず
		return nil, err
	}

	payload := &SayPayload{}
	if resp.MessageId != "" {
		payload.ClientMutationID = proto.String(resp.GetMessageId())
	}
	payload.MessageBody = resp.GetMessageBody()
	payload.Received = *resp.GetReceived()

	return payload, nil
}
