package db_client

import (
	"errors"
	"fmt"

	"github.com/scottdelly/goql/models"
)

const LikesTableArtist = "likes_artists"
const LikesTableSong = "likes_songs"

const UserIdColumn = "user_id"
const ArtistIdColumn = "artist_id"
const SongIdColumn = "song_id"

type LikesClient struct {
	*DBClient
}

func NewLikesClient(dbc *DBClient) *LikesClient {
	lc := new(LikesClient)
	lc.DBClient = dbc
	return lc
}

func (lc *LikesClient) GetLikesForUser(lm models.ReadUserLikesMessage, limit uint64) ([]models.ModelId, error) {
	var result = make([]models.ModelId, 0)
	err := lc.db.Select(relIdColumnForType(lm.ObjectType)).
		From(tableNameForType(lm.ObjectType)).
		Where(fmt.Sprintf(`%s = $1`, UserIdColumn), lm.User.Id).
		Limit(limit).
		QuerySlice(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lc *LikesClient) GetUsersByLikes(lm models.ReadObjectLikesMessage, limit uint64) ([]models.ModelId, error) {
	var result = make([]models.ModelId, 0)
	err := lc.db.Select(UserIdColumn).
		From(tableNameForType(lm.ObjectType)).
		Where(fmt.Sprintf(`%s = $1`, relIdColumnForType(lm.ObjectType)), lm.Object.Identifier()).
		Limit(limit).
		QuerySlice(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lc *LikesClient) LikeCount(lm models.ReadObjectLikesMessage) (int, error) {
	var count int
	err := lc.db.Select(`COUNT(*)`).
		From(tableNameForType(lm.ObjectType)).
		Where(fmt.Sprintf(`%s = $1`, relIdColumnForType(lm.ObjectType)), lm.Object.Identifier()).
		QueryScalar(&count)
	return count, err
}

func (lc *LikesClient) CreateUserLike(lm models.CreateLikeMessage) error {
	var err error
	if _, err = lc.validateUser(lm.User.Id); err != nil {
		return err
	}

	if err = lc.checkExistingLikes(lm); err != nil {
		return err
	}
	if err = lc.createLike(lm); err != nil {
		return err
	}
	return nil
}

func (lc *LikesClient) ResolveObject(objectId models.ModelId, likeType models.LikeObjectType) (models.Likable, error) {
	var likable models.Likable
	var err error
	switch likeType {
	case models.LikeTypeArtist:
		if likable, err = NewArtistClient(lc.DBClient).GetArtistById(objectId); err != nil {
			return nil, err
		}

	case models.LikeTypeSong:
		if likable, err = NewSongClient(lc.DBClient).GetSongById(objectId); err != nil {
			return nil, err
		}
	}
	return likable, nil
}

func (lc *LikesClient) validateUser(userId models.ModelId) (*models.User, error) {
	uc := NewUserClient(lc.DBClient)
	if user, err := uc.GetUserById(userId); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (lc *LikesClient) checkExistingLikes(lm models.CreateLikeMessage) error {
	count := 0
	objectKey := relIdColumnForType(lm.ObjectType)
	if err := lc.db.Select(`COUNT(*)`).
		From(tableNameForType(lm.ObjectType)).
		Where(fmt.Sprintf(`%s = $1 AND %s = $2`, UserIdColumn, objectKey), lm.User.Id, lm.Object.Identifier()).
		QueryScalar(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("like relation already exists for user: %d, %s: %d", lm.User.Id, objectKey, lm.Object.Identifier()))
	}
	return nil
}

func (lc *LikesClient) createLike(lm models.CreateLikeMessage) error {
	if _, err := lc.db.InsertInto(tableNameForType(lm.ObjectType)).
		Columns(UserIdColumn, relIdColumnForType(lm.ObjectType)).
		Values(lm.User.Id, lm.Object.Identifier()).
		Exec(); err != nil {
		return err
	}
	return nil
}

func tableNameForType(likeType models.LikeObjectType) string {
	switch likeType {
	case models.LikeTypeArtist:
		return LikesTableArtist
	case models.LikeTypeSong:
		return LikesTableSong
	}
	panic(errors.New("like type/table mismatch"))

}

func relIdColumnForType(likeType models.LikeObjectType) string {
	switch likeType {
	case models.LikeTypeArtist:
		return ArtistIdColumn
	case models.LikeTypeSong:
		return SongIdColumn
	}
	panic(errors.New("like type/relId column mismatch"))
}
