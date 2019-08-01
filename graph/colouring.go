package graph

import (
	"container/heap"
)

//IsProperColouring checks if the vertex colouring is a proper colouring of the graph g. It assumes that all colours are \geq 0 and that a colour <0 is a mistake.
//This is because the colour -1 is often used to indicate no colour.
func IsProperColouring(g Graph, colouring []int) bool {
	n := g.N()

	if colouring == nil || len(colouring) != n {
		return false
	}
	for i := 0; i < n; i++ {
		if colouring[i] < 0 {
			return false
		}
		neighbours := g.Neighbours(i)
		for _, v := range neighbours {
			if v > i {
				break
			}
			if colouring[v] == colouring[i] {
				return false
			}
		}
	}
	return true
}

//GreedyColor greedily colours the graph G colouring the vertices in the given order. That is, this reads the vertices in the given order and assigns to each vertex the minimum colour such that none of its neighbour have this colour.
func GreedyColor(g Graph, order []int) (int, []int) {
	n := g.N()
	if n != len(order) {
		panic("order does not have length  equal to g.N()")
	}
	c := make([]int, n)
	for i := range c {
		c[i] = -1
	}
	seenColours := make([]bool, n)
	maxColour := -1
	max := 0
	for _, v := range order {
		max = 0
		for _, u := range g.Neighbours(v) {
			if c[u] > -1 {
				seenColours[c[u]] = true
				if c[u] > max {
					max = c[u]
				}
			}
		}
		i := 0
		for i = 0; i < n; i++ {
			if seenColours[i] == false {
				c[v] = i
				if i > maxColour {
					maxColour = i
				}
				break
			}
			seenColours[i] = false
		}
		for ; i <= max; i++ {
			seenColours[i] = false
		}
	}
	return maxColour, c
}

//ChromaticNumber returns the minimum number of colours needed in a proper vertex colouring of g (known as the Chromatic Number χ) and a colouring that uses this many colours ([0, 1, ..., χ -1]).
//Note that a colouring with the minimum number of colours is not necessarily unique and the colouring returned here is arbitrary.
func ChromaticNumber(g Graph) (chromaticNumber int, colouring []int) {
	pc := make([]int, g.N())
	for i := range pc {
		pc[i] = -1
	}
	cn := CliqueNumber(g)
	return dfsDsatur(g, cn, g.N()+1, pc)
}

//IsKColorable returns true if there is a proper colouring with k colours and an example colouring, else it returns false, nil.
func IsKColorable(g Graph, k int) (ok bool, colouring []int) {
	pc := make([]int, g.N())
	for i := range pc {
		pc[i] = -1
	}
	cn, c := dfsDsatur(g, k, k, pc)
	if cn == -1 {
		return false, nil
	}
	return true, c
}

//uncolouredVertex is a type used in the dfsDsatur. It stores information about currently uncoloured vertices.
type uncolouredVertex struct {
	v                   int
	numberOfSeenColours int
	seenColours         []int
	degree              int
}

//uncolouredHeap stores the information about uncoloured vertices and a heap of which vertices are uncoloured. Information on vertices isn't modified when they are removed from the heap so it can be used again when the vertices are added back to the heap.
type uncolouredHeap struct {
	uv      []uncolouredVertex
	intHeap []int
}

func (uh uncolouredHeap) Len() int { return len(uh.intHeap) }

func (uh uncolouredHeap) Less(i, j int) bool {
	//Order the vertices in descending numberOfSeenColours and descending degree if they have the same numberOfSeenColours.
	if uh.uv[uh.intHeap[i]].numberOfSeenColours != uh.uv[uh.intHeap[j]].numberOfSeenColours {
		return uh.uv[uh.intHeap[i]].numberOfSeenColours > uh.uv[uh.intHeap[j]].numberOfSeenColours
	}
	return uh.uv[uh.intHeap[i]].degree > uh.uv[uh.intHeap[j]].degree
}

func (uh uncolouredHeap) Swap(i, j int) {
	uh.intHeap[i], uh.intHeap[j] = uh.intHeap[j], uh.intHeap[i]
}

func (uh *uncolouredHeap) Push(x interface{}) {
	v := x.(int)
	uh.intHeap = append(uh.intHeap, v)
}

func (uh *uncolouredHeap) Pop() interface{} {
	old := uh.intHeap
	n := len(old)
	v := old[n-1]
	uh.intHeap = old[0 : n-1]
	return v
}

//dfsDsatur runs a DFS search over valid colourings where the order is given by DSATUR until a colouring using at most lowerBound colours is found or all options have been exhausted.
//The search will stop when a colouring is found using lowerBound colours.
//The serach will not consider any colourings which use more than upperBound colours (but they are allowed to use upperBound colours).
//The search will start using the colouring partialColouring. -1 should be used for colours which aren't yet fixed.
//If no valid colouring can be found, the returned values are -1, nil.
func dfsDsatur(g Graph, lowerBound int, upperBound int, partialColouring []int) (chromaticNumber int, colouring []int) {
	//The code is going to assume that we already have a colouring with upperBound colours and will not look for another one.
	//We want to find a colouring with upperBound colours if it exists or we will return -1, nil.
	//We'll just increment upperBound by 1.
	upperBound++
	n := g.N()
	if n == 0 {
		return 0, []int{}
	}
	degrees := g.Degrees()

	bestColouring := make([]int, n)
	for i := range bestColouring {
		bestColouring[i] = -1
	}

	colouring = make([]int, n)
	copy(colouring, partialColouring)

	precolouredVertices := make([]int, 0)

	uncolouredV := make([]uncolouredVertex, n)
	intHeap := make([]int, 0)
	var tmp []int
	maxColourUsed := -1
	for v, c := range partialColouring {
		if c == -1 {
			tmp = make([]int, upperBound)
			uncolouredV[v] = uncolouredVertex{v, 0, tmp, degrees[v]}
			intHeap = append(intHeap, v)
		} else {
			if c > maxColourUsed {
				maxColourUsed = c
			}
			colouring[v] = c
			precolouredVertices = append(precolouredVertices, v)
		}
	}

	if maxColourUsed+1 >= upperBound {
		return -1, nil
	}

	for i := range uncolouredV {
		seenColours := uncolouredV[i].seenColours
		v := uncolouredV[i].v
		for _, u := range precolouredVertices {
			if g.IsEdge(u, v) {
				seenColours[colouring[u]]++
				if seenColours[colouring[u]] == 1 {
					uncolouredV[i].numberOfSeenColours++
				}

			}
		}
		uncolouredV[i].seenColours = seenColours
	}

	uh := uncolouredHeap{uncolouredV, intHeap}
	heap.Init(&uh)

	chosenVertices := []int{}
	currentChoice := []int{}
	choices := [][]int{}
	c := make([]int, n)
dfsLoop:
	for {
		v := -1
		c := c[:0]
		if len(uh.intHeap) > 0 {
			v = uh.intHeap[0]
			vertex := uh.uv[v]
			maxOption := upperBound - 2
			if maxColourUsed+1 < maxOption {
				maxOption = maxColourUsed + 1
			}
			for j := 0; j <= maxOption; j++ {
				b := vertex.seenColours[j]
				if b == 0 {
					c = append(c, j)
				}
			}
			if len(c) > 0 {
				heap.Remove(&uh, 0)
			}
		} else {
			copy(bestColouring, colouring)
			upperBound = maxColourUsed + 1
			if upperBound <= lowerBound {
				return upperBound, bestColouring
			}
		}
		if len(c) == 0 {
			//Backtrack
			mustChange := len(currentChoice) - 1
			for i := range chosenVertices {
				//If the colour is at least upperBound - 1, then the colouring must use at least upperBound colours and we can't improve. We must change the colour here but the only changes left for this colour are higher so we need to change something before this one.
				if colouring[chosenVertices[i]] >= upperBound-1 {
					mustChange = i - 1
					break
				}
			}

			for i := mustChange; i >= 0; i-- {
				if currentChoice[i] < len(choices[i])-1 && choices[i][currentChoice[i]+1]+1 < upperBound {
					toColour := choices[i][currentChoice[i]+1]

					//Backtrack
					for j := len(chosenVertices) - 1; j > i; j-- {
						for _, u := range uh.intHeap {
							if g.IsEdge(u, chosenVertices[j]) {
								uh.uv[u].seenColours[colouring[chosenVertices[j]]]--
								if uh.uv[u].seenColours[colouring[chosenVertices[j]]] == 0 {
									uh.uv[u].numberOfSeenColours--
								}
							}
						}
						colouring[chosenVertices[j]] = -1
						heap.Push(&uh, chosenVertices[j])
					}
					currentChoice = currentChoice[:i+1]
					choices = choices[:i+1]
					chosenVertices = chosenVertices[:i+1]

					//Change the choice at position i
					for _, u := range uh.intHeap {
						if g.IsEdge(u, chosenVertices[i]) {
							uh.uv[u].seenColours[colouring[chosenVertices[i]]]--
							if uh.uv[u].seenColours[colouring[chosenVertices[i]]] == 0 {
								uh.uv[u].numberOfSeenColours--
							}
							uh.uv[u].seenColours[toColour]++
							if uh.uv[u].seenColours[toColour] == 1 {
								uh.uv[u].numberOfSeenColours++
							}
						}
					}
					heap.Init(&uh)
					currentChoice[i]++
					colouring[chosenVertices[i]] = toColour

					maxColourUsed = 0
					for _, u := range chosenVertices {
						if colouring[u] > maxColourUsed {
							maxColourUsed = colouring[u]
						}
					}
					continue dfsLoop
				}
			}
			if bestColouring[0] == -1 {
				return -1, nil
			}
			return upperBound, bestColouring
		}
		tmp := make([]int, len(c))
		copy(tmp, c)
		choices = append(choices, tmp)
		currentChoice = append(currentChoice, 0)
		chosenVertices = append(chosenVertices, v)
		toColour := c[0]
		colouring[v] = toColour
		if toColour > maxColourUsed {
			maxColourUsed = toColour
		}

		//Update the seen colours.
		for k, u := range uh.intHeap {
			if g.IsEdge(u, v) {
				uh.uv[u].seenColours[toColour]++
				if uh.uv[u].seenColours[toColour] == 1 {
					uh.uv[u].numberOfSeenColours++
				}
			}
			heap.Fix(&uh, k)
		}

	}
}

//ChromaticIndex returns the minimum number of colours needed in a proper edge colouring of g (known as the Chromatic Index χ') and a colouring that uses this many colours ([1, ..., χ']).
//The colouring is returned in the form of an edge array with 0 for non-edges and a colour in [1, 2,..., χ'] for the edges.
//Note that a colouring with the minimum number of colours is not necessarily unique and the colouring returned here is arbitrary.
func ChromaticIndex(g Graph) (chromaticIndex int, colouredEdges []byte) {
	h := LineGraphDense(g)
	ci, colouring := ChromaticNumber(h)
	if ci == -1 {
		return -1, nil
	}
	n := 0
	colouringIndex := 0
	colouredEdges = make([]byte, n*(n-1)/2)
	index := 0
	for j := 1; j < n; j++ {
		for i := 0; i < j; i++ {
			if g.IsEdge(i, j) {
				colouredEdges[i] = byte(colouring[colouringIndex] + 1)
				colouringIndex++
			}
			index++
		}
	}
	return ci, colouredEdges
}

//ChromaticPolynomial returns the coefficients of the chromatic polynomial.
//This is a very basic implementation.
func ChromaticPolynomial(g EditableGraph) []int {
	n := g.N()
	poly := make([]int, n+1)
	type holder struct {
		g    EditableGraph
		sign int
	}
	toCheck := make([]holder, 1)
	toCheck[0] = holder{g, 1}
	var hold holder
	for len(toCheck) > 0 {
		hold, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]
		h := hold.g
		//Check if we know the chromatic polynomial of this graph.
		if h.M() == 0 {
			poly[h.N()] += hold.sign
			continue
		}

		//Choose an edge.
		var i int
		var j int
	edgeLoop:
		for i = 0; i < h.N(); i++ {
			for j = 0; j < i; j++ {
				if h.IsEdge(i, j) {
					break edgeLoop
				}
			}
		}

		//Contract and delete the edge.
		tmp := h.Copy()
		tmp.RemoveEdge(i, j)
		toCheck = append(toCheck, holder{tmp, hold.sign})

		tmp = h.Copy()
		neighbours := h.Neighbours(j)
		for _, v := range neighbours {
			tmp.AddEdge(i, v)
		}
		tmp.RemoveVertex(j)
		toCheck = append(toCheck, holder{tmp, -hold.sign})
	}
	return poly
}
