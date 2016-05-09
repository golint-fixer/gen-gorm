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
	colRows := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE", "character_maximum_length", "COLUMN_KEY"}).
		AddRow("id", "int", nil, "PRI").
		AddRow("name", "varchar", "45", "")
	// you cant reuse mocked rows
	colRowsTwo := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE", "character_maximum_length", "COLUMN_KEY"}).
		AddRow("id", "int", nil, "PRI").
		AddRow("name", "varchar", "45", "").
		AddRow("user_id", "int", nil, "FOREIGN KEY")
	foreignKeys := sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "REFERENCED_TABLE_NAME", "REFERENCED_COLUMN_NAME"}).
		AddRow("posts", "user_id", "users", "id")

	mock.ExpectQuery("SELECT tablename FROM pg_catalog.pg_tables where tableowner").WillReturnRows(tableRows)
	mock.ExpectQuery("select kc.column_name, c.data_type, c.character_maximum_length, tc.constraint_type from information_schema.table_constraints").WillReturnRows(colRows)
	mock.ExpectQuery("select kc.column_name, c.data_type, c.character_maximum_length, tc.constraint_type from information_schema.table_constraints").WillReturnRows(colRowsTwo)
	mock.ExpectQuery("select tc.table_name, kcu.column_name, ccu.table_name AS references_table, ccu.column_name AS references_field").WillReturnRows(foreignKeys)
	p := Postgres{}
	s, u := "some_schema", "adama"
	c := ConnConfig{Schema: &s, Username: &u}
	fmt.Println(database.Vertices["Posts"].Edges[0].DestinationTable.HasMany)
	data := p.createModel(db, c)
	assert.Equal(t, database, data, "should be equal")
}
