package grpcapi

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ echopb.EchoServer = (*echoServiceImpl)(nil)

func NewEchoServer() (echopb.EchoServer, error) {
	return &echoServiceImpl{}, nil
}

type echoServiceImpl struct {
}

func (srv *echoServiceImpl) Say(ctx context.Context, req *echopb.SayRequest) (*echopb.SayResponse, error) {
	ctx, span := trace.StartSpan(ctx, "echoService/Say")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("messageId", req.MessageId))
	span.AddAttributes(trace.StringAttribute("messageBody", req.MessageBody))

	if strings.Contains(req.MessageBody, "error") {
		span.SetStatus(trace.Status{
			Code:    int32(codes.InvalidArgument),
			Message: "messageBody contains 'error'",
		})
		return nil, status.Error(codes.InvalidArgument, "messageBody contains 'error'")
	}

	return &echopb.SayResponse{
		MessageId:   req.MessageId,
		MessageBody: req.MessageBody,
		Received:    ptypes.TimestampNow(),
	}, nil
}

func (srv *echoServiceImpl) SaySecond(ctx context.Context, req *echopb.SayRequest) (*echopb.SayResponse, error) {
	return srv.Say(ctx, req)
}
