package graph

//Graph is the interface that represents a graph which cannot be copied or edited.
type Graph interface {
	//N() returns the number of vertices in the graph.
	N() int
	//M() returns the number of edges in the graph.
	M() int

	//IsEdge returns true if the edge is in the graph and false otherwise. It may panic if i or j are not in the appropriate range (although one could argue the correct response is really false).
	IsEdge(i, j int) bool
	//Neighbours returns the neighbours of the vertex v.
	Neighbours(v int) []int
	//Degrees returns the degree sequence of g.
	Degrees() []int
}

//EditableGraph is the interface which represent a graph which can be copied and edited.
type EditableGraph interface {
	Graph

	//AddVertex modifies the graph by appending one new vertex with edges from the new vertex to the vertices in neighbours.
	AddVertex(neighbours []int)
	//RemoveVertex modifies the graph by removing the speicified vertex. The index of a vertex u > v becomes u - 1 while the index of u < v is unchanged.
	RemoveVertex(i int)
	//AddEdge modifies the graph by adding the edge (i, j) if it is not already present.
	AddEdge(i, j int)
	//RemoveEdge modifies the graph by removing the edge (i, j) if it is present.
	RemoveEdge(i, j int)

	//InducedSubgraph returns a deep copy of the induced subgraph of g with vertices given in order by V. This must not modify g.
	InducedSubgraph(V []int) EditableGraph
	//Copy returns a deep copy of the graph g.
	Copy() EditableGraph
}
