package main

import (
	"os"
	"testing"

	"github.com/kmulvey/gen-gorm/graph"
	"github.com/stretchr/testify/assert"
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
		Edges: []graph.Edge{graph.Edge{DestinationTable: &users, DestinationCol: &id, OriginCol: &userID}},
	},
},
}

func TestProcessTemplates(t *testing.T) {
	t.Parallel()
	var dir = "dist"

	// delete it if it exists
	if _, err := os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}
	processTemplates(database, "dist")

	_, err := os.Stat(dir)
	assert.NoError(t, err, "directory should exist")
	f, err := os.Stat(dir + "/models.go")
	assert.NoError(t, err, "models.go should exist")
	assert.True(t, f.Size() > 0, "file should not be empty")

	// clean up
	os.RemoveAll(dir)
}
