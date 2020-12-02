package main

import (
	"flag"
	"grpc-mock-example/client/signal"
	"log"
	"net/http"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "Server address in the format of host:port")
)

func main() {
	flag.Parse()

	signal.New(*serverAddr)
	defer signal.Close()

	http.HandleFunc("/ping/", signal.Ping)
	http.HandleFunc("/echo/", signal.Echo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
