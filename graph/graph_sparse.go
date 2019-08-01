package graph

import (
	"sort"

	"github.com/Tom-Johnston/gigraph/sortints"
)

//SparseGraph is a data structure for representing a simple undirected graph. *SparseGraph implements the graph insterface.
//SparseGraph stores the number of vertices, the number of edges, the degree sequence of the graph and the neighbourhood of edge vertex.
//Most modifications are slow as the edges are stored as SortedInts and not a data structure with log(n) modifications but sparse graphs take up relatively little space. Checking if an individual edge is present is logarithmic in the degree of the vertices but returning the neighbours is quick so most algorithms run fairly quickly.
//TODO Do we want to switch the neighbourhoods to using some kind of heap with quicker insertions and deletions?
type SparseGraph struct {
	NumberOfVertices int
	NumberOfEdges    int

	Neighbourhoods []sortints.SortedInts
	DegreeSequence []int
}

//NewSparse create the graph on n vertices with the specified neighbours. If neighbourhoods is nil, the empty graph is created.
//TODO Check that this is a graph. Check SortedInts are in fact sorted.
func NewSparse(n int, neighbourhoods []sortints.SortedInts) *SparseGraph {
	if neighbourhoods == nil {
		neighbourhoods = make([]sortints.SortedInts, n)
		for i := range neighbourhoods {
			neighbourhoods[i] = []int{}
		}
	}

	if len(neighbourhoods) != n {
		panic("Length of neighbourhoods doesn't match the number of vertices.")
	}

	tmpNeighbourhoods := make([]sortints.SortedInts, n)

	for i := range neighbourhoods {
		tmpNeighbourhoods[i] = sortints.NewSortedInts(neighbourhoods[i]...)
	}

	degreeSequence := make([]int, n)
	for i := range tmpNeighbourhoods {
		degreeSequence[i] = len(tmpNeighbourhoods[i])
	}

	return &SparseGraph{NumberOfVertices: n, NumberOfEdges: intsSum(degreeSequence) / 2, Neighbourhoods: tmpNeighbourhoods, DegreeSequence: degreeSequence}

}

//N returns the number of vertices in g.
func (g SparseGraph) N() int {
	return g.NumberOfVertices
}

//M returns the number of edges in g.
func (g SparseGraph) M() int {
	return g.NumberOfEdges
}

//IsEdge returns if there is an edge between i and j in the graph.
func (g SparseGraph) IsEdge(i, j int) bool {
	if g.DegreeSequence[i] > g.DegreeSequence[j] {
		return sortints.ContainsSingle(g.Neighbourhoods[i], j)
	}
	return sortints.ContainsSingle(g.Neighbourhoods[j], i)
}

//Neighbours returns the neighbours of the given vertex.
func (g SparseGraph) Neighbours(v int) []int {
	tmpNeighbours := make([]int, len(g.Neighbourhoods[v]))
	copy(tmpNeighbours, g.Neighbourhoods[v])
	return tmpNeighbours
}

//Degrees returns the degree sequence of the graph.
func (g SparseGraph) Degrees() []int {
	tmpDegreeSequence := make([]int, len(g.DegreeSequence))
	copy(tmpDegreeSequence, g.DegreeSequence)
	return tmpDegreeSequence
}

//AddVertex adds a vertex to the graph with the specified neighbours.
func (g *SparseGraph) AddVertex(neighbours []int) {
	g.NumberOfVertices++
	tmp := sortints.NewSortedInts(neighbours...)

	g.NumberOfEdges += len(tmp)

	for _, v := range tmp {
		g.Neighbourhoods[v] = append(g.Neighbourhoods[v], g.NumberOfVertices-1)
		g.DegreeSequence[v]++
	}

	g.Neighbourhoods = append(g.Neighbourhoods, sortints.SortedInts(tmp))
	g.DegreeSequence = append(g.DegreeSequence, len(neighbours))

}

//RemoveVertex removes the specified vertex. The index of a vertex u > v becomes u - 1 while the index of u < v is unchanged.
func (g *SparseGraph) RemoveVertex(i int) {
	g.NumberOfVertices--
	g.NumberOfEdges -= g.DegreeSequence[i]

	for _, v := range g.Neighbourhoods[i] {
		g.Neighbourhoods[v].Remove(i)
	}

	g.Neighbourhoods = g.Neighbourhoods[:i+copy(g.Neighbourhoods[i:], g.Neighbourhoods[i+1:])]

	for j := range g.Neighbourhoods {
		startIndex := sort.SearchInts(g.Neighbourhoods[j], i)
		for k := startIndex; k < len(g.Neighbourhoods[j]); k++ {
			g.Neighbourhoods[j][k]--
		}
	}

}

//AddEdge modifies the graph by adding the edge (i, j) if it is not already present.
//If the edge is already present (or i == j), this does nothing.
func (g *SparseGraph) AddEdge(i, j int) {
	if !g.IsEdge(i, j) {
		g.Neighbourhoods[i].Add(j)
		g.Neighbourhoods[j].Add(i)
		g.NumberOfEdges++
		g.DegreeSequence[i]++
		g.DegreeSequence[j]++
	}
}

//RemoveEdge modifies the graph by removing the edge (i, j) if it is present.
//If the edge is not already present, this does nothing.
func (g *SparseGraph) RemoveEdge(i, j int) {
	if g.IsEdge(i, j) {
		g.Neighbourhoods[i].Remove(j)
		g.Neighbourhoods[j].Remove(i)
		g.NumberOfEdges--
		g.DegreeSequence[i]--
		g.DegreeSequence[j]--
	}
}

//InducedSubgraph returns a deep copy of the induced subgraph of g with vertices given in order by V.
//This can also be used to return relabellings of the graph if len(V) = g.N().
func (g SparseGraph) InducedSubgraph(V []int) EditableGraph {
	n := len(V)
	values, indices := intsSort(V)
	tmpDegreeSequence := make([]int, n)
	tmpNeighbourhoods := make([]sortints.SortedInts, n)
	m := 0
	for i, v := range V {
		tmpNeighbourhoods[i] = intersectionByIndex(g.Neighbours(v), values, indices)
		tmpDegreeSequence[i] = len(tmpNeighbourhoods[i])
		m += tmpDegreeSequence[i]
	}

	h := &SparseGraph{NumberOfVertices: n, NumberOfEdges: m / 2, Neighbourhoods: tmpNeighbourhoods, DegreeSequence: tmpDegreeSequence}
	return h
}

//Copy returns a deep copy of the graph g.
func (g SparseGraph) Copy() EditableGraph {
	tmpNeighbourhoods := make([]sortints.SortedInts, len(g.Neighbourhoods))
	for i := range g.Neighbourhoods {
		tmpNeighbourhoods[i] = make(sortints.SortedInts, len(g.Neighbourhoods[i]))
		copy(tmpNeighbourhoods[i], g.Neighbourhoods[i])
	}
	tmpDegreeSequence := make([]int, len(g.DegreeSequence))
	copy(tmpDegreeSequence, g.DegreeSequence)

	return &SparseGraph{NumberOfVertices: g.NumberOfVertices, NumberOfEdges: g.NumberOfEdges, Neighbourhoods: tmpNeighbourhoods, DegreeSequence: tmpDegreeSequence}
}
