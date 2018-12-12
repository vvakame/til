package main

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/vvakame/til/appengine/go111-internal-grpc/echopb"
)

var _ echopb.EchoServer = (*echoServiceImpl)(nil)

type echoServiceImpl struct {
}

func (srv *echoServiceImpl) Say(ctx context.Context, req *echopb.SayRequest) (*echopb.SayResponse, error) {
	return &echopb.SayResponse{
		MessageId:   req.MessageId,
		MessageBody: req.MessageBody,
		Received:    ptypes.TimestampNow(),
	}, nil
}
