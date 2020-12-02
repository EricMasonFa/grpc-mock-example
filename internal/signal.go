package internal

import (
	"context"
	"errors"
	"grpc-mock-example/pb"
)

var (
	ErrInvalidPing = errors.New("invalid ping message")
	ErrEmptyEcho   = errors.New("message can't be empty")
)

// SignalService implements pb.SignalServer
type SignalService struct {
	pb.UnimplementedSignalServer
}

// Ping returns "PONG" message only when recives "PING" Message, and
// error otherwise
func (s *SignalService) Ping(_ context.Context, sr *pb.PingRequest) (*pb.PingResponse, error) {
	if sr != nil && sr.Message == "PING" {
		return &pb.PingResponse{Message: "PONG"}, nil
	}

	return nil, ErrInvalidPing
}

// Echo returns the same message that is received, and error if the
// message is empty
func (s *SignalService) Echo(_ context.Context, sr *pb.EchoRequest) (*pb.EchoResponse, error) {
	if sr != nil && sr.Message != "" {
		return &pb.EchoResponse{Message: sr.Message}, nil
	}

	return nil, ErrEmptyEcho
}
