package server

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/atos-digital/NHSS-scigateway/internal/auth"
)

func (s *Server) middleware() {
	s.r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   s.conf.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
		middleware.Recoverer,
		middleware.Compress(5),
		middleware.Logger,
		auth.Authenicate,
	)
}
