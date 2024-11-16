package server

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"

	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"github.com/nandesh-dev/subtle/generated/proto/media/mediaconnect"
	"github.com/nandesh-dev/subtle/generated/proto/subtitle/subtitleconnect"
	"github.com/nandesh-dev/subtle/internal/server/media"
	"github.com/nandesh-dev/subtle/internal/server/subtitle"
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

func (s *server) Listen() error {
	apiMux := http.NewServeMux()

	path, handler := mediaconnect.NewMediaServiceHandler(media.ServiceHandler{})
	apiMux.Handle(path, handler)

	path, handler = subtitleconnect.NewSubtitleServiceHandler(subtitle.ServiceHandler{})
	apiMux.Handle(path, handler)

	if config.Config().Server.Web.EnableGRPCReflection {
		path, handler = grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(
			mediaconnect.MediaServiceName))
		apiMux.Handle(path, handler)
	}

	fileServer := http.FileServer(http.Dir(config.Config().Web.ServeDirectory))
	serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := apiMux.Handler(r)
		if pattern != "" {
			apiMux.ServeHTTP(w, r)
			return
		}

		if filepath.Ext(r.URL.Path) == "" || filepath.Ext(r.URL.Path) == ".html" {
			http.ServeFile(w, r, filepath.Join(config.Config().Web.ServeDirectory, "index.html"))
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: config.Config().Server.Web.COROrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})

	http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%v", config.Config().Server.Web.Port),
		h2c.NewHandler(corsMiddleware.Handler(serverHandler), &http2.Server{}),
	)

	return nil
}
