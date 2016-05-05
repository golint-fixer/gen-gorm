package backends

import (
	"fmt"
	"testing"

	"github.com/kmulvey/gen-gorm/graph"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var users = graph.Vertex{
	Name: "Users", HasMany: "Posts", Cols: map[string]graph.Col{
		"Id":   {Name: "Id", Type: "int", Key: "PRI"},
		"Name": {Name: "Name", Type: "string", Key: ""},
	},
}
var id = graph.Col{Name: "Id", Type: "int", Key: "PRI"}
var userID = graph.Col{Name: "UserId", Type: "int", Key: "MUL"}

var database = graph.Graph{Name: "some_schema", Vertices: map[string]*graph.Vertex{
	"Users": &users,
	"Posts": {
		Name: "Posts", HasMany: "", Cols: map[string]graph.Col{
			"Id":     id,
			"Name":   {Name: "Name", Type: "string", Key: ""},
			"UserId": userID,
		},
		Edges: []graph.Edge{{DestinationTable: &users, DestinationCol: &id, OriginCol: &userID}},
	},
},
}

func TestGetTableInfoMysql(t *testing.T) {
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

	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").WillReturnRows(tableRows)
	mock.ExpectQuery("SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE").WillReturnRows(colRows)
	mock.ExpectQuery("SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE").WillReturnRows(colRowsTwo)
	mock.ExpectQuery("SELECT TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE").WillReturnRows(foreignKeys)
	m := Mysql{}
	s := "some_schema"
	c := ConnConfig{Schema: &s}
	data := m.createModel(db, c)
	assert.Equal(t, database, data, "should be equal")
}
