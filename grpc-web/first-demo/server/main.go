package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/vvakame/til/grpc-web/first-demo/server/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	chatServer := &chatServer{}

	chat.RegisterChatServiceServer(s, chatServer)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

var _ chat.ChatServiceServer = (*chatServer)(nil)

type chatServer struct {
	Speaks    []*chat.SpeakResponse
	Listeners []chat.ChatService_ListenSpeakServer
}

func (cs *chatServer) Speak(ctx context.Context, req *chat.SpeakRequest) (*chat.SpeakResponse, error) {
	now, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}
	resp := &chat.SpeakResponse{
		SpeakId:  int32(len(cs.Speaks) + 1),
		UserName: req.UserName,
		Message:  req.Message,
		SpeakAt:  now,
	}
	cs.Speaks = append(cs.Speaks, resp)

	for _, listener := range cs.Listeners {
		err := listener.Send(resp)
		log.Printf("error on listener.Send: %s", err.Error())
	}

	return resp, nil
}

func (cs *chatServer) ListenSpeak(req *chat.ListenSpeakRequest, stream chat.ChatService_ListenSpeakServer) error {
	cs.Listeners = append(cs.Listeners, stream)
	return nil
}
