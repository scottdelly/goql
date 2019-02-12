package db_client

import (
	"database/sql"
	"time"

	"gopkg.in/mgutz/dat.v2/dat"
	"gopkg.in/mgutz/dat.v2/sqlx-runner"

	"github.com/scottdelly/goql/models"
)

var IdColumn = "id"

type DBClient struct {
	db *runner.DB
}

func (d *DBClient) Start() {
	// create a normal database connection through database/sql
	db, err := sql.Open("postgres", "dbname=dat_test user=dat password=!test host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}

	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	d.db = runner.NewDB(db, "postgres")
}

func (d *DBClient) Create(i models.CRUDModel) error {
	_, err := d.db.
		InsertInto(i.TableName(models.CRUDCreate)).
		Columns(i.ColumnNames(models.CRUDCreate)...).
		Values(i.Values(models.CRUDCreate)...).
		Exec()
	return err
}

func (d *DBClient) Read(i models.CRUDModel) *dat.SelectBuilder {
	return d.db.
		Select(i.ColumnNames(models.CRUDRead)...).
		From(i.TableName(models.CRUDRead))
}

func (d *DBClient) GetByID(i models.CRUDModel, id models.ModelId, response interface{}) error {
	err := d.
		Read(i).
		Where(`$1 = $2`, IdColumn, id).
		Limit(1).
		QueryStruct(response)
	return err
}

func (d *DBClient) Update(i models.CRUDModel) *dat.UpdateBuilder {
	columns := i.ColumnNames(models.CRUDUpdate)
	values := i.Values(models.CRUDUpdate)
	updateMap := make(map[string]interface{}, 0)
	for i, col := range columns {
		updateMap[col] = values[i]
	}

	return d.db.
		Update(i.TableName(models.CRUDUpdate)).
		SetMap(updateMap)
}

func (d *DBClient) UpdateById(i models.CRUDModel) error {
	_, err := d.
		Update(i).
		Where(`$1 = $2`, "id", i.Identifier()).
		Exec()
	return err
}

func (d *DBClient) Delete(i models.CRUDModel) *dat.DeleteBuilder {
	return d.db.
		DeleteFrom(i.TableName(models.CRUDDelete))
}

func (d *DBClient) DeleteByID(i models.CRUDModel) error {
	_, err := d.
		Delete(i).
		Where(`$1 = $2`, "id", i.Identifier()).
		Exec()
	return err
}
