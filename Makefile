.PHONY: proto run

proto:
	rm -rf generated
	mkdir generated
	protoc ./proto/**/*.proto --go_out=. --go-grpc_out=. 
	rm -rf web/generated
	mkdir web/generated
	protoc --grpc-web_out=import_style=typescript,mode=grpcwebtext:web/generated/  proto/**/*.proto

run:
	go run ./cmd/subtle
