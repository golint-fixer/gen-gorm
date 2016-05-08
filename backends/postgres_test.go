package backends

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetTableInfoPostgres(t *testing.T) {
	t.Parallel()

	// Open new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	// columns to be used for result
	tableRows := sqlmock.NewRows([]string{"table_name"}).
		AddRow("users").
		AddRow("posts")
	colRows := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE", "COLUMN_KEY"}).
		AddRow("id", "int", "PRI").
		AddRow("name", "varchar", "")
	// you cant reuse mocked rows
	colRowsTwo := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE", "COLUMN_KEY"}).
		AddRow("id", "int", "PRI").
		AddRow("name", "varchar", "").
		AddRow("user_id", "int", "MUL")
	foreignKeys := sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "REFERENCED_TABLE_NAME", "REFERENCED_COLUMN_NAME"}).
		AddRow("posts", "user_id", "users", "id")

	mock.ExpectQuery("SELECT tablename FROM pg_catalog.pg_tables where tableowner").WillReturnRows(tableRows)
	mock.ExpectQuery("select kc.column_name, c.data_type, tc.constraint_type from information_schema.table_constraints tc").WillReturnRows(colRows)
	mock.ExpectQuery("select kc.column_name, c.data_type, tc.constraint_type from information_schema.table_constraints tc").WillReturnRows(colRowsTwo)
	mock.ExpectQuery("select tc.table_name, kcu.column_name, ccu.table_name AS references_table, ccu.column_name AS references_field").WillReturnRows(foreignKeys)
	p := Postgres{}
	s, u := "some_schema", "adama"
	c := ConnConfig{Schema: &s, Username: &u}
	fmt.Println(c)
	data := p.createModel(db, c)
	assert.Equal(t, database, data, "should be equal")
}
