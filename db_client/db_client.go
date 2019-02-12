package db_client

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/mgutz/dat.v2/dat"
	"gopkg.in/mgutz/dat.v2/sqlx-runner"

	"github.com/scottdelly/goql/models"
)

type DBClient struct {
	db *runner.DB
}

func (d *DBClient) Start(user, pass, host string) {
	// create a normal database connection through database/sql
	println(fmt.Sprintf("DB Starting"))

	db, err := sql.Open("postgres",
		fmt.Sprintf("dbname=postgres user=%s password=%s host=%s sslmode=disable", user, pass, host))
	if err != nil {
		panic(err)
	}

	runner.MustPing(db)
	println(fmt.Sprintf("DB Reached"))

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
		Where(`"id" = $1`, id).
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
		Where(`"id" = $2`, i.Identifier()).
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
		Where(`"id" = $1`, i.Identifier()).
		Exec()
	return err
}
