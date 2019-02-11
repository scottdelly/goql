package models

type ModelId int

type Model struct {
	ID   ModelId `db:"id" json:"id"`
	Name string  `db:"name" json:"name"`
}

func (m *Model) Identifier() ModelId {
	return m.ID
}

type CRUDOperation int

const (
	CRUDCreate = 0
	CRUDRead   = 1
	CRUDUpdate = 2
	CRUDDelete = 3
)

type CRUDModel interface {
	Identifier() ModelId
	TableName(operation CRUDOperation) string
	ColumnNames(operation CRUDOperation) []string
	Values(operation CRUDOperation) []interface{}
}
