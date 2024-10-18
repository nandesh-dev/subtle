package server

import (
	"fmt"
	"net"

	"github.com/nandesh-dev/subtle/generated/api/media"
	media_service "github.com/nandesh-dev/subtle/internal/server/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	listener   *net.Listener
	grpcServer *grpc.Server
}

func New() *server {
	return &server{}
}

func (s *server) Listen(port int, enableReflection bool) error {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	s.listener = &listener

	s.grpcServer = grpc.NewServer()

	if enableReflection {
		reflection.Register(s.grpcServer)
	}

	mediaService := media_service.MediaServiceServer{}
	media.RegisterMediaServiceServer(s.grpcServer, &mediaService)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
