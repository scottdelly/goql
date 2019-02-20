package models

type LikeObjectType int

const LikeTypeArtist LikeObjectType = 0
const LikeTypeSong LikeObjectType = 1

type LikeMessage struct {
	User       *User
	ObjectType LikeObjectType
	Object     Identifiable
}

func (u *User) LikeArtist(artist *Artist) LikeMessage {
	return LikeMessage{
		User:       u,
		ObjectType: LikeTypeArtist,
		Object:     artist,
	}
}

func (u *User) LikeSong(song *Song) LikeMessage {
	return LikeMessage{
		User:       u,
		ObjectType: LikeTypeSong,
		Object:     song,
	}
}
