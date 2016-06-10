package backends

import (
	"database/sql"

	"github.com/kmulvey/gen-gorm/graph"
)

var users = graph.Vertex{
	Name: "Users", HasMany: "Posts", Cols: map[string]graph.Col{
		"Id":   {Name: "Id", Type: "int", Key: "primary_key", MaxLen: sql.NullInt64{Int64: 0, Valid: false}},
		"Name": {Name: "Name", Type: "string", Key: "", MaxLen: sql.NullInt64{Int64: 45, Valid: true}},
	},
}
var id = graph.Col{Name: "Id", Type: "int", Key: "primary_key", MaxLen: sql.NullInt64{Int64: 0, Valid: false}}
var userID = graph.Col{Name: "UserId", Type: "int", Key: "MULTIPLE", MaxLen: sql.NullInt64{Int64: 0, Valid: false}}

var database = graph.Graph{Name: "some_schema", Vertices: map[string]*graph.Vertex{
	"Users": &users,
	"Posts": {
		Name: "Posts", HasMany: "", Cols: map[string]graph.Col{
			"Id":     id,
			"Name":   {Name: "Name", Type: "string", Key: "", MaxLen: sql.NullInt64{Int64: 45, Valid: true}},
			"UserId": userID,
		},
		Edges: []graph.Edge{{DestinationTable: &users, DestinationCol: &id, OriginCol: &userID}},
	},
},
}
