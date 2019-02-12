package models

type User struct {
	Model
	Email string `db:"email" json:"email"`
}

//Conforms to CRUDModel
func (u *User) TableName(operation CRUDOperation) string {
	return "user"
}

func (u *User) ColumnNames(operation CRUDOperation) []string {
	return []string{"id", "name", "email"}

}

func (u *User) Values(operation CRUDOperation) []interface{} {
	return []interface{}{u.ID, u.Name, u.Email}
}
