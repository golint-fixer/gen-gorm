package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kmulvey/gen-gorm/graph"
	"github.com/stretchr/testify/assert"
)

var database = graph.Graph{Name: "some_schema", Vertices: map[string]graph.Vertex{
	"Users": graph.Vertex{
		Name: "Users", Cols: map[string]graph.Col{
			"Id":   graph.Col{Name: "Id", Type: "int"},
			"Name": graph.Col{Name: "Name", Type: "string"},
		},
	},
	"Posts": graph.Vertex{
		Name: "Posts", Cols: map[string]graph.Col{
			"Id":     graph.Col{Name: "Id", Type: "int"},
			"Name":   graph.Col{Name: "Name", Type: "string"},
			"UserId": graph.Col{Name: "UserId", Type: "string"},
		},
	},
},
}

func TestConvertType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "string", convertType("time"), "should be string")
	assert.Equal(t, "string", convertType("datetime"), "should be string")
	assert.Equal(t, "string", convertType("char"), "should be string")
	assert.Equal(t, "string", convertType("varchar"), "should be string")
	assert.Equal(t, "string", convertType("blob"), "should be string")
	assert.Equal(t, "int", convertType("integer"), "should be integer")
	assert.Equal(t, "int", convertType("int"), "should be integer")
	assert.Equal(t, "int", convertType("timestamp"), "should be integer")
	assert.Equal(t, "float", convertType("float"), "should be float")
	assert.Equal(t, "float", convertType("double"), "should be float")
	assert.Equal(t, "bool", convertType("boolean"), "should be bool")
	assert.Equal(t, "string", convertType("enum"), "should be string")
	assert.Equal(t, "string", convertType("other"), "should be string")
}

func TestCapFirst(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "String", capFirst("string"), "should be \"String\"")
	assert.Equal(t, "Int", capFirst("int"), "should be \"Int\"")
}

func TestTormatColName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Id", formatColName("id"), "should be Id")
	assert.Equal(t, "UserName", formatColName("user_name"), "should be UserName")
	assert.Equal(t, "Username", formatColName("username"), "should be Username")
	assert.Equal(t, "Col5", formatColName("col_5"), "should be Col5")
	assert.Equal(t, "Code", formatColName("code"), "should be Code")
}

func TestGetTableInfo(t *testing.T) {
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
	colRows := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE"}).
		AddRow("id", "int").
		AddRow("name", "varchar")
	// you cant reuse mocked rows
	colRowsTwo := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE"}).
		AddRow("id", "int").
		AddRow("name", "varchar").
		AddRow("user_id", "varchar")
	foreignKeys := sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "REFERENCED_TABLE_NAME", "REFERENCED_COLUMN_NAME"}).
		AddRow("posts", "user_id", "users", "id")

	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").WillReturnRows(tableRows)
	mock.ExpectQuery("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE").WillReturnRows(colRows)
	mock.ExpectQuery("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE").WillReturnRows(colRowsTwo)
	mock.ExpectQuery("SELECT TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE").WillReturnRows(foreignKeys)
	data := getTableInfo(db, "some_schema")
	assert.EqualValues(t, database, data, "should be equal")
}

/*

func TestProcessTemplates(t *testing.T) {
	t.Parallel()
	var dir = "dist"

	// delete it if it exists
	if _, err := os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}
	processTemplates(expected, "dist")

	_, err := os.Stat(dir)
	assert.NoError(t, err, "directory should exist")
	f, err := os.Stat(dir + "/struct.go")
	assert.NoError(t, err, "struct.go should exist")
	assert.True(t, f.Size() > 0, "file should not be empty")

	// clean up
	os.RemoveAll(dir)
}
*/

func TestHandleError(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() {
		handleErr(errors.New("some error"))
	}, "Calling handleErr() should panic")
}
