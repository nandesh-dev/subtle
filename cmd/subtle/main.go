package main

import (
	"log"
	"net"

	"github.com/nandesh-dev/subtle/internal/pb/library"
	"github.com/nandesh-dev/subtle/internal/services"
	"github.com/nandesh-dev/subtle/internal/tesseract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	client := tesseract.NewClient()
	defer client.Close()

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln("failed to create listener: ", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	libraryService := services.Services().Library
	library.RegisterLibraryServiceServer(s, &libraryService)

	if err := s.Serve(listener); err != nil {
		log.Fatalln("failed to serve: ", err)
	}
}
