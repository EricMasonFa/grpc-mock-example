TEST?=./...

default: test

test:
	go list $(TEST) | xargs -t -n4 go test -v $(TESTARGS) -timeout=2m -parallel=4 -count=1

protobuf:
	protoc -I/opt/include -I/usr/local/include -Ipb \
		--go_out=./pb --go_opt=paths=source_relative \
    	--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
    	pb/*.proto

genmocks:
	@GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4
	@mkdir -p ./mocks
	mockgen -source=pb/signal_grpc.pb.go -package=mocks SignalClient > ./mocks/signal_client_mock.go;

all: protobuf genmocks test

.PHONY: test protobuf genmocks all
