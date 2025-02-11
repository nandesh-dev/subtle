package api

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"github.com/nandesh-dev/subtle/generated/proto/web/webconnect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type WebServiceHandler struct {
	webconnect.UnimplementedWebServiceHandler
}

type APIServer struct {
	handler http.Handler
}

type APIServerOptions struct {
	EnableGRPCReflection bool
	COROrigins           []string
}

func NewAPIServer(options APIServerOptions) *APIServer {
	mux := http.NewServeMux()

	path, handler := webconnect.NewWebServiceHandler(WebServiceHandler{})
	mux.Handle(path, handler)

	if options.EnableGRPCReflection {
		path, handler = grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(webconnect.WebServiceName))
		mux.Handle(path, handler)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: options.COROrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})

	return &APIServer{
		handler: h2c.NewHandler(corsMiddleware.Handler(mux), &http2.Server{}),
	}
}

func (s *APIServer) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.handler)
}
