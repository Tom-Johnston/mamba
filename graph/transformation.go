package graph

import "github.com/Tom-Johnston/gigraph/sortints"

//SplitEdge modifies the graph G by removing the edge ij (if it is present) and adding a new vertex connected to i and j.
func SplitEdge(g EditableGraph, i, j int) {
	//Make j > i
	if j == i {
		panic("Multiedges are not supported")
	}

	g.RemoveEdge(i, j)
	g.AddVertex([]int{i, j})
}

//Contract modifies g by adding edges from i to all the neighbours of j and then removing the vertex j.
//The main use case will be contracting an edge although this doesn't require the edge (i,j) to be present. One may wish to consider the order of i and j for performance reasons e.g. in a DenseGraph removing the vertex of larger index requires fewer copies and is likely to be quicker.
func Contract(g EditableGraph, i, j int) {
	for _, v := range g.Neighbours(j) {
		g.AddEdge(i, v)
	}
	g.RemoveVertex(j)
}

type complement struct {
	g Graph
}

func (c complement) N() int {
	return c.g.N()
}

func (c complement) M() int {
	n := c.g.N()
	return (n*(n-1))/2 - c.g.M()
}

func (c complement) IsEdge(i, j int) bool {
	return !c.g.IsEdge(i, j)
}

//Neighbours returns the neighbours of the vertex v.
func (c complement) Neighbours(v int) []int {
	n := c.g.N()
	neighbours := sortints.Complement(n, c.g.Neighbours(v))
	neighbours.Remove(v)
	return neighbours
}

//Degrees returns the degree sequence of g.
func (c complement) Degrees() []int {
	n := c.g.N()
	degrees := c.g.Degrees()
	for i := range degrees {
		degrees[i] = (n - 1) - degrees[i]
	}
	return degrees
}

//Complement returns the complement of a graph. Updating the original graph, changes the complement.
func Complement(g Graph) Graph {
	return complement{g: g}
}

//ComplementDense returns the graph on the same vertices as g where an edge is present if and only if it isn't present in g.
func ComplementDense(g Graph) *DenseGraph {
	n := g.N()
	m := n*(n-1)/2 - g.M()
	oldDegrees := g.Degrees()
	degrees := make([]int, n)
	for i := range degrees {
		degrees[i] = (n - 1) - oldDegrees[i]
	}
	edges := make([]byte, (n*(n-1))/2)
	index := 0
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if !g.IsEdge(i, j) {
				edges[index] = 1
			}
			index++
		}
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: m, DegreeSequence: degrees, Edges: edges}
}

//LineGraphDense returns a graph L which has the edges of g as its vertices and two vertices are joined in L iff the corresponding edges in g share a vertex.
func LineGraphDense(g Graph) *DenseGraph {
	m := g.M()
	edges := make([]byte, (m*(m-1))/2)
	lVerticesLower := make([]int, 0, m)
	lVerticesUpper := make([]int, 0, m)
	mIndex := 0
	for j := 0; j < g.N(); j++ {
		for i := 0; i < j; i++ {
			if g.IsEdge(i, j) {
				for k, v := range lVerticesLower {
					if i == v {
						edges[(mIndex*(mIndex-1))/2+k] = 1
					}
				}
				for k, v := range lVerticesUpper {
					if i == v {
						edges[(mIndex*(mIndex-1))/2+k] = 1
					} else if i < v {
						break
					}
				}

				for k := len(lVerticesUpper) - 1; k >= 0; k-- {
					if lVerticesUpper[k] == j {
						edges[(mIndex*(mIndex-1))/2+k] = 1
					} else {
						break
					}
				}
				lVerticesLower = append(lVerticesLower, i)
				lVerticesUpper = append(lVerticesUpper, j)
				mIndex++
			}
		}
	}
	return NewDense(m, edges)
}
