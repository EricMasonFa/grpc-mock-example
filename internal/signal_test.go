package internal

import (
	"context"
	"fmt"
	"grpc-mock-example/mocks"
	"grpc-mock-example/pb"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	fuzz "github.com/google/gofuzz"
)

func setup() {

}

func TestSignalService_Echo(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		sr  *pb.EchoRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.EchoResponse
		wantErr bool
	}{
		{
			name: "empty message",
			args: args{
				ctx: ctx,
				sr:  &pb.EchoRequest{Message: ""},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple message",
			args: args{
				ctx: ctx,
				sr:  &pb.EchoRequest{Message: "Hi there"},
			},
			want:    &pb.EchoResponse{Message: "Hi there"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SignalService{}
			got, err := s.Echo(tt.args.ctx, tt.args.sr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Echo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Echo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type echoRequest struct {
	message string
}

func (r *echoRequest) Matches(msg interface{}) bool {
	m, ok := msg.(*pb.EchoRequest)
	if !ok {
		return false
	}

	return m.Message == r.message
}

func (r *echoRequest) String() string {
	return fmt.Sprintf("is %s %v", r.message, r.message)
}

func TestSignalService_Ping(t *testing.T) {
	ctx := context.Background()
	var msg string
	f := fuzz.New()
	ctrl := gomock.NewController(t)
	signalClient := mocks.NewMockSignalClient(ctrl)

	f.Fuzz(&msg)

	signalClient.
		EXPECT().
		Echo(gomock.Any(), &echoRequest{message: ""}, gomock.Any()).
		Return(nil, ErrEmptyEcho).
		AnyTimes()

	signalClient.
		EXPECT().
		Echo(gomock.Any(), &echoRequest{message: msg}, gomock.Any()).
		Return(&pb.EchoResponse{Message: msg}, nil).
		AnyTimes()

	type args struct {
		ctx context.Context
		sr  *pb.PingRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.PingResponse
		wantErr bool
	}{
		{
			name: "lowercase",
			args: args{
				ctx: ctx,
				sr:  &pb.PingRequest{Message: "PiNG"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non ping message",
			args: args{
				ctx: ctx,
				sr:  &pb.PingRequest{Message: "POING"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ping",
			args: args{
				ctx: ctx,
				sr:  &pb.PingRequest{Message: "PING"},
			},
			want:    &pb.PingResponse{Message: "PONG"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SignalService{}
			got, err := s.Ping(tt.args.ctx, tt.args.sr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}
