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
	return []string{"id", "name", "like_count"}
}

func (a *Artist) Values(operation CRUDOperation) []interface{} {
	return []interface{}{a.ID, a.Name, a.LikeCount}
}
