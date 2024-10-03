.PHONY: proto run

proto:
	protoc ./proto/**/*.proto --go_out=. --go-grpc_out=. 

run:
	go run ./cmd/subtle
