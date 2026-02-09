package server

import (
	"context"
	"my_blog_backend/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, serverCfg *config.HttpServerCfg) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + serverCfg.Port,
			Handler:      handler,
			ReadTimeout:  serverCfg.ReadTimeout,
			WriteTimeout: serverCfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
