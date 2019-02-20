package models

type Artist struct {
	Model
}

//Conforms to CRUDModel
func (a *Artist) TableName(operation CRUDOperation) string {
	return "artists"
}

func (a *Artist) ColumnNames(operation CRUDOperation) []string {
	columns := a.Model.ColumnNames(operation)
	return columns
}

func (a *Artist) Values(operation CRUDOperation) []interface{} {
	values := a.Model.Values(operation)
	return values
}
