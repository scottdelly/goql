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
	columns := s.Model.ColumnNames(operation)
	columns = append(columns, "artist_id", "duration")
	return columns
}

func (s *Song) Values(operation CRUDOperation) []interface{} {
	values := s.Model.Values(operation)
	values = append(values, s.ArtistId, s.Duration)
	return values
}
