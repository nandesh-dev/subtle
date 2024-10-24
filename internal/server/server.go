package server

import (
	"fmt"
	"net"
	"net/http"

	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"github.com/nandesh-dev/subtle/generated/proto/media/mediaconnect"
	"github.com/nandesh-dev/subtle/internal/server/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type server struct {
	listener   *net.Listener
	grpcServer *grpc.Server
}

func New() *server {
	return &server{}
}

func (s *server) Listen(port int, enableReflection bool) error {
	mux := http.NewServeMux()

	path, handler := mediaconnect.NewMediaServiceHandler(media.ServiceHandler{})
	mux.Handle(path, handler)

	if config.Config().Server.GRPCReflection {
		path, handler = grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(
			mediaconnect.MediaServiceName))
		mux.Handle(path, handler)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: config.Config().Server.COROrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})

	http.ListenAndServe(
		fmt.Sprintf("localhost:%v", port),
		h2c.NewHandler(corsMiddleware.Handler(mux), &http2.Server{}),
	)

	return nil
}
