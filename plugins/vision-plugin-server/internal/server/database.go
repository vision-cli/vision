package server

import (
	"gorm.io/gorm"

	"github.com/atos-digital/NHSS-scigateway/internal/database"
	"github.com/atos-digital/NHSS-scigateway/internal/models/tables"
)

func (s *Server) setupDB(gormConf *gorm.Config, models ...any) error {
	db, err := database.Postgres(s.conf.DatabaseURL, gormConf)
	if err != nil {
		return err
	}
	err = db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Server) TeardownDB() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func defaultModels() []any {
	return []any{
		&tables.Patients{},
	}
}
