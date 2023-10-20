package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/atos-digital/NHSS-scigateway/internal/config"
)

type Server struct {
	r    *chi.Mux
	srv  *http.Server
	conf config.Config
	db   *gorm.DB
}

func New(conf config.Config) (*Server, error) {
	r := chi.NewRouter()

	// create a server so we can control its startup and shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Handler: r,
	}

	// create an instance of our server to wrap the http.Server
	s := &Server{
		r:    r,
		srv:  srv,
		conf: conf,
	}
	// setupDB will set the db var in the server struct
	err := s.setupDB(nil, defaultModels()...)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) ListenAndServe() error {
	s.middleware()
	s.Routes()
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
