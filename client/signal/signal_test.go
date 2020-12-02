package signal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grpc-mock-example/internal"
	"grpc-mock-example/mocks"
	"grpc-mock-example/pb"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	fuzz "github.com/google/gofuzz"
)

// echoRequest keeps the minimum requirement for matching Echo requests.
// It implements gomock.Matcher
type echoRequest struct {
	message string
}

// Matches compares echoRequest.message and pb.EchoRequest.Message
// returns true if equal and false otherwise
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

// echoRequest keeps the minimum requirement for matching Echo requests.
// It implements gomock.Matcher
type pingRequest struct {
	message string
}

// Matches compares echoRequest.message and pb.EchoRequest.Message
// returns true if equal and false otherwise
func (r *pingRequest) Matches(msg interface{}) bool {
	m, ok := msg.(*pb.PingRequest)
	if !ok {
		return false
	}

	return m.Message == r.message
}

func (r *pingRequest) String() string {
	return fmt.Sprintf("is %s %v", r.message, r.message)
}

type customResponseWriter struct {
	header     http.Header
	body       *bytes.Buffer
	statusCode int
}

func (c customResponseWriter) Header() http.Header {
	return c.header
}

func (c customResponseWriter) Write(body []byte) (int, error) {
	return c.body.Write(body)
}

func (c customResponseWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
}

func TestEcho(t *testing.T) {
	var msg string
	f := fuzz.New()
	ctrl := gomock.NewController(t)
	signalClient := mocks.NewMockSignalClient(ctrl)

	f.Fuzz(&msg)
	jsonMsg := &bytes.Buffer{}
	json.HTMLEscape(jsonMsg, []byte(msg))

	signalClient.
		EXPECT().
		Echo(gomock.Any(), &echoRequest{message: ""}, gomock.Any()).
		Return(nil, internal.ErrEmptyEcho).
		AnyTimes()

	signalClient.
		EXPECT().
		Echo(gomock.Any(), &echoRequest{message: msg}, gomock.Any()).
		Return(&pb.EchoResponse{Message: msg}, nil).
		AnyTimes()

	client = signalClient

	type args struct {
		w customResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{
				w: customResponseWriter{
					header: http.Header{},
					body:   &bytes.Buffer{},
				},
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/echo/",
					},
				},
			},
			want: "\"message can't be empty\"\n",
		},
		{
			name: "msg",
			args: args{
				w: customResponseWriter{
					header: http.Header{},
					body:   &bytes.Buffer{},
				},
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/echo/" + msg,
					},
				},
			},
			want: fmt.Sprintf("{\"message\":\"%s\"}\n", jsonMsg),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Echo(tt.args.w, tt.args.r)

			got := tt.args.w.body.String()
			if !reflect.DeepEqual(tt.args.w.body.String(), tt.want) {
				t.Errorf("TestPing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing(t *testing.T) {
	var msg string
	f := fuzz.New()
	ctrl := gomock.NewController(t)
	signalClient := mocks.NewMockSignalClient(ctrl)

	f.Fuzz(&msg)

	signalClient.
		EXPECT().
		Ping(gomock.Any(), &pingRequest{message: "Ping"}, gomock.Any()).
		Return(nil, internal.ErrInvalidPing).
		AnyTimes()

	signalClient.
		EXPECT().
		Ping(gomock.Any(), &pingRequest{message: "PING"}, gomock.Any()).
		Return(&pb.PingResponse{Message: "PONG"}, nil).
		AnyTimes()

	client = signalClient

	type args struct {
		w customResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Ping",
			args: args{
				w: customResponseWriter{
					header: http.Header{},
					body:   &bytes.Buffer{},
				},
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/ping/Ping",
					},
				},
			},
			want: "\"invalid ping message\"\n",
		},
		{
			name: "PING",
			args: args{
				w: customResponseWriter{
					header: http.Header{},
					body:   &bytes.Buffer{},
				},
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/ping/PING",
					},
				},
			},
			want: "{\"message\":\"PONG\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Ping(tt.args.w, tt.args.r)

			got := tt.args.w.body.String()
			if !reflect.DeepEqual(tt.args.w.body.String(), tt.want) {
				t.Errorf("TestPing() = %v, want %v", got, tt.want)
			}
		})
	}
}
