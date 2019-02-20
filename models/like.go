package models

type LikeObjectType int

const LikeTypeArtist LikeObjectType = 0
const LikeTypeSong LikeObjectType = 1

type CreateLikeMessage struct {
	User       *User
	ObjectType LikeObjectType
	Object     Identifiable
}

type ReadUserLikesMessage struct {
	User       *User
	ObjectType LikeObjectType
}

type ReadObjectLikesMessage struct {
	ObjectType LikeObjectType
	Object     Identifiable
}

func (u *User) LikeArtist(artist *Artist) CreateLikeMessage {
	return CreateLikeMessage{
		User:       u,
		ObjectType: LikeTypeArtist,
		Object:     artist,
	}
}

func (u *User) LikeSong(song *Song) CreateLikeMessage {
	return CreateLikeMessage{
		User:       u,
		ObjectType: LikeTypeSong,
		Object:     song,
	}
}

func (u *User) LikesOfType(t LikeObjectType) ReadUserLikesMessage {
	return ReadUserLikesMessage{
		User:       u,
		ObjectType: t,
	}
}

func (a *Artist) Likes() ReadObjectLikesMessage {
	return ReadObjectLikesMessage{
		ObjectType: LikeTypeArtist,
		Object:     a,
	}
}

func (s *Song) Likes() ReadObjectLikesMessage {
	return ReadObjectLikesMessage{
		ObjectType: LikeTypeSong,
		Object:     s,
	}
}
