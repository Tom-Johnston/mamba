package graph

//DenseGraph is a data structure representing a simple undirected labelled graph.
//DenseGraph stores the number of vertices, the number of edges, the degree sequence of the graph and stores the edges in a []byte array which has an indicator of an edge being present. The edges are in the order 01, 02, 12, 03, 13, 23... so the edge ij with i < j is in the (j*(j-1))/2 + i place.
//Adding or removing edges are quick operations. Adding a vertex may be quick if the backing array doesn't have to grow but may require copying the entire adjacency matrix. Removing a vertex is generally slow.
//*DenseGraph implements the Graph interface.
type DenseGraph struct {
	NumberOfVertices int
	NumberOfEdges    int
	DegreeSequence   []int
	Edges            []byte
}

//NewDense returns a pointer to a DenseGraph representation of the graph with n vertices and the edges as given in edges.
//The edges are in the order 01, 02, 12, 03, 13, 23... so the edge ij with i < j is in the (j*(j-1))/2 + i place. The *DenseGraph uses its own copy of edges and modifications to edges won't change the current graph.
//*DenseGraph implements the Graph interface.
func NewDense(n int, edges []byte) *DenseGraph {
	if edges == nil {
		edges = make([]byte, (n*(n-1))/2)
		return &DenseGraph{NumberOfVertices: n, NumberOfEdges: 0, DegreeSequence: make([]int, n), Edges: edges}
	}
	if len(edges) != (n*(n-1))/2 {
		panic("Wrong number of edges")
	}
	degrees := make([]int, n)
	m := 0
	copyOfEdges := make([]byte, len(edges))
	copy(copyOfEdges, edges)
	index := 0
	for j := 0; j < n; j++ {
		for i := 0; i < j; i++ {
			if edges[index] > 0 {
				degrees[i]++
				degrees[j]++
				m++
			}
			index++
		}
	}

	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: m, DegreeSequence: degrees, Edges: edges}
}

//N returns the number of vertices in the graph.
func (g DenseGraph) N() int {
	return g.NumberOfVertices
}

//M returns the number of vertices in the graph.
func (g DenseGraph) M() int {
	return g.NumberOfEdges
}

//IsEdge returns true if the undirected edge (i, j) is present in the graph and false otherwise.
func (g DenseGraph) IsEdge(i, j int) bool {
	if i >= g.N() || j >= g.N() || i < 0 || j < 0 {
		return false
	}
	if i < j && g.Edges[(j*(j-1))/2+i] > 0 {
		return true
	} else if i > j && g.Edges[(i*(i-1))/2+j] > 0 {
		return true
	}
	return false
}

//Neighbours returns the neighbours of v i.e. the vertices u such that (u,v) is an edge.
func (g DenseGraph) Neighbours(v int) []int {
	degrees := g.Degrees()
	r := make([]int, 0, degrees[v])
	tmp := (v * (v - 1)) / 2
	for i := 0; i < v; i++ {
		index := tmp + i
		if g.Edges[index] > 0 {
			r = append(r, i)
		}
	}

	for i := v + 1; i < g.N(); i++ {
		index := (i*(i-1))/2 + v
		if g.Edges[index] > 0 {
			r = append(r, i)
		}
	}
	return r
}

//Degrees returns the slice containing the degrees (number of edges incident with the vertex) of each vertex.
func (g DenseGraph) Degrees() []int {
	tmpDegreeSequence := make([]int, len(g.DegreeSequence))
	copy(tmpDegreeSequence, g.DegreeSequence)
	return tmpDegreeSequence
}

//AddEdge modifies the graph by adding the edge (i, j) if it is not already present.
//If the edge is already present (or i == j), this does nothing.
func (g *DenseGraph) AddEdge(i, j int) {
	if i == j || g.IsEdge(i, j) {
		return
	}
	g.DegreeSequence[i]++
	g.DegreeSequence[j]++
	g.NumberOfEdges++

	if i < j {
		g.Edges[(j*(j-1))/2+i] = 1
	} else if i > j {
		g.Edges[(i*(i-1))/2+j] = 1
	}
}

//RemoveEdge modifies the graph by removing the edge (i, j) if it is present.
//If the edge is not already present, this does nothing.
func (g *DenseGraph) RemoveEdge(i, j int) {
	if !g.IsEdge(i, j) {
		return
	}
	if i < j {
		g.Edges[(j*(j-1))/2+i] = 0
	} else if i > j {
		g.Edges[(i*(i-1))/2+j] = 0
	}
	g.DegreeSequence[i]--
	g.DegreeSequence[j]--
	g.NumberOfEdges--
}

//AddVertex modifies the graph by appending one new vertex with edges from the new vertex to the vertices in neighbours.
func (g *DenseGraph) AddVertex(neighbours []int) {
	tmp := make([]byte, g.NumberOfVertices)
	for _, v := range neighbours {
		tmp[v] = 1
		g.DegreeSequence[v]++
	}
	g.Edges = append(g.Edges, tmp...)
	g.DegreeSequence = append(g.DegreeSequence, len(neighbours))
	g.NumberOfVertices++
	g.NumberOfEdges += len(neighbours)
}

//RemoveVertex modifies the graph by removing the speicified vertex. The index of a vertex u > v becomes u - 1 while the index of u < v is unchanged.
func (g *DenseGraph) RemoveVertex(v int) {
	if v >= g.NumberOfVertices {
		panic("No such vertex")
	}

	//Update the degree sequences and number of edges.
	g.NumberOfEdges -= g.DegreeSequence[v]
	neighbours := g.Neighbours(v)
	for _, u := range neighbours {
		g.DegreeSequence[u]--
	}
	copy(g.DegreeSequence[v:], g.DegreeSequence[v+1:])
	g.DegreeSequence = g.DegreeSequence[:len(g.DegreeSequence)-1]

	//Update the backing array.
	oldIndex := (v*(v+1))/2 - 1
	newIndex := (v * (v - 1)) / 2

	for j := v + 1; j < g.NumberOfVertices; j++ {
		tmp := (j*(j-1))/2 + v
		newIndex += copy(g.Edges[newIndex:], g.Edges[oldIndex+1:tmp])
		oldIndex = tmp
	}
	copy(g.Edges[newIndex:], g.Edges[oldIndex+1:])
	g.NumberOfVertices--
	g.Edges = g.Edges[:(g.NumberOfVertices*(g.NumberOfVertices-1))/2]
}

//InducedSubgraph returns a deep copy of the induced subgraph of g with vertices given in order by V.
//This can also be used to return relabellings of the graph if len(V) = g.N().
func (g *DenseGraph) InducedSubgraph(V []int) EditableGraph {
	n := len(V)
	m := 0
	degrees := make([]int, n)
	edges := make([]byte, (n*(n-1))/2)
	index := 0
	for j := 1; j < len(V); j++ {
		for i := 0; i < j; i++ {
			if g.IsEdge(V[i], V[j]) {
				edges[index] = 1
				m++
				degrees[i]++
				degrees[j]++
			}
			index++
		}
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: m, DegreeSequence: degrees, Edges: edges}
}

//Copy returns a deep copy of the graph g.
func (g *DenseGraph) Copy() EditableGraph {
	newEdges := make([]byte, len(g.Edges))
	copy(newEdges, g.Edges)
	newDegrees := make([]int, len(g.DegreeSequence))
	copy(newDegrees, g.DegreeSequence)
	return &DenseGraph{NumberOfVertices: g.NumberOfVertices, NumberOfEdges: g.NumberOfEdges, DegreeSequence: newDegrees, Edges: newEdges}
}

//Helper functions for implementing the required functions

//String returns a human readable representation of the graph.
// func (g DenseGraph) String() string {
// 	var buffer bytes.Buffer
// 	buffer.WriteString(fmt.Sprintf("Degree: %v \n", g.NumberOfVertices))
// 	for i := 0; i < g.NumberOfVertices; i++ {
// 		for j := 0; j < g.NumberOfVertices; j++ {
// 			if j < i {
// 				buffer.WriteString("  ")
// 			} else if j == i {
// 				buffer.WriteString("0 ")
// 			} else {
// 				buffer.WriteString(fmt.Sprintf("%v ", g.Edges[(j*(j-1))/2+i]))
// 			}
// 		}
// 		buffer.WriteString("\n")
// 	}
// 	return buffer.String()
// }
