package db_client

import (
	"github.com/scottdelly/goql/models"
)

type UserClient struct {
	DBClient
}

func NewUserClient(dbc *DBClient) *UserClient {
	uc := new(UserClient)
	uc.DBClient = *dbc
	return uc
}

func newUser() *models.User {
	return new(models.User)
}
func emptyUsers() []*models.User {
	return make([]*models.User, 0)
}

func (u *UserClient) GetUsers(limit uint64, where map[string]interface{}, args ...interface{}) ([]*models.User, error) {
	users := emptyUsers()
	err := u.
		Read(newUser()).
		Where(where, args...).
		Limit(limit).
		QueryStructs(&users)
	return users, err
}

func (u *UserClient) GetUserById(id models.ModelId) (*models.User, error) {
	user := newUser()
	err := u.GetByID(user, id, user)
	return user, err
}
