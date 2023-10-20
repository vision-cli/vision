package server

import (
	"io/fs"
	"net/http"
)

func (s *Server) Routes() error {
	// special case: handker static content for Chi router with subroute, default go router does not require this step
	contentStatic, err := fs.Sub(fs.FS(static), "static")
	if err != nil {
		return err
	}
	s.r.Method(http.MethodGet, "/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(contentStatic))))

	s.r.Get("/", s.handleServeIndex())
	s.r.Get("/ping", s.handlePing())
	s.r.Get("/patients", s.handleGetPatient())
	s.r.Post("/patients", s.handleCreatePatient())
	return nil
}
