package tables

import (
	"gorm.io/datatypes"
)

type Patients struct {
	ID   string `gorm:"primaryKey"`
	Data datatypes.JSON
}
