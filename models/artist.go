package models

type Artist struct {
	Model
	LikeCount int `db:"like_count" json:"like_count"`
}

//Conforms to CRUDModel
func (a *Artist) TableName(operation CRUDOperation) string {
	return "artists"
}

func (a *Artist) ColumnNames(operation CRUDOperation) []string {
	columns := a.Model.ColumnNames(operation)
	columns = append(columns, "like_count")
	return columns
}

func (a *Artist) Values(operation CRUDOperation) []interface{} {
	values := a.Model.Values(operation)
	values = append(values, a.LikeCount)
	return values
}
