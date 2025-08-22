package server

import (
	"context"
	"my_blog_backend/internal/config"
	"net/http"
	"os"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, httpServer config.HttpServer) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + os.Getenv("HTTP_PORT"),
			Handler:      handler,
			ReadTimeout:  httpServer.ReadTimeout,
			WriteTimeout: httpServer.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
