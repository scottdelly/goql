package db_client

type UserClient struct {
	DB
}

func (u *UserClient) fetchUsers(limit int) {
	_limit := limit ? limit : 10
}
