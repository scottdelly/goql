package db_client

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/mgutz/dat.v2/dat"
	"gopkg.in/mgutz/dat.v2/sqlx-runner"

	"github.com/scottdelly/goql/models"
)

const dbName = "postgres"

type DBClient struct {
	db *runner.DB
}

func (d *DBClient) Start(user, pass, host string, bootstrap bool) {
	// create a normal database connection through database/sql
	log.Println("==> DB Starting")

	db, err := sql.Open("postgres",
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", dbName, user, pass, host))
	if err != nil {
		panic(err)
	}

	runner.MustPing(db)
	log.Println("==> DB Connected")

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

	if bootstrap {
		d.bootstrapDB()
	}
	log.Println("==> DB Ready")
}

func (d *DBClient) Create(i models.CRUDModel) error {
	err := d.db.
		InsertInto(i.TableName(models.CRUDCreate)).
		Columns(i.ColumnNames(models.CRUDCreate)...).
		Values(i.Values(models.CRUDCreate)...).
		Returning(i.ColumnNames(models.CRUDRead)...).
		QueryStruct(i)
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

func (d *DBClient) bootstrapDB() {
	log.Println("==> Bootstrapping DB")
	resetSQL, err := ioutil.ReadFile("./migrations/reset.sql")
	if err != nil {
		panic(err)
	}
	bootstrapSQL, err := ioutil.ReadFile("./migrations/bootstrap.sql")
	if err != nil {
		panic(err)
	}
	tx, err := d.db.DB.Begin()
	_, err = tx.Exec(string(resetSQL))
	_, err = tx.Exec(string(bootstrapSQL))
	if err != nil {
		_ = tx.Rollback()
		panic(err)
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	log.Println("==> Bootstrapping Complete")
}
