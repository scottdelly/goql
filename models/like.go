package models

/*
The representation of a Like is not a 1:1 Model like other DB representation. Instead, the Like Model is syntax sugar
to create "messages" that the LikeClient uses to build queries for likes.

The main intention for this is to make it easier to reason about querying for likes at a higher level. For example:
- this user likes this song
- show me users who like this artist
- show me songs liked by this user

Instead of making unique functions for each type that can be liked, a type that can be liked need only conform to the
Likable interface, and the existing functions will work for it.
*/

type LikeObjectType int

const LikeTypeArtist LikeObjectType = 0
const LikeTypeSong LikeObjectType = 1

type Likable interface {
	Identifiable
	LikeObjectType() LikeObjectType
}

//Declare conformance to likable
var _ Likable = (*Artist)(nil)
var _ Likable = (*Song)(nil)

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
