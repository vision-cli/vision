package database

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres returns a gorm.DB connection to a postgres database
func Postgres(url string, conf *gorm.Config) (*gorm.DB, error) {
	if conf == nil {
		conf = &gorm.Config{}
	}
	return gorm.Open(postgres.Open(url), conf)
}

// Json returns a datatypes.JSON from a value v
// https://gorm.io/docs/v2_release_note.html#DataTypes-JSON-as-example
func Json(v any) (datatypes.JSON, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}

// UnJson decodes a datatypes.JSON into a value v
// https://gorm.io/docs/v2_release_note.html#DataTypes-JSON-as-example
func UnJson(d datatypes.JSON, v any) error {
	return json.Unmarshal(d, &v)
}
