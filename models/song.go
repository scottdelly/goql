package models

import (
	"time"
)

type Song struct {
	Model
	ArtistId ModelId       `db:"artist_id" json:"artist_id"`
	Duration time.Duration `db:"duration" json:"duration"`
}

//Conforms to CRUDModel
func (s *Song) TableName(operation CRUDOperation) string {
	return "songs"
}

func (s *Song) ColumnNames(operation CRUDOperation) []string {
	return []string{"id", "name", "duration"}
}

func (s *Song) Values(operation CRUDOperation) []interface{} {
	return []interface{}{s.ID, s.Name, s.Duration}
}
