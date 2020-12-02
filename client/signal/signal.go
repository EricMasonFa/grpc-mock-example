package signal

import (
	"context"
	"encoding/json"
	"grpc-mock-example/pb"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

var (
	conn   *grpc.ClientConn
	client pb.SignalClient
)

// New gets the Signal's gPRC server as serverAddr, opens a connection and
// return a pb.SignalClient
func New(serverAddr string) pb.SignalClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client = pb.NewSignalClient(conn)

	return client
}

// Ping is a http.HandlerFunc which calls the pb.SignalClient.Ping with url
// querystring and returns the json formatted result
func Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := r.URL.Path[len("/ping/"):]
	ping, err := client.Ping(ctx, &pb.PingRequest{Message: message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e, ok := status.FromError(err)

		if ok {
			json.NewEncoder(w).Encode(e.Message())
			return
		}

		em := err.Error()
		log.Printf("%s", em)
		json.NewEncoder(w).Encode(em)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(ping)
}

// Echo is a http.HandlerFunc which calls the pb.SignalClient.Echo with url
// querystring and returns the json formatted result
func Echo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := r.URL.Path[len("/echo/"):]
	echo, err := client.Echo(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e, ok := status.FromError(err)

		if ok {
			json.NewEncoder(w).Encode(e.Message())
			return
		}

		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(echo)
}

func Close() {
	conn.Close()
}
