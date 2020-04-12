package graph

import (
	"sort"

	"github.com/Tom-Johnston/mamba/ints"
	"github.com/Tom-Johnston/mamba/sortints"
)

//NumberOfCycles returns a slice where the ith element contains the number of cycles of length i.
//Any cycle is contained in a biconnected component so the algorithm first splits the graph into biconnected components. The algorithm involves finding a spanning tree and this implementation doesn't check it finds every vertex so splitting into at least components is necessary. Spliting into biconnected components should help prevent unnecessary XORing when checking all the cycles.
//We find a set of fundamental cycles according to the paper K. Paton, An algorithm for finding a fundamental set of cycles for an undirected linear graph, Comm. ACM 12 (1969), pp. 514-518. We can then find all cycles XORing together every combination of fundamental cycles and ignoring ones which are made of copies of 2 or more disjoing cycles. This is done according to Gibb's Algorithm from ALGORITHMIC APPROACHES TO CIRCUIT ENUMERATION PROBLEMS AND APPLICATIONS by Boon Chai Lee avaible here: http://dspace.mit.edu/bitstream/handle/1721.1/68106/FTL_R_1982_07.pdf.
//This effectively finds every cycle in the graph and could be adapted to output every cycle if required. Remember to switch from the labels in the biconnected component to the labels of the graph and that the edges are not stored in the order in the cycle.
func NumberOfCycles(g EditableGraph) []int {
	n := g.N()
	numberFound := make([]int, n+1)
	if n == 0 {
		return numberFound
	}

	//Find the biconnected components and then perform the count on each one.
	bicoms, _ := BiconnectedComponents(g)
	for _, bicom := range bicoms {
		a := g.InducedSubgraph(bicom)
		n = a.N()
		if n < 3 {
			//Can't have any cycles with fewer than 3 vertices.
			continue
		}
		h := a.Copy() //This is the working copy.

		//The fundamental cycles. They will be stored as a list of edges where the edge (i,j) with i < j is encoded as (j*(j-1))/2 + i.
		fundCycles := make([][]int, 0, 1)

		T := make([]int, n) //T[u] will store the parent vertex of each vertex in the tree (except the root 0). If the vertex v is not in the tree T[v] = -1.
		for i := 1; i < n; i++ {
			T[i] = -1
		}
		depth := make([]int, n) //Depth of the vertices in the tree. No need to set the value for vertices not in the tree as the check is always done on T.

		X := make([]int, 1, n) //This holds what the paper calls T intersection X, the vertices not yet examined and in the tree.
		var v int
		for len(X) > 0 {
			X, v = X[:len(X)-1], X[len(X)-1]
			for _, u := range h.Neighbours(v) {
				if T[u] != -1 {
					length := depth[v] - depth[T[u]] + 2 //As noted in the paper the back edge leads to something distance exactly one from the path to v.
					cycle := make([]int, length)
					//To make it easy to backtrack given the array of parents the cycle is given as the parent of u, u, v,... parent of u.
					if T[u] < u {
						cycle[0] = (u*(u-1))/2 + T[u]
					} else {
						cycle[0] = (T[u]*(T[u]-1))/2 + u
					}

					if u < v {
						cycle[1] = (v*(v-1))/2 + u
					} else {
						cycle[1] = (u*(u-1))/2 + v
					}

					previous := v
					for i := 2; i < length; i++ {
						if previous < T[previous] {
							cycle[i] = (T[previous]*(T[previous]-1))/2 + previous
						} else {
							cycle[i] = (previous*(previous-1))/2 + T[previous]
						}
						previous = T[previous]
					}
					sort.Ints(cycle) //We sort the edges of the cycle according to their encoding to make it easy to XOR.
					fundCycles = append(fundCycles, cycle)
				} else {
					T[u] = v
					X = append(X, u)
					depth[u] = depth[v] + 1
				}
				h.RemoveEdge(u, v)
			}
		}
		if len(fundCycles) == 0 {
			//No need to do anything.
			continue
		}

		//Now find all cycles from the fundamental cycles using Gibb's algorithm.
		S := [][]int{fundCycles[0]}
		Q := [][]int{fundCycles[0]}
		R := [][]int{}
		P := [][]int{} //This is what the reference calls R* but that isn't convenient for programming.

		var V []int
		for i := 1; i < len(fundCycles); i++ {
			//Step 2
			for _, t := range Q {
				tmp := sortints.XOR(t, fundCycles[i])

				if len(tmp) != len(t)+len(fundCycles[i]) {
					//They have some intersection
					R = append(R, tmp)
				} else {
					P = append(P, tmp)
				}
				Q = append(Q, tmp)
			}
			//Step 3
			for j := len(R) - 1; j >= 0; j-- {
				V = R[j]
				for k := 0; k < len(R); k++ {
					if k == j {
						continue
					}
					if sortints.ContainsSorted(V, R[k]) {
						R[j] = R[len(R)-1]
						R = R[:len(R)-1]
						P = append(P, V)
						break
					}
				}
			}

			//Step 4
			S = append(S, R...) //TODO: Remove this and replace it with just counting the lengths?
			S = append(S, fundCycles[i])
			Q = append(Q, fundCycles[i])
			R = R[:0]
			P = P[:0]
		}

		//We have now found every cycle and we check the length of each one.
		for _, V = range S {
			numberFound[len(V)]++
		}
	}

	return numberFound
}

//NumberOfInducedPaths returns a slice of length n containing the number of induced paths in g which are of length at most k.
//The length of a path is the number of edges in the path.
func NumberOfInducedPaths(g Graph, maxLength int) []int {
	n := g.N()
	if maxLength < 0 || maxLength > n-1 {
		maxLength = n - 1
	}
	r := make([]int, n)
	type path struct {
		p                []int
		length           int
		bannedNeighbours sortints.SortedInts
	}
	com := ConnectedComponents(g)
	for _, v := range com {
		h := InducedSubgraph(g, v)
		n := h.N()
		for i := 0; i < n; i++ {
			var p path
			toCheck := make([]path, 1)
			toCheck[0] = path{[]int{i}, 0, []int{i}} //Look for paths in h starting at i.
			for len(toCheck) > 0 {
				p, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]

				options := sortints.SetMinus(h.Neighbours(p.p[len(p.p)-1]), p.bannedNeighbours)

				r[p.length+1] += len(options)

				if p.length >= maxLength-1 {
					continue
				}

				for _, v := range options {
					tmpP := make([]int, p.length+2)
					copy(tmpP, p.p)
					tmpP[p.length+1] = v
					tmpBannedNeighbours := sortints.Union(p.bannedNeighbours, h.Neighbours(p.p[len(p.p)-1]))
					tmpBannedNeighbours.Add(v)
					toCheck = append(toCheck, path{tmpP, p.length + 1, tmpBannedNeighbours})
				}
			}
		}
	}
	for i := 1; i < len(r); i++ {
		r[i] /= 2
	}
	r[0] = n
	return r
}

//NumberOfInducedCycles returns a slice of length n containing the number of induced cycles in g which are of length at most k.
//The length of a cycle is the number of edges in the cycle, or equivalently, the number of vertices in the cycle.
func NumberOfInducedCycles(g Graph, maxLength int) []int {
	n := g.N()
	if maxLength < 0 || maxLength > n {
		maxLength = n
	}
	r := make([]int, n+1)
	type cycle struct {
		p                []int
		length           int
		allowedEnds      sortints.SortedInts
		bannedNeighbours sortints.SortedInts
	}
	com := ConnectedComponents(g)
	for _, v := range com {
		h := InducedSubgraph(g, v)
		n := h.N()
		for i := 0; i < n; i++ {
			var p cycle
			toCheck := make([]cycle, 1)
			toCheck[0] = cycle{[]int{i}, 0, h.Neighbours(i), []int{i}} //Look for paths in h starting at i.
			for len(toCheck) > 0 {
				p, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]

				if p.length > 0 {
					numCycles := len(sortints.Intersection(h.Neighbours(p.p[len(p.p)-1]), p.allowedEnds))
					r[p.length+2] += numCycles
				}

				if p.length >= maxLength-2 {
					continue
				}

				options := sortints.SetMinus(h.Neighbours(p.p[len(p.p)-1]), p.bannedNeighbours)
				for _, v := range options {
					tmpP := make([]int, p.length+2)
					copy(tmpP, p.p)
					tmpP[p.length+1] = v
					tmpBannedNeighbours := sortints.Union(p.bannedNeighbours, h.Neighbours(p.p[len(p.p)-1]))
					tmpBannedNeighbours.Add(v)
					var tmpAllowedEnds sortints.SortedInts
					if p.length > 0 {
						tmpAllowedEnds = sortints.SetMinus(p.allowedEnds, h.Neighbours(p.p[len(p.p)-1]))
					} else {
						tmpAllowedEnds = make([]int, len(p.allowedEnds))
						copy(tmpAllowedEnds, p.allowedEnds)
						tmpAllowedEnds.Remove(v)
					}

					toCheck = append(toCheck, cycle{tmpP, p.length + 1, tmpAllowedEnds, tmpBannedNeighbours})
				}
			}
		}
	}
	for i := 1; i < len(r); i++ {
		r[i] /= 2 * i
	}
	return r
}

type inducedSubgraph struct {
	verts   []int
	sortedV []int
	indices []int
	g       Graph
}

func (h inducedSubgraph) N() int {
	return len(h.verts)
}

func (h inducedSubgraph) Degrees() []int {
	degrees := make([]int, len(h.verts))
	for i, v := range h.verts {
		degrees[i] = sortints.IntersectionSize(h.g.Neighbours(v), h.sortedV)
	}
	return degrees
}

func (h inducedSubgraph) M() int {
	degrees := h.Degrees()
	sum := ints.Sum(degrees)
	return sum / 2
}

func (h inducedSubgraph) IsEdge(i, j int) bool {
	return h.g.IsEdge(h.verts[i], h.verts[j])
}

func (h inducedSubgraph) Neighbours(v int) []int {
	return intersectionByIndex(h.g.Neighbours(h.verts[v]), h.sortedV, h.indices)
}

//InducedSubgraph returns a graph which represents the subgraph of g induced by the vertices in V in the order they are in V.
//The properties of the induced subgraph are calculated from g when called and reflect the current state of g. If a vertex in V is no longer in the graph, the behaviour of this function is unspecified.
func InducedSubgraph(g Graph, V []int) Graph {
	values, indices := intsSort(V)
	return inducedSubgraph{verts: V, sortedV: values, indices: indices, g: g}
}

type sortWithIndex struct {
	values  *[]int
	indices *[]int
}

func (s sortWithIndex) Len() int { return len(*s.values) }
func (s sortWithIndex) Swap(i, j int) {
	(*s.values)[i], (*s.values)[j] = (*s.values)[j], (*s.values)[i]
	(*s.indices)[i], (*s.indices)[j] = (*s.indices)[j], (*s.indices)[i]
}
func (s sortWithIndex) Less(i, j int) bool { return (*s.values)[i] < (*s.values)[j] }

func intsSort(a []int) (values, indices []int) {
	values = make([]int, len(a))
	copy(values, a)
	indices = make([]int, len(a))
	for i := range indices {
		indices[i] = i
	}

	toSort := sortWithIndex{values: &values, indices: &indices}
	sort.Sort(toSort)
	return values, indices
}

func intersectionByIndex(a, b sortints.SortedInts, indicesOfB []int) sortints.SortedInts {
	rV := make([]int, 0, sortints.IntersectionSize(a, b))
	r := sortints.SortedInts(rV)
	i := 0 //Point in a
	j := 0 //Point in b
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			r.Add(indicesOfB[j])
			i++
			j++
		} else if a[i] > b[j] {
			j++
		} else {
			i++
		}
	}
	return r
}
