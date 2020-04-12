package graph

import (
	"sort"

	"github.com/Tom-Johnston/mamba/ints"
)

//ConnectedComponent returns the connected component in g containing v.
func ConnectedComponent(g Graph, v int) []int {
	toCheck := make([]int, 1, g.N()-1)
	toCheck[0] = v
	unseen := make([]int, g.N())
	for i := range unseen {
		unseen[i] = i
	}
	unseen[v] = unseen[len(unseen)-1]
	unseen = unseen[:len(unseen)-1]
	seen := make([]int, 1, g.N())
	seen[0] = v
	u := 0
	w := 0
	for len(toCheck) > 0 {
		toCheck, u = toCheck[:len(toCheck)-1], toCheck[len(toCheck)-1]
		for i := len(unseen) - 1; i >= 0; i-- {
			w = unseen[i]
			if g.IsEdge(u, w) {
				unseen[i] = unseen[len(unseen)-1]
				unseen = unseen[:len(unseen)-1]
				toCheck = append(toCheck, w)
				seen = append(seen, w)
			}
		}
	}
	sort.Ints(seen)
	return seen
}

//ConnectedComponents returns all the connected components of g.
func ConnectedComponents(g Graph) [][]int {
	if g.N() == 0 {
		return [][]int{}
	} else if g.N() == 1 {
		return [][]int{[]int{0}}
	}
	components := make([][]int, 0, 1)

	toCheck := make([]int, 0, g.N()-1)
	seen := make([]int, 1, g.N())

	unseen := make([]int, g.N())
	for i := range unseen {
		unseen[i] = i
	}

	for len(unseen) > 0 {
		v := unseen[len(unseen)-1]
		unseen = unseen[:len(unseen)-1]
		toCheck = toCheck[:1]
		toCheck[0] = v
		seen = seen[:1]
		seen[0] = v
		u := 0
		w := 0
		for len(toCheck) > 0 {
			toCheck, u = toCheck[:len(toCheck)-1], toCheck[len(toCheck)-1]
			for i := len(unseen) - 1; i >= 0; i-- {
				w = unseen[i]
				if g.IsEdge(u, w) {
					unseen[i] = unseen[len(unseen)-1]
					unseen = unseen[:len(unseen)-1]
					toCheck = append(toCheck, w)
					seen = append(seen, w)
				}
			}
		}
		sort.Ints(seen)
		tmp := make([]int, len(seen))
		copy(tmp, seen)
		components = append(components, tmp)
	}
	return components
}

//BiconnectedComponents returns the biconnected components and the articulation vertices.
func BiconnectedComponents(g Graph) ([][]int, []int) {
	articulationPoints := make([]int, 0)
	biconnectedComponents := make([][]int, 0)
	//Split into connected components.
	components := ConnectedComponents(g)
	for _, com := range components {
		h := InducedSubgraph(g, com)
		n := h.N()
		toCheck := make([]int, 1, n)
		toCheck[0] = 0
		depths := make([]int, n)
		for i := 1; i < n; i++ {
			depths[i] = -1
		}
		lowpoints := make([]int, n)
		parents := make([]int, n)
		isArticulation := make([]bool, n)

		childCount := 0
		bicoms := make([][]int, 1, 1)
		tmp := make([]int, 0, n)
		bicoms[0] = tmp
	DFS:
		for len(toCheck) > 0 {
			v := toCheck[len(toCheck)-1]
			// fmt.Println("v", v)
			tmpLowPoint := lowpoints[v]
			for _, u := range h.Neighbours(v) {
				// fmt.Println("u", u)
				// fmt.Println(lowpoints)
				// fmt.Println(depths)
				// fmt.Println(isArticulation)
				// fmt.Println("bicoms", bicoms)
				// fmt.Println("childCount", childCount)
				if depths[u] == -1 {
					if v == 0 {
						childCount++
					}
					toCheck = append(toCheck, u)
					depths[u] = depths[v] + 1
					lowpoints[u] = depths[v] + 1
					parents[u] = v
					if len(bicoms[len(bicoms)-1]) > 0 {
						tmp = make([]int, 0, n)
						bicoms = append(bicoms, tmp)
					}
					continue DFS
				} else if u != parents[v] {
					if lowpoints[u] < tmpLowPoint {
						tmpLowPoint = lowpoints[u]
					}
					if v != 0 && parents[u] == v && lowpoints[u] >= depths[v] {
						// fmt.Println("Add")
						// fmt.Println(lowpoints)
						parents[u] = -1
						bicoms[len(bicoms)-1] = append(bicoms[len(bicoms)-1], v)
						for i := range bicoms[len(bicoms)-1] {
							bicoms[len(bicoms)-1][i] = com[bicoms[len(bicoms)-1][i]]
						}
						sort.Ints(bicoms[len(bicoms)-1])
						biconnectedComponents = append(biconnectedComponents, bicoms[len(bicoms)-1])
						bicoms[len(bicoms)-1] = make([]int, 0, n)
						isArticulation[v] = true
					}
				}
			}
			lowpoints[v] = tmpLowPoint
			toCheck = toCheck[:len(toCheck)-1]
			// fmt.Println("bicoms", bicoms)
			if v != 0 {
				for i := len(bicoms) - 2; i >= -1; i-- {
					// fmt.Println(bicoms, i, v)

					if i > -1 && depths[bicoms[i][len(bicoms[i])-1]] == depths[v]+1 {
						//fmt.Println(depths[bicoms[i][len(bicoms[i])-1]])
						bicoms[len(bicoms)-1] = append(bicoms[len(bicoms)-1], bicoms[i]...)
					} else {
						bicoms[i+1] = bicoms[len(bicoms)-1]
						bicoms = bicoms[:i+2]
						break
					}
				}
			}

			bicoms[len(bicoms)-1] = append(bicoms[len(bicoms)-1], v)
			// fmt.Println("bicoms", bicoms)
			// fmt.Println(biconnectedComponents)
		}

		for i := 0; i < len(bicoms)-1; i++ {
			bicoms[i] = append(bicoms[i], 0)
		}

		for i := range bicoms {
			for j := range bicoms[i] {
				bicoms[i][j] = com[bicoms[i][j]]
			}
			sort.Ints(bicoms[i])
		}

		biconnectedComponents = append(biconnectedComponents, bicoms...)

		if childCount < 2 {
			//fmt.Println("here2")
			isArticulation[0] = false
		} else {
			//fmt.Println("here")
			isArticulation[0] = true
		}
		for i, b := range isArticulation {
			if b {
				articulationPoints = append(articulationPoints, com[i])
			}
		}
	}
	return biconnectedComponents, articulationPoints
}

//MinDegree returns the minimum degree of g.
func MinDegree(g Graph) int {
	minDegree := g.N()
	for _, v := range g.Degrees() {
		if v < minDegree {
			minDegree = v
		}
	}
	return minDegree
}

//MaxDegree returns the maximum degree of g.
func MaxDegree(g Graph) int {
	maxDegree := 0
	for _, v := range g.Degrees() {
		if v > maxDegree {
			maxDegree = v
		}
	}
	return maxDegree
}

//Equal returns true if the two lablled graphs are exactly equal and false otherwise.
//To test is two graphs are isomorphic the graphs need to be transformed into their canonical isomorphs first (e.g. by using g.InducedSubgraph(g.CanonicalIsomorph()))
func Equal(g, h Graph) bool {
	if g.N() != h.N() {
		return false
	}

	n := g.N()
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if g.IsEdge(i, j) != h.IsEdge(i, j) {
				return false
			}
		}
	}
	return true
}

//Degeneracy returns the smallest integer d such that every ordering of the vertices contains a vertex preceeded by at least d neighbours. It also returns an ordering where no vertex is proceeded by d + 1 neighbours.
func Degeneracy(g Graph) (d int, order []int) {
	//Extract the information.
	n := g.N()
	if n == 0 {
		return 0, nil
	}
	degreeSequence := g.Degrees()
	maxDegree := ints.Max(degreeSequence)

	//Initialise the degeneracy and an optimum ordering.
	d = 0
	order = make([]int, n)

	//Initialise the bins and an array keeping track of the new degrees.
	bins := make([][]int, maxDegree+1)
	for i := range bins {
		bins[i] = make([]int, 0)
	}
	for i, v := range degreeSequence {
		bins[v] = append(bins[v], i)
	}
	degrees := make([]int, n)
	copy(degrees, degreeSequence)

	//Repeatedly remove a vertex with the fewest neighbours not in the list.
	for i := 0; i < n; i++ {
		//Find the first non-empty to bin.
		var j int
		for j = range bins {
			if len(bins[j]) != 0 {
				break
			}
		}
		//Increase the degeneracy if necessary.
		if j > d {
			d = j
		}

		//Prepend the vertex to the order.
		v := bins[j][len(bins[j])-1]
		order[n-1-i] = v
		//Remove the vertex from the bins and set it to -1 in degrees.
		bins[j] = bins[j][:len(bins[j])-1]
		degrees[v] = -1
		//Update the neighbours of v.
		neighbours := g.Neighbours(v)
		for _, u := range neighbours {
			if degrees[u] == -1 {
				continue
			}

			for k, w := range bins[degrees[u]] {
				if w != u {
					continue
				}
				bins[degrees[u]][k] = bins[degrees[u]][len(bins[degrees[u]])-1]
				bins[degrees[u]] = bins[degrees[u]][:len(bins[degrees[u]])-1]
				degrees[u]--
				bins[degrees[u]] = append(bins[degrees[u]], u)
				break
			}
		}
	}
	return d, order
}
