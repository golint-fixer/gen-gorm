package backends

import "github.com/kmulvey/gen-gorm/graph"

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
