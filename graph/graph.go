package graph

// Graph is a databse
type Graph struct {
	Vertices map[string]*Vertex
	Name     string
}

// Vertex is a table
type Vertex struct {
	Name    string
	Cols    map[string]Col
	HasMany string
	Edges   []Edge
}

// Col is a column ... duh?!?
type Col struct {
	Name string
	Type string
	Key  string
}

// Edge is a foreign key
type Edge struct {
	DestinationTable *Vertex
	DestinationCol   *Col
	OriginCol        *Col
}
