package models

import (
	"time"
)

type Song struct {
	Model
	ArtistId ModelId
	Duration time.Duration `db:"duration" json:"duration"`
}

//Conforms to CRUDModel
func (s *Song) TableName(operation CRUDOperation) string {
	return "song"
}

func (s *Song) ColumnNames(operation CRUDOperation) []string {
	return []string{"id", "name", "duration"}
}

func (s *Song) Values(operation CRUDOperation) []interface{} {
	return []interface{}{s.ID, s.Name, s.Duration}
}
