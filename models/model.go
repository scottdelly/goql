package models

import (
	"time"
)

type ModelId int

type Model struct {
	ID      ModelId   `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	Created time.Time `db:"created" json:"created"`
}

type Identifiable interface {
	Identifier() ModelId
}

func (m Model) Identifier() ModelId {
	return m.ID
}

//Demonstrate and check conformance to Identifiable
var _ Identifiable = (*User)(nil)
var _ Identifiable = (*Artist)(nil)
var _ Identifiable = (*Song)(nil)

type Nameable interface {
	NameValue() string
}

func (m Model) NameValue() string {
	return m.Name
}

//Demonstrate and check conformance to Nameable
var _ Nameable = (*User)(nil)
var _ Nameable = (*Artist)(nil)
var _ Nameable = (*Song)(nil)

type Historical interface {
	DateValue() time.Time
}

func (m Model) DateValue() time.Time {
	//time.Now().UTC()
	return m.Created
}

//Demonstrate and check conformance to Historical
var _ Historical = (*User)(nil)
var _ Historical = (*Artist)(nil)
var _ Historical = (*Song)(nil)

type CRUDOperation int

const (
	CRUDCreate = 0
	CRUDRead   = 1
	CRUDUpdate = 2
	CRUDDelete = 3
)

type CRUDModel interface {
	Identifiable
	TableName(operation CRUDOperation) string
	ColumnNames(operation CRUDOperation) []string
	Values(operation CRUDOperation) []interface{}
}

//Demonstrate and check conformance to CRUDModel
var _ CRUDModel = (*User)(nil)
var _ CRUDModel = (*Artist)(nil)
var _ CRUDModel = (*Song)(nil)

func (m Model) ColumnNames(operation CRUDOperation) []string {
	switch operation {
	case CRUDRead:
		return []string{"id", "name", "created"}
	default: //create or update should not change id or created date
		return []string{"name"}
	}
}

func (m Model) Values(operation CRUDOperation) []interface{} {
	switch operation {
	case CRUDRead:
		return []interface{}{m.ID, m.Name, m.Created}
	default: //create or update should not change id or created date
		return []interface{}{m.Name}
	}
}
