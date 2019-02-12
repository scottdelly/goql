package db_client

import (
	"github.com/scottdelly/goql/models"
)

type SongClient struct {
	DBClient
}

func NewSongClient(dbc *DBClient) *SongClient {
	sc := new(SongClient)
	sc.DBClient = *dbc
	return sc
}

func newSong() *models.Song {
	return new(models.Song)
}
func emptySongs() []*models.Song {
	return make([]*models.Song, 0)
}

func (s *SongClient) GetSongs(limit uint64, where interface{}, args ...interface{}) ([]*models.Song, error) {
	songs := emptySongs()

	builder := s.Read(newSong())
	if where != nil {
		builder.
			Where(where, args...)
	}
	err :=
		builder.
			Limit(limit).
			QueryStructs(&songs)
	return songs, err
}

func (s *SongClient) GetSongById(id models.ModelId) (*models.Song, error) {
	song := newSong()
	err := s.GetByID(song, id, song)
	return song, err
}

func (s *SongClient) SongsBy(artistId models.ModelId) ([]*models.Song, error) {
	songs := emptySongs()
	err := s.
		Read(newSong()).
		Where(`"artist_id"" = $1`, artistId).
		QueryStructs(songs)
	if err != nil {
		return nil, err
	}
	return songs, nil
}
