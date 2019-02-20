package models

type LikeObjectType int

const LikeTypeArtist LikeObjectType = 0
const LikeTypeSong LikeObjectType = 1

type Likable interface {
	Identifiable
	LikeObjectType() LikeObjectType
}

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

func (u *User) LikeObject(likable Likable) CreateLikeMessage {
	return CreateLikeMessage{
		User:       u,
		ObjectType: likable.LikeObjectType(),
		Object:     likable,
	}
}

func (u *User) LikesOfType(t LikeObjectType) ReadUserLikesMessage {
	return ReadUserLikesMessage{
		User:       u,
		ObjectType: t,
	}
}

func LikesOn(l Likable) ReadObjectLikesMessage {
	return ReadObjectLikesMessage{
		ObjectType: l.LikeObjectType(),
		Object:     l,
	}
}

//Artist Conforms To Likable
func (a *Artist) LikeObjectType() LikeObjectType {
	return LikeTypeArtist
}

//Song Conforms To Likable
func (s *Song) LikeObjectType() LikeObjectType {
	return LikeTypeSong
}
