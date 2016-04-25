package graph

type Graph struct {
	Vertices []Vertex
}

type Vertex struct {
	Name  string
	Edges []Vertex
}
