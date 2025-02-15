package api

import (
	"net/http"
	"path/filepath"

	connectcors "connectrpc.com/cors"
	"connectrpc.com/grpcreflect"
	"github.com/nandesh-dev/subtle/generated/embed"
	"github.com/nandesh-dev/subtle/generated/proto/web/webconnect"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type WebServiceHandler struct {
	webconnect.UnimplementedWebServiceHandler
	configFile *configuration.File
}

type APIServer struct {
	handler http.Handler
}

type APIServerOptions struct {
	EnableGRPCReflection bool
	COROrigins           []string
}

func NewAPIServer(configFile *configuration.File, options APIServerOptions) *APIServer {
	mux := http.NewServeMux()

	path, handler := webconnect.NewWebServiceHandler(WebServiceHandler{configFile: configFile})
	mux.Handle(path, handler)

	if options.EnableGRPCReflection {
		path, handler = grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(webconnect.WebServiceName))
		mux.Handle(path, handler)
	}

	frontendFilesystem := embed.GetFrontendFilesystem()
	frontendFileServer := http.FileServer(http.FS(frontendFilesystem))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler, handlerPattern := mux.Handler(r)
		if handlerPattern != "" && handlerPattern != "/" {
			handler.ServeHTTP(w, r)
			return
		}

		if filepath.Ext(r.URL.Path) == "" {
			http.ServeFileFS(w, r, frontendFilesystem, "index.html")
      return
		}

		frontendFileServer.ServeHTTP(w, r)
	})

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
