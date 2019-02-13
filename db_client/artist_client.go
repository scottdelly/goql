package db_client

import (
	"github.com/scottdelly/goql/models"
)

type ArtistClient struct {
	DBClient
}

func NewArtistClient(dbc *DBClient) *ArtistClient {
	ac := new(ArtistClient)
	ac.DBClient = *dbc
	return ac
}

func newArtist() *models.Artist {
	return new(models.Artist)
}
func emptyArtists() []*models.Artist {
	return make([]*models.Artist, 0)
}

func (a *ArtistClient) GetArtists(limit uint64, where interface{}, args ...interface{}) ([]*models.Artist, error) {
	artists := emptyArtists()
	builder := a.Read(newArtist())
	if where != nil {
		builder.Where(where, args...)
	}
	err := builder.
		Limit(limit).
		QueryStructs(&artists)
	return artists, err
}

func (a *ArtistClient) GetArtistById(id models.ModelId) (*models.Artist, error) {
	artist := newArtist()
	err := a.GetByID(artist, id, artist)
	return artist, err
}

func (a *ArtistClient) ArtistFor(song *models.Song) (*models.Artist, error) {
	return a.GetArtistById(song.ArtistId)
}
