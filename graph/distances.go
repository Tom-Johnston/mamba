package graph

import (
	"container/list"

	"github.com/Tom-Johnston/gigraph/ints"
)

//Distance returns the length of the shortest path from i to j in g and -1 if there is no path
//It finds this by doing a BFS search.
func Distance(g Graph, i, j int) int {
	if i == j {
		return 0
	}
	n := g.N()
	distances := make([]int, n)
	verticesToCheck := list.New()
	verticesToCheck.PushBack(i)
	for verticesToCheck.Len() > 0 {
		k := verticesToCheck.Remove(verticesToCheck.Front()).(int)
		for _, v := range g.Neighbours(k) {
			if distances[v] == 0 {
				if v == j {
					return distances[k] + 1
				}
				distances[v] = distances[k] + 1
				verticesToCheck.PushBack(v)
			}
		}
	}
	return -1
}

//Eccentricity returns a slice giving the eccentricity of each vertex and a slice of -1s if the graph is disconnected.
//The eccentricity of a vertex v is the maximum over vertices u of the length of the shortest path from v to u.
func Eccentricity(g Graph) []int {
	n := g.N()
	eccentricity := make([]int, n)
	distances := make([]int, n)
	for i := 0; i < n; i++ {
		e := 0
		seenVertices := 0
		if i != 0 {
			for j := range distances {
				distances[j] = 0
			}
		}
		verticesToCheck := list.New()
		verticesToCheck.PushBack(i)
		for verticesToCheck.Len() > 0 {
			k := verticesToCheck.Remove(verticesToCheck.Front()).(int)
			for _, l := range g.Neighbours(k) {
				if l != i && distances[l] == 0 {
					seenVertices++
					tmp := distances[k] + 1
					distances[l] = tmp
					if tmp > e {
						e = tmp
					}
					verticesToCheck.PushBack(l)
				}
			}
		}
		if seenVertices == n-1 {
			eccentricity[i] = e
		} else {
			eccentricity[i] = -1
		}
	}
	return eccentricity
}

//Diameter retuns the diameter of a graph or -1 for a disconnected graph.
//The diameter is the maximum eccentricity of the vertices of the graph i.e. it is the length of the longest shortest path in g.
func Diameter(g Graph) int {
	if g.N() == 0 {
		return 0
	}
	e := Eccentricity(g)
	if ints.Min(e) == -1 {
		return -1
	}
	return ints.Max(Eccentricity(g))
}

//Radius returns the raidus of a graph or -1 for a disconnected graph.
//The radius is the minimum eccentricity of the vertices of the graph.
func Radius(g Graph) int {
	if g.N() == 0 {
		return 0
	}
	return ints.Min(Eccentricity(g))
}

//Girth returns the size of the shortest cycle in g.
//It finds the shortest cycle by using a BFS algorithm to find the shortest cycle containing each vertex v. This is probably not the most efficient way.
//Directed graps are currently not supported.
func Girth(g Graph) int {
	n := g.N()

	//Compute the girth using a BFS algorithm starting from each vertex except the last 2.
	if n < 3 {
		return -1
	}
	girth := n + 2
	distances := make([]int, n)
	parentVertices := make([]int, n)
	verticesToCheck := list.New()
	for i := 0; i < n-2; i++ { //Note that starting from the last 2 vertices is a waste of time because the minimum cycle must include another of the vertices.
		zeroOut(distances)
		verticesToCheck.PushBack(i)
		for verticesToCheck.Len() > 0 {
			k := verticesToCheck.Remove(verticesToCheck.Front()).(int)
			for _, j := range g.Neighbours(k) {
				if j != parentVertices[k] {
					if j == i && distances[k]+1 < girth {
						girth = distances[k] + 1
					} else if j != i && distances[j] == 0 {
						if distances[k]+2 < girth {
							//This vertex hasn't already been encountered.
							parentVertices[j] = k
							distances[j] = distances[k] + 1
							verticesToCheck.PushBack(j)
						}
					} else if j != i && distances[k]+distances[j]+1 < girth {
						girth = distances[k] + distances[j] + 1
					}
				}
			}
		}
	}
	if girth == n+2 {
		return -1
	}
	return girth
}
