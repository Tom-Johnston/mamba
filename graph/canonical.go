package graph

import (
	"sort"

	"github.com/Tom-Johnston/mamba/disjoint"
	"github.com/Tom-Johnston/mamba/ints"
)

//OrderedPartition is a wrapper containing the information of an ordered partition and the path that the CanonicalIsomorph function used to reach the partition.
type OrderedPartition struct {
	Order      []int
	BinSizes   []int
	Path       []int
	SplitPoint int
}

//split creates a deep copy of the OrderedPartition where the ith bin has been split after the first element.
func (op OrderedPartition) split(i int) OrderedPartition {
	order := make([]int, len(op.Order))
	copy(order, op.Order)
	binSizes := make([]int, len(op.BinSizes)+1)
	copy(binSizes[:i], op.BinSizes[:i])
	copy(binSizes[i+1:], op.BinSizes[i:])
	binSizes[i] = 1
	binSizes[i+1]--
	path := make([]int, len(op.Path)+1)
	copy(path, op.Path)
	return OrderedPartition{order, binSizes, path, 0}
}

//zeroOut sets all the entries of a to be 0.
func zeroOut(a []int) {
	for i := range a {
		a[i] = 0
	}
}

//isConstant returns true if all the degreeWrapper.degrees are the same and false otherwise.
func isConstant(a degreeWrappers) bool {
	for i := 1; i < len(a); i++ {
		if len(a[0].degrees) != len(a[i].degrees) {
			return false
		}
		for j := 0; j < len(a[0].degrees); j++ {
			if a[0].degrees[j] != a[i].degrees[j] {
				return false
			}
		}
	}
	return true
}

//degreeWrapper contains a []int degrees such that the number of edges incident with the appropriate vertex of type i is given by degrees[i] and a int position used to identify the degreeWrapper in the sort.
type degreeWrapper struct {
	degrees  []int
	position int
}

//degreeWrappers is a []degreeWrapper with a sort functionality.
type degreeWrappers []degreeWrapper

func (a degreeWrappers) Len() int      { return len(a) }
func (a degreeWrappers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a degreeWrappers) Less(i, j int) bool {
	for k := 0; k < len(a[i].degrees); k++ {
		if a[i].degrees[k] != a[j].degrees[k] {
			return a[i].degrees[k] < a[j].degrees[k]
		}
	}
	return false
}

//equitableRefinementProcedure will find the largest coarsest equitable refinement of an OrderedPartition over the graph g which has edges colours in {0,1,...,maxEdgeColour}
//dws is a []degreeWrapper which is already initialised to be of the appropriate size to hold the appropriate degree data and nbs is an []int which will store the sizes of bins in the function. These are included to reduce the number of memory allocations necessary.
//This currently splits on the lexicographically least pair of bins (i,j) such that bin i is shattered by (I think) bin j.
func equitableRefinementProcedure(g Graph, op OrderedPartition, edgeColours func(i, j int) int, dws []degreeWrapper, nbs []int) OrderedPartition {
	iBinStartIndex := 0
	for i := 0; i < len(op.BinSizes); i++ {
		if op.BinSizes[i] < 2 {
			iBinStartIndex += op.BinSizes[i]
			continue
		}
		jBinStartIndex := 0
		for j := 0; j < len(op.BinSizes); j++ {
			for k := 0; k < op.BinSizes[i]; k++ {
				v := op.Order[k+iBinStartIndex]
				zeroOut(dws[k].degrees)
				dws[k].position = v
				for l := jBinStartIndex; l < jBinStartIndex+op.BinSizes[j]; l++ {
					u := op.Order[l]
					//fmt.Printf("u: %v \n", u)
					if g.IsEdge(u, v) {
						dws[k].degrees[edgeColours(u, v)-1]++
					}
				}
			}
			tmp := degreeWrappers(dws[:op.BinSizes[i]])
			//fmt.Printf("(%v, %v)", i, j)
			//fmt.Println(tmp)
			if !isConstant(tmp) {

				newBinSizes := nbs[:0]
				currentBinSize := 1
				sort.Sort(tmp)
				for k := 1; k < len(tmp); k++ {
					if ints.Equal(tmp[k-1].degrees, tmp[k].degrees) {
						currentBinSize++
					} else {
						newBinSizes = append(newBinSizes, currentBinSize)
						currentBinSize = 1
					}
				}
				newBinSizes = append(newBinSizes, currentBinSize)
				for k := 0; k < len(tmp); k++ {
					op.Order[k+iBinStartIndex] = tmp[k].position
				}
				op.BinSizes = append(op.BinSizes, newBinSizes[1:]...)
				copy(op.BinSizes[i+len(newBinSizes):], op.BinSizes[i+1:])
				copy(op.BinSizes[i:i+len(newBinSizes)], newBinSizes)
				return equitableRefinementProcedure(g, op, edgeColours, dws, nbs)
			}
			jBinStartIndex += op.BinSizes[j]
		}
		iBinStartIndex += op.BinSizes[i]
	}
	return op
}

//permuteEdges partially fills backing with the edges in the labelled subgraph of g given by the vertices in perm.
func permuteEdges(g DenseGraph, perm []int, backing []byte) []byte {
	index := 0
	for j := 1; j < len(perm); j++ {
		for i := 0; i < j; i++ {
			if perm[i] < perm[j] {
				backing[index] = g.Edges[(perm[j]*(perm[j]-1))/2+perm[i]]
			} else {
				backing[index] = g.Edges[(perm[i]*(perm[i]-1))/2+perm[j]]
			}
			index++
		}
	}
	return backing[:index]
}

type edge struct {
	from     int
	to       int
	position int
	weight   int
}

type edgeList []edge

func (a edgeList) Len() int           { return len(a) }
func (a edgeList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a edgeList) Less(i, j int) bool { return a[i].position < a[j].position }

//These must be sorted.
func compare(a, b edgeList) int {
	aIndex := 0
	bIndex := 0
	for true {
		if aIndex >= len(a) && bIndex >= len(b) {
			return 0
		}

		if aIndex >= len(a) {
			for ; bIndex < len(b); bIndex++ {
				if b[bIndex].weight > 0 {
					return -1
				} else if b[bIndex].weight < 0 {
					return 1
				}
			}
			return 0
		}
		if bIndex >= len(b) {
			for ; aIndex < len(a); aIndex++ {
				if a[aIndex].weight > 0 {
					return 1
				} else if a[aIndex].weight < 0 {
					return -1
				}
			}
			return 0
		}

		if a[aIndex].weight == 0 {
			aIndex++
			continue
		}

		if b[bIndex].weight == 0 {
			bIndex++
			continue
		}

		if a[aIndex].position < b[bIndex].position {
			if a[aIndex].weight < 0 {
				return -1
			}

			return 1
		}

		if b[bIndex].position < a[aIndex].position {
			if b[bIndex].weight < 0 {
				return 1
			}

			return -1
		}

		if a[aIndex].weight > b[bIndex].weight {
			return 1
		}
		if a[aIndex].weight < b[bIndex].weight {
			return -1
		}

		aIndex++
		bIndex++
	}
	return 0
}

func permuteEdgesList(e edgeList, perm []int, backingEdgeList edgeList, space []int) edgeList {
	//Recalculate weights
	backingEdgeList = backingEdgeList[:0]
	var position int

	if len(perm)*len(perm)/2 < len(e) {
		for j := 1; j < len(perm); j++ {
			for i := 0; i < j; i++ {
				position = 0
				if perm[i] < perm[j] {
					position = (perm[j]*(perm[j]-1))/2 + perm[i]
				} else {
					position = (perm[i]*(perm[i]-1))/2 + perm[j]
				}
				index := sort.Search(len(e), func(i int) bool { return e[i].position >= position })
				if index < len(e) && e[index].position == position {
					backingEdgeList = append(backingEdgeList, edge{from: i, to: j, position: (j*(j-1))/2 + i, weight: e[index].weight})
				}
			}

		}
		return backingEdgeList
	}

	for i := range space {
		space[i] = -1
	}

	for i, v := range perm {
		space[v] = i
	}

	for _, v := range e {
		i := v.from
		j := v.to
		if space[i] == -1 || space[j] == -1 {
			continue
		}

		indexI := space[i]
		indexJ := space[j]

		if indexI < indexJ {
			position = (indexJ*(indexJ-1))/2 + indexI
		} else {
			position = (indexI*(indexI-1))/2 + indexJ
		}
		backingEdgeList = append(backingEdgeList, edge{from: indexI, to: indexJ, position: position, weight: v.weight})
	}
	sort.Sort(backingEdgeList)
	return backingEdgeList
}

//CanonicalIsomorphWithEdgeColours returns, given a graph g with edge colours in {0,1,..,maxEdgeColour}, a []byte "canonical isomorph of g" which is constant on isomorphism classes (which repect the edge colouring).
//This returns the lexicographically most edge array formed by permuting the vertices with some additional conditions (such as increasing degree), in particular this gives a valid edge array for the representation of the isomorphism class of g.
//The canonical isomorph returned may change without warning between versions of this code so one should always update any stored canonical isomorphs whenever they change versions.
// func CanonicalIsomorphWithEdgeColours(g Graph, maxEdgeColour int) []int {
// 	n := g.N()
// 	unitPartition := make([]int, n)
// 	for i := 0; i < g.N(); i++ {
// 		unitPartition[i] = i
// 	}
// 	ci, _, _ := CanonicalIsomorphCustom(g, maxEdgeColour, OrderedPartition{unitPartition, []int{n}, []int{}, 0})
// 	return ci
// }

//CanonicalIsomorph returns a permutation which when applied to the graph gives the canonical isomorph.
//To get the actual canonical isomorph do g.InducedSubgraph(g.CanonicalIsomorph()).
func CanonicalIsomorph(g Graph) []int {
	n := g.N()
	unitPartition := make([]int, n)
	for i := 0; i < n; i++ {
		unitPartition[i] = i
	}

	f := func(i, j int) int {
		return 1
	}
	ci, _, _ := CanonicalIsomorphCustom(g, f, 1, OrderedPartition{unitPartition, []int{n}, []int{}, 0})
	return ci
}

//CanonicalIsomorphCustom is the main function with all the options which the other exposed functions wrap.
//Output is the permutation which when applied to the graph gives the canonical isomorph, the orbits of the vertices and a set of generators for the automorphism group.
func CanonicalIsomorphCustom(g Graph, edgeColours func(i, j int) int, maxEdgeColour int, initialPartition OrderedPartition) ([]int, disjoint.Set, [][]int) {
	if g.N() == 0 {
		//TODO
		return nil, nil, nil
	}

	n := g.N()

	generators := make([][]int, 0, n-1)

	currentBest := make(edgeList, g.M())

	if g.M() == 0 {
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			perm[i] = i
		}
		ds := make([]int, n)
		ds[0] = -2
		if n == 1 {
			return perm, disjoint.Set(ds), [][]int{[]int{0}}
		}
		generators := make([][]int, 2)
		tmp := make([]int, n)
		for i := range tmp {
			tmp[i] = i + 1
		}
		tmp[n-1] = 0
		generators[0] = tmp

		tmp = make([]int, n)
		for i := range tmp {
			tmp[i] = i
		}
		tmp[0] = 1
		tmp[1] = 0
		generators[1] = tmp
		return perm, disjoint.Set(ds), generators
	}

	dws := make([]degreeWrapper, n)
	for i := 0; i < n; i++ {
		dws[i].degrees = make([]int, maxEdgeColour)
	}

	nbs := make([]int, n)

	toCheck := make([]OrderedPartition, 1)
	toCheck[0] = equitableRefinementProcedure(g, initialPartition, edgeColours, dws, nbs)
	var op OrderedPartition

	var currentBestPath []int
	var currentBestPerm []int
	currentBestPermInv := make([]int, n)
	var currentBestOrbits disjoint.Set

	firstLeaf := make(edgeList, g.M())
	firstLeafPermInv := make([]int, n)
	firstLeafOrbits := disjoint.New(n)
	var firstLeafPath []int

	el := make(edgeList, 0, g.M())
	backingEL := make(edgeList, 0, g.M())
	for i := 0; i < n; i++ {
		neighbours := g.Neighbours(i)
		for _, j := range neighbours {
			if j > i {
				break
			}
			el = append(el, edge{from: i, to: j, position: (i*(i-1))/2 + j, weight: 1})
		}
	}

	space := make([]int, n)
loop:
	for len(toCheck) > 0 {

		op, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]
		if len(op.Path) > 1 && ints.HasPrefix(firstLeafPath, op.Path[:len(op.Path)-1]) {
			y := firstLeafOrbits.Find(op.Order[op.SplitPoint])
			index := op.SplitPoint + op.BinSizes[op.SplitPoint+1]
			for i := 0; i < op.BinSizes[op.SplitPoint+1]-op.Path[len(op.Path)-1]; i++ {
				if y == firstLeafOrbits.Find(op.Order[index]) {
					continue loop
				}
				index--
			}
		} else if len(op.Path) > 1 && ints.HasPrefix(currentBestPath, op.Path[:len(op.Path)-1]) {

			y := currentBestOrbits.Find(op.Order[op.SplitPoint])
			index := op.SplitPoint + op.BinSizes[op.SplitPoint+1]
			for i := 0; i < op.BinSizes[op.SplitPoint+1]-op.Path[len(op.Path)-1]; i++ {
				if y == currentBestOrbits.Find(op.Order[index]) {
					continue loop
				}
				index--
			}
		}
		op = equitableRefinementProcedure(g, op, edgeColours, dws, nbs)
		if len(op.BinSizes) == n {
			s := permuteEdgesList(el, op.Order, backingEL, space)
			if comp := compare(s, currentBest); comp == 1 {
				copy(currentBest, s)
				currentBestPath = op.Path
				currentBestPerm = op.Order
				for i := range op.Order {
					currentBestPermInv[op.Order[i]] = i
				}
				currentBestOrbits = disjoint.New(n)
				if len(firstLeafPath) == 0 {
					copy(firstLeaf, s)
					firstLeafPath = op.Path
					copy(firstLeafPermInv, currentBestPermInv)
					// fmt.Println("first", op)
					continue
				}
			} else if comp == 0 {
				//Heuristic 1
				index := 0
				for i := 0; i < n; i++ {
					if op.Path[i] != currentBestPath[i] {
						index = i
						break
					}
				}
				for len(toCheck) > 0 {
					if !ints.HasPrefix(toCheck[len(toCheck)-1].Path, op.Path[:index+1]) {
						break
					}
					toCheck = toCheck[:len(toCheck)-1]
				}
				//Update the orbits
				for i := 0; i < n; i++ {
					if tmp := op.Order[currentBestPermInv[i]]; currentBestOrbits.Find(tmp) != currentBestOrbits.Find(i) {
						currentBestOrbits.Union(i, tmp)
					}
				}

				mergesOrbits := false
				//Update the orbits
				for i := 0; i < n; i++ {
					if tmp := op.Order[currentBestPermInv[i]]; firstLeafOrbits.Find(tmp) != firstLeafOrbits.Find(i) {
						firstLeafOrbits.Union(i, tmp)
						mergesOrbits = true
					}
				}
				if mergesOrbits {
					tmp := make([]int, n)
					for i := range op.Order {
						tmp[i] = op.Order[currentBestPermInv[i]]
					}
					generators = append(generators, tmp)
				}

			}
			if comp := compare(s, firstLeaf); comp == 0 {
				//Heuristic 1
				//We will find the point in the path where they first differ. This must be somewhere as they are not the same leaf.

				index := 0
				for i := 0; i < n; i++ {
					if op.Path[i] != firstLeafPath[i] {
						index = i
						break
					}
				}
				for len(toCheck) > 0 {
					if !ints.HasPrefix(toCheck[len(toCheck)-1].Path, op.Path[:index+1]) {
						break
					}
					toCheck = toCheck[:len(toCheck)-1]
				}
				mergesOrbits := false
				//Update the orbits
				for i := 0; i < n; i++ {
					if tmp := op.Order[firstLeafPermInv[i]]; firstLeafOrbits.Find(tmp) != firstLeafOrbits.Find(i) {
						firstLeafOrbits.Union(i, tmp)
						mergesOrbits = true
					}
				}
				if mergesOrbits {
					tmp := make([]int, n)
					for i := range op.Order {
						tmp[i] = op.Order[firstLeafPermInv[i]]
					}
					generators = append(generators, tmp)
				}
			}
			continue
		}
		for i := 0; i < len(op.BinSizes); i++ {
			if op.BinSizes[i] > 1 {
				s := permuteEdgesList(el, op.Order[:i], backingEL, space)
				if len(currentBest) > 0 && compare(s, currentBest[:len(s)]) == -1 && compare(s, firstLeaf[:len(s)]) != 0 {
					break
				}
				//Split here
				for j := 0; j < op.BinSizes[i]; j++ {
					cpy := op.split(i)
					cpy.SplitPoint = i
					cpy.Path[len(cpy.Path)-1] = j
					cpy.Order[i], cpy.Order[i+j] = cpy.Order[i+j], cpy.Order[i]
					toCheck = append(toCheck, cpy)
				}
				break
			}
		}
	}

	return currentBestPerm, firstLeafOrbits, generators
}
