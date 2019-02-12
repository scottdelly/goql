package models

type User struct {
	Model
	Email string `db:"email" json:"email"`
}

//Conforms to CRUDModel
func (u *User) TableName(operation CRUDOperation) string {
	return "users"
}

func (u *User) ColumnNames(operation CRUDOperation) []string {
	columns := u.Model.ColumnNames(operation)
	columns = append(columns, "email")
	return columns

}

func (u *User) Values(operation CRUDOperation) []interface{} {
	values := u.Model.Values(operation)
	values = append(values, u.Email)
	return values
}
