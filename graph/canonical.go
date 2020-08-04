package graph

import (
	"fmt"
	"math/bits"

	"github.com/Tom-Johnston/mamba/disjoint"
	"github.com/Tom-Johnston/mamba/ints"
	"github.com/Tom-Johnston/mamba/sortints"
)

//CanonicalOrderedPartition contains the current ordered partition used in the canonical isomorph function.
//A CanonicalOrderedPartition should be generated using the function NewOrderedPartition and can be reset for a second use by calling the function Reslice.
//It will be modified when used.
type CanonicalOrderedPartition struct {
	order                 []int               //The order of the vertices. The order will always be sorted within each bin.
	binDividers           []int               //The bins are [0, BinDividers[0]), [BinDividers[0], BinDividers[1]), ... This must be of capacity n.
	binAges               []int               //The time when the bin was created. This must be of capacity n and of the same size as BinDividers. Bins which are present at the start of the algorithm should have age 0.
	binsToCheck           sortints.SortedInts //The bins which have been subdivided. //TODO: Replace this with a BitSet at least for small values.
	age                   int                 //This should be initialised to 0.
	value                 []int               //This must have capacity g.M() and should be empty when initialised. TODO: Replace this with a BitSet at least for small values.
	singletonPrefixLength int                 //Starting with the first bin, how many are singletons? This is useful for keeping track of the value.
	inCell                []int               //An array giving the cell index which contains the element i.
	//TODO Is it worth having an array of singletons?
	//TODO Maybe we should move to linked lists to avoid all the copying.
}

//NewOrderedPartition creates a new ordered partition which is large enough for a graph with n vertices and m edges and initialises it with the given vertexClasses.
//If the vertexClasses is nil, all the vertices are treated as being in the same vertex class.
//A CanonicalOrderedPartition can be reset for smaller graphs but not larger graphs so searches may want to initialise a large CanonicalOrderedPartition at the start.
func NewOrderedPartition(n, m int, vertexClasses [][]int) *CanonicalOrderedPartition {
	if n == 0 {
		return nil
	}
	order := make([]int, n)
	binDividers := make([]int, n)
	inCell := make([]int, n)
	if vertexClasses == nil {
		for i := 0; i < n; i++ {
			order[i] = i
		}
		binDividers = binDividers[:1]
		binDividers[0] = n
	} else {
		binDividers = binDividers[:len(vertexClasses)]
		index := 0
		for i := range vertexClasses {
			for j := range vertexClasses[i] {
				v := vertexClasses[i][j]
				order[index] = v
				inCell[v] = j
				index++
			}
			binDividers[i] = index
		}
	}

	binAges := make([]int, len(binDividers), n)
	for i := range binAges {
		binAges[i] = 0
	}
	binsToCheck := make([]int, 1, n)
	binsToCheck[0] = 0
	value := make([]int, 0, m)
	return &CanonicalOrderedPartition{order: order, binDividers: binDividers, binAges: binAges, binsToCheck: binsToCheck, value: value, inCell: inCell}
}

//Reset resets the CanonicalOrderedPartition to the initial state for a graph with n vertices, m edges and using the given vertex classes.
//A CanonicalOrderedPartition can be reset to handle a smaller graph but will panic if reset for a larger graph.
//TODO: The case of 0.
func (op *CanonicalOrderedPartition) Reset(n, m int, vertexClasses [][]int) {
	if cap(op.order) < n {
		s := fmt.Sprintf("the partition is too small for graphs with %v vertices (cap: %v)", n, cap(op.order))
		panic(s)
	}
	if cap(op.value) < m {
		s := fmt.Sprintf("the partition is too small for graphs with %v edges (cap: %v)", m, cap(op.value))
		panic(s)
	}
	op.order = op.order[:n]
	op.inCell = op.inCell[:n]
	if vertexClasses == nil {
		for i := 0; i < n; i++ {
			op.order[i] = i
		}
		if n > 0 {
			op.binDividers = op.binDividers[:1]
			op.binDividers[0] = n
		}

		for i := range op.inCell {
			op.inCell[i] = 0
		}
	} else {
		op.binDividers = op.binDividers[:len(vertexClasses)]
		index := 0
		for i := range vertexClasses {
			for j := range vertexClasses[i] {
				v := vertexClasses[i][j]
				op.order[index] = v
				op.inCell[v] = j
				index++
			}
			op.binDividers[i] = index
		}
	}

	op.binAges = op.binAges[:len(op.binDividers)]
	for i := range op.binAges {
		op.binAges[i] = 0
	}

	if n > 0 {
		op.binsToCheck = op.binsToCheck[:1]
		op.binsToCheck[0] = 0
	}

	op.value = op.value[:0]
	op.age = 0
	op.singletonPrefixLength = 0
}

//splitBin splits the bin containing position i into a bin containing i and the rest of the partition.
//This ages the partition.
func (op *CanonicalOrderedPartition) splitBin(i int, neighbours [][]int, currentBest, firstLeaf []int) bool {

	//Age the partition.
	op.age++

	//This will be the bin containing position i
	binNumber := 0

	if op.binDividers[0] <= i {
		for j := 1; j < len(op.binDividers); j++ {
			if i < op.binDividers[j] {
				binNumber = j
				break
			}
		}
	}

	//Rearrange the values appropriately.
	binStart := 0
	if binNumber > 0 {
		binStart = op.binDividers[binNumber-1]
	}
	tmp := op.order[i]
	copy(op.order[binStart+1:], op.order[binStart:i])
	op.order[binStart] = tmp

	//Update the inCells
	for j := binStart + 1; j < len(op.order); j++ {
		op.inCell[op.order[j]]++
	}

	//Make the BinDividers longer by 1.
	op.binDividers = op.binDividers[:len(op.binDividers)+1]
	//Move the later dividers along by 1.
	copy(op.binDividers[binNumber+1:], op.binDividers[binNumber:])
	//Set the new BinDivider
	op.binDividers[binNumber] = binStart + 1

	//Do the same for the ages.
	op.binAges = op.binAges[:len(op.binAges)+1]
	copy(op.binAges[binNumber+1:], op.binAges[binNumber:])
	op.binAges[binNumber] = op.age

	//The modified bins could shatter other bins.
	op.binsToCheck.Union([]int{binNumber, binNumber + 1})

	if binNumber == op.singletonPrefixLength {
		worse := op.expandValue(neighbours, currentBest, firstLeaf)
		return worse
	}
	return false
}

//expandValue extends the current value of the CanonicalOrderedPartition. It also compares it with the currentBest and firstLeaf and returns false if there is no reason to continue exploring this branch.
//The value might be incomplete if worse is false.
func (op *CanonicalOrderedPartition) expandValue(neighbours [][]int, currentBest []int, firstLeaf []int) (worse bool) {
	for j := op.singletonPrefixLength; j < len(op.order); j++ {
		binSize := 0
		if j == 0 {
			binSize = op.binDividers[j]
		} else {
			binSize = op.binDividers[j] - op.binDividers[j-1]
		}
		if binSize != 1 {
			op.singletonPrefixLength = j
			return false
		}
		u := op.order[j]
		nbrs := neighbours[u]
		startValue := len(op.value)
		for _, v := range nbrs {
			if k := op.inCell[v]; k < j {
				op.value = append(op.value, (j*(j-1))/2+k)
			}
		}
		ints.Sort(op.value[startValue:])
		if len(currentBest) > 0 && ints.Compare(op.value, currentBest[:len(op.value)]) == -1 && ints.Compare(op.value, firstLeaf[:len(op.value)]) != 0 {
			return true
		}
	}
	op.singletonPrefixLength = len(op.order)
	return false
}

//deage rolls the age back by one, returning the partition to the state just before it reached its current age.
func (op *CanonicalOrderedPartition) deage() {
	age := op.age
	j := 0
	prev := -1
	prevDiv := 0
	for i := 0; i < len(op.binAges); i++ {
		if op.binAges[i] != age {
			//Keep this one
			op.binDividers[j] = op.binDividers[i]
			op.binAges[j] = op.binAges[i]

			//Check if we have merged with other bins
			if i-prev > 1 {
				//We have merged.
				//This bin is no longer a singleton.
				if j < op.singletonPrefixLength {
					op.singletonPrefixLength = j
					maxPos := ((j - 1) * j) / 2
					k := len(op.value) - 1
					for ; k >= 0; k-- {
						if op.value[k] < maxPos {
							break
						}
					}
					k++
					op.value = op.value[:k]
				}
				//Order the elements in the bin
				ints.Sort(op.order[prevDiv:op.binDividers[j]])
			}
			prev = i
			prevDiv = op.binDividers[j]
			j++
		}
	}
	op.binDividers = op.binDividers[:j]
	op.binAges = op.binAges[:j]
	op.binsToCheck = op.binsToCheck[:0]

	//Update the inCell array. Not the best implementation here.
	currBin := 0
	for i := range op.order {
		if op.binDividers[currBin] == i {
			currBin++
		}
		op.inCell[op.order[i]] = currBin
	}

	op.age--
}

//keyValue holds a key value pair. There are functions for sorting a []keyValues at the end of this file.
type keyValue struct {
	value int
	key   int
}

//equitableRefinementProcedure will find the largest coarsest equitable refinement of an CanonicalOrderedPartition over the the graph with the given neighbourhoods.
//The refinement will return false if there is no point investigating this part of the tree further due to the value we will get. If it returns false, it may have terminated early.
//If the option CheckViability is specified, the procedure will also check if the vertex n - 1 is still in the same bin as the first vertex in ViableBits. This is useful when doing a canonical deletion algorithm.
//The majority of the paremeters here are storage space to avoid allocations and speed up the code.
//TODO: Can we reuse some of the workspace and remove some parameters?
func equitableRefinementProcedure(neighbours [][]int, op *CanonicalOrderedPartition, dws []keyValue, nbs, space, timesSeen, maxCell, numberOfMax []int, currentBest []int, firstLeaf []int, options *CanonicalOptions) (worse bool) {
	//Check if the bin i shatters a bin j.
	n := len(op.order)

	for len(op.binsToCheck) > 0 {
		zeroOut(timesSeen)
		zeroOut(maxCell)
		zeroOut(numberOfMax)
		maxCell = maxCell[:len(op.binDividers)]
		numberOfMax = numberOfMax[:len(op.binDividers)]
		i := op.binsToCheck[len(op.binsToCheck)-1]
		op.binsToCheck = op.binsToCheck[:len(op.binsToCheck)-1]
		iBinStart := 0
		if i > 0 {
			iBinStart = op.binDividers[i-1]
		}

		for wIndex := iBinStart; wIndex < op.binDividers[i]; wIndex++ {
			w := op.order[wIndex]
			for _, v := range neighbours[w] {
				timesSeen[v]++
				cell := op.inCell[v]
				if timesSeen[v] > maxCell[cell] {
					numberOfMax[cell] = 1
					maxCell[cell] = timesSeen[v]
				} else if timesSeen[v] == maxCell[cell] {
					numberOfMax[cell]++
				}
			}
		}

		//Check if bin i shatters bin j.
		//TODO We could store which jBins we actually need to check. But we need a better set for this.
		for j := len(op.binDividers) - 1; j >= 0; j-- {
			binStart := 0
			if j > 0 {
				binStart = op.binDividers[j-1]
			}
			binSize := op.binDividers[j] - binStart
			if binSize == 1 || maxCell[j] == 0 || numberOfMax[j] == binSize {
				continue
			}

			//We should also handle the case maxCell[j] == 2
			if maxCell[j] == 1 {
				dws = dws[:binSize]
				zeroIndex := 0
				oneIndex := binSize - numberOfMax[j]
				for k := 0; k < binSize; k++ {
					v := op.order[binStart+k]
					if timesSeen[v] == 0 {
						dws[zeroIndex] = keyValue{value: 0, key: v}
						zeroIndex++
					} else {
						dws[oneIndex] = keyValue{value: 1, key: v}
						oneIndex++
					}
				}
			} else {
				dws = dws[:binSize]
				for k := 0; k < binSize; k++ {
					v := op.order[binStart+k]
					dws[k].value = timesSeen[v]
					dws[k].key = v
				}
				//Sort the new elements.
				stable(dws, len(dws))
			}

			//Create the new bins and update the order.
			//Also mark the new bins as needing to be checked.
			nbsIndex := 0
			nbs = nbs[:n]
			op.order[binStart] = dws[0].key
			for k := 1; k < binSize; k++ {
				op.order[binStart+k] = dws[k].key
				if dws[k].value != dws[k-1].value {
					nbs[nbsIndex] = binStart + k
					nbsIndex++
				}
			}

			nbs = nbs[:nbsIndex]

			for k := len(op.binsToCheck) - 1; k >= 0; k-- {
				if op.binsToCheck[k] <= j {
					break
				}
				op.binsToCheck[k] += nbsIndex
			}

			op.binDividers = op.binDividers[:len(op.binDividers)+len(nbs)]
			copy(op.binDividers[j+len(nbs):], op.binDividers[j:])
			copy(op.binDividers[j:], nbs)

			op.binAges = op.binAges[:len(op.binAges)+len(nbs)]
			copy(op.binAges[j+len(nbs):], op.binAges[j:])

			for i := 0; i < len(nbs); i++ {
				op.binAges[j+i] = op.age
			}

			space = space[:nbsIndex+1]
			for k := j; k < j+nbsIndex+1; k++ {
				space[k-j] = k
			}

			op.binsToCheck.Union(space)

			//Update the inCell array. Not the best implementation here.
			currBin := 0
			for i := range op.order {
				if op.binDividers[currBin] == i {
					currBin++
				}
				op.inCell[op.order[i]] = currBin
			}

			if j == op.singletonPrefixLength {
				worse = op.expandValue(neighbours, currentBest, firstLeaf)
				if worse {
					return true
				}
			}

			if options.CheckViability {
				//Check where the last element is.
				cell := op.inCell[len(op.order)-1]
				viable := options.ViableBits
				x := uint(viable)
				for x != 0 {
					//y isolates the least significant bit in x
					y := x & -x
					//The index is the number of trailing zeros.
					v := bits.TrailingZeros(x)
					//Unset the least significant bit
					x ^= y
					if op.inCell[v] < cell {
						return true
					}
					if op.inCell[v] > cell {
						viable ^= y
					}
				}
			}
		}
	}
	return false
}

//CanonicalStorage contains the working space for the CanonicalIsomorphAllocated function and can be reused multiple times. It will be suitably reset by the call to CanonicalIsomorphAllocated.
//To create the storage space for a graph on at most n vertices and with at most m edges, call NewStorage(n, m)
type CanonicalStorage struct {
	path       []int   //Cap >= n advised
	choices    []int   //Cap >= n advised
	generators [][]int //Cap >= n-1 and ideally each []int of size n. For one off calls it might be wasteful to completely initialise this so the function will allocate as necessary.

	currentBest        []int        //Cap >= m
	currentBestPath    []int        //Cap >= n
	currentBestPerm    []int        //Cap >= n
	currentBestPermInv []int        //Cap >= n
	currentBestOrbits  disjoint.Set //Cap >= n. Doesn't need to be reset.

	firstLeaf        []int        //Cap >= m
	firstLeafPermInv []int        //Cap >= n
	firstLeafOrbits  disjoint.Set //Cap >= n. Doesn't need to be reset
	firstLeafPath    []int        //Cap >=n

	space       []int      //Cap >= n
	dws         []keyValue //Cap >= n
	nbs         []int      //Cap >= n
	timesSeen   []int      //Cap >= n
	maxCell     []int      //Cap >= n
	numberOfMax []int      //Cap >= n
}

//NewStorage creates the necessary storage space for a call to CanonicalIsomorphAllocated with a graph with at most n vertices and m edges.
func NewStorage(n, m int) *CanonicalStorage {
	cs := new(CanonicalStorage)
	cs.path = make([]int, 0, n)
	cs.choices = make([]int, 0, n)
	if n > 0 {
		cs.generators = make([][]int, n-1)
	}

	cs.currentBest = make([]int, 0, m)
	cs.currentBestPath = make([]int, n)
	cs.currentBestPerm = make([]int, n)
	cs.currentBestPermInv = make([]int, n)
	cs.currentBestOrbits = disjoint.New(n)

	cs.firstLeaf = make([]int, 0, m)
	cs.firstLeafPermInv = make([]int, n)
	cs.firstLeafOrbits = disjoint.New(n)
	cs.firstLeafPath = make([]int, n)

	cs.space = make([]int, n)
	cs.dws = make([]keyValue, n)
	cs.nbs = make([]int, n)
	cs.timesSeen = make([]int, n)
	cs.maxCell = make([]int, n)
	cs.numberOfMax = make([]int, n)

	return cs
}

//CanonicalOptions is a struct containing the options for a call to CanonicalIsomorphAllocated.
//The zero value gives the default settings.
type CanonicalOptions struct {
	CheckViability bool //If CheckViability is true, the first call to equitableRefinementPartition will check that the last vertex (n-1) is in the same bin as the earliest vertex in ViableBits. If the vertex is in a different bin, CanonicalIsomorphAllocated returns nil, nil, nil (and this is the only time it does so). No guess has been made at this point so if two vertices are in different bins, they are in different orbits. This is useful for a canonical deletion search. Note that since ViableBits is a uint, this setting is valid for small-ish graphs.
	ViableBits     uint //Put a one in bit i (starting from the low bits) if you want the function to return early if the vertex i is in an earlier bin than the vertex n - 1.
}

//CanonicalIsomorph returns a permutation which gives the canonical isomorph when applied to the graph. The actual canonical isomorph can be obtained by running InducedSubgraph(CanonicalIsomorph(g)).
//The canonical isomorph is the chosen representation of the isomorphism class containing g. In particular, the canonical isomorph of two graphs is the same if and only if they are isomorphic.
func CanonicalIsomorph(g Graph) []int {
	ci, _, _ := CanonicalIsomorphFull(g, nil)
	return ci
}

//CanonicalIsomorphFull returns the permutation which when applied to g gives the canonical isomorph, a disjoint.Set giving the vertex orbits and a set of generators for the autmorphism group of g.
func CanonicalIsomorphFull(g Graph, vertexClasses [][]int) ([]int, disjoint.Set, [][]int) {
	op := NewOrderedPartition(g.N(), g.M(), vertexClasses)
	neighbours := make([][]int, g.N())
	for i := range neighbours {
		neighbours[i] = g.Neighbours(i)
	}
	return CanonicalIsomorphAllocated(g.N(), g.M(), neighbours, op, NewStorage(g.N(), g.M()), new(CanonicalOptions))
}

//CanonicalIsomorphAllocated returns the same values as CanonicalIsomorphFull but requires more setup. This means that the CanonicalStorage and CanonicalOrderedPartition can be resued to reduce the number of allocations when calling this for many graphs. This function is also currently the only function which allows the setting of options.
//It is not recommended to call this function unless you know you need to reduce the allocations or you need to set options. See the source for CanonicalIsomorphFull for an example of how to call this function.
//Note that op, storage and options may be modified when calling this function and modifying storage may modify the output of this function.
func CanonicalIsomorphAllocated(n, m int, neighbours [][]int, op *CanonicalOrderedPartition, storage *CanonicalStorage, options *CanonicalOptions) ([]int, disjoint.Set, [][]int) {
	if n == 0 {
		return []int{}, nil, nil
	}

	count := 0

	generators := storage.generators[:0]
	currentBest := storage.currentBest[:0]

	//Handle the special case where m = 0.
	//TODO: Check if this is necessary.
	if m == 0 {
		//Return the identity permutation.
		perm := storage.currentBestPerm[:n]
		for i := 0; i < n; i++ {
			perm[i] = i
		}
		//Every vertex is in the same orbit.
		ds := storage.firstLeafOrbits[:n]
		ds[0] = -2
		for i := 1; i < n; i++ {
			ds[i] = 0
		}

		if n == 1 {
			return perm, ds, storage.generators[:0]
		}

		generators := storage.generators[:1]
		tmp := generators[0]
		if cap(tmp) < n {
			tmp = make([]int, n)
		} else {
			tmp = tmp[:n]
		}
		for i := range tmp {
			tmp[i] = i + 1
		}
		tmp[n-1] = 0
		generators[0] = tmp

		if n == 2 {
			return perm, ds, generators
		}

		generators = generators[:2]

		tmp = generators[1]
		if cap(tmp) < n {
			tmp = make([]int, n)
		} else {
			tmp = tmp[:n]
		}
		for i := range tmp {
			tmp[i] = i
		}
		tmp[0] = 1
		tmp[1] = 0
		generators[1] = tmp
		return perm, ds, generators
	}

	path := storage.path[:0]
	choices := storage.choices[:0]

	currentBestPath := storage.currentBestPath[:n]
	currentBestPerm := storage.currentBestPerm[:n]
	currentBestPermInv := storage.currentBestPermInv[:n]
	currentBestOrbits := storage.currentBestOrbits[:n]

	firstLeaf := storage.firstLeaf[:m]
	firstLeafPermInv := storage.firstLeafPermInv[:n]
	firstLeafOrbits := storage.firstLeafOrbits[:n]
	firstLeafPath := storage.firstLeafPath[:n]

	space := storage.space[:n]
	dws := storage.dws[:n]
	nbs := storage.nbs[:n]
	timesSeen := storage.timesSeen[:n]
	maxCell := storage.maxCell[:n]
	numberOfMax := storage.numberOfMax[:n]

	skipDeage := false

	//Split the partition.
	//We split here and at the end of the loop so we can easily handle the CheckViable option. It wouldn't be hard to check it the other way but might require a
	worse := equitableRefinementProcedure(neighbours, op, dws, nbs, space, timesSeen, maxCell, numberOfMax, currentBest, firstLeaf, options)
	if options.CheckViability {
		//Disable the check for any further iterations
		options.CheckViability = false
		if worse {
			return nil, nil, nil
		}
	}
	for {
		if !worse && len(op.binDividers) == n {
			count++
			//Are we the new best?
			if comp := ints.Compare(op.value, currentBest); comp == 1 {
				currentBest = currentBest[:m]
				copy(currentBest, op.value)
				copy(currentBestPath, path)
				copy(currentBestPerm, op.order)
				for i := range op.order {
					currentBestPermInv[op.order[i]] = i
					currentBestOrbits[i] = -1
				}
				if count == 1 {
					copy(firstLeaf, op.value)
					copy(firstLeafPath, path)
					copy(firstLeafPermInv, currentBestPermInv)
					copy(firstLeafOrbits, currentBestOrbits)
				}
			} else if comp == 0 {
				//We are the same. This means we can update the automorphism group based on the currentBest.

				//Update the orbits
				for i := 0; i < n; i++ {
					if tmp := op.order[currentBestPermInv[i]]; currentBestOrbits.FindBuffered(tmp, space) != currentBestOrbits.FindBuffered(i, space) {
						currentBestOrbits.UnionBuffered(i, tmp, space)
					}
				}

				mergesOrbits := false
				//Update the orbits
				for i := 0; i < n; i++ {
					if tmp := op.order[currentBestPermInv[i]]; firstLeafOrbits.FindBuffered(tmp, space) != firstLeafOrbits.FindBuffered(i, space) {
						firstLeafOrbits.UnionBuffered(i, tmp, space)
						mergesOrbits = true
					}
				}
				if mergesOrbits {
					generators = generators[:len(generators)+1]
					tmp := generators[len(generators)-1]
					if cap(tmp) >= n {
						tmp = tmp[:n]
					} else {
						tmp = make([]int, n)
					}
					for i := range op.order {
						tmp[i] = op.order[currentBestPermInv[i]]
					}
					generators[len(generators)-1] = tmp
				}

				//TODO Maybe we want to store this instead of calculating it every time?
				//Heuristic 1
				index := len(path) - 1
				for i := 0; i < len(path)-1; i++ {
					if path[i] != currentBestPath[i] {
						index = i
						break
					}
				}

				// toTrim := 0
				for i := len(path) - 1; i > index; i-- {
					op.deage()
					// toTrim += path[i]
				}
				path = path[:index+1]
				choices = choices[:index+1]
			} else if comp2 := ints.Compare(op.value, firstLeaf); comp2 == 0 {
				//We will find the point in the path where they first differ. This must be somewhere as they are not the same leaf.

				mergesOrbits := false
				//Update the orbits
				for i := 0; i < n; i++ {
					if tmp := op.order[firstLeafPermInv[i]]; firstLeafOrbits.FindBuffered(tmp, space) != firstLeafOrbits.FindBuffered(i, space) {
						firstLeafOrbits.UnionBuffered(i, tmp, space)
						mergesOrbits = true
					}
				}
				if mergesOrbits {
					generators = generators[:len(generators)+1]
					tmp := generators[len(generators)-1]
					if cap(tmp) >= n {
						tmp = tmp[:n]
					} else {
						tmp = make([]int, n)
					}
					for i := range op.order {
						tmp[i] = op.order[firstLeafPermInv[i]]
					}
					generators[len(generators)-1] = tmp
				}
				//Heuristic 1
				index := len(path) - 1
				for i := 0; i < len(path)-1; i++ {
					if path[i] != firstLeafPath[i] {
						index = i
						break
					}
				}
				// toTrim := 0
				for i := len(path) - 1; i > index; i-- {
					op.deage()
					// toTrim += path[i]
				}
				path = path[:index+1]
				choices = choices[:index+1]
			}
		} else if !worse {
			//We are not a leaf so we need to split.
			prevBinStart := 0
			for i := 0; i < len(op.binDividers); i++ {
				binSize := op.binDividers[i] - prevBinStart
				if binSize > 1 {
					//We will split on this bin.
					// for j := 0; j < binSize; j++ {
					// 	choices = append(choices, prevBinStart+j)
					// }
					choices = append(choices, op.binDividers[i])
					path = append(path, binSize)
					skipDeage = true
					break
				}
				prevBinStart = op.binDividers[i]
			}
		}

		//Try the next node in the dfs.
	stepLoop:
		for {
			// fmt.Println(path, choices)
			if len(path) == 0 {
				return currentBestPerm, firstLeafOrbits, generators
			}
			//Potential to step to.

		jLoop:
			for j := path[len(path)-1] - 1; j >= 0; j-- {
				// fmt.Println(path, choices)
				if !skipDeage {
					op.deage()
				} else {
					skipDeage = false
				}

				choices[len(choices)-1]--
				choicePosition := choices[len(choices)-1]
				choiceElement := op.order[choicePosition]
				//Can we step there?
				//Heuristic 2
				if count > 0 && ints.HasPrefix(firstLeafPath, path[:len(path)-1]) {
					if firstLeafOrbits[choiceElement] >= 0 {
						skipDeage = true
						continue jLoop
					}

				}

				//Do the same for the currentBest
				//Heuristic 2
				if count > 0 && ints.HasPrefix(currentBestPath, path[:len(path)-1]) {
					if currentBestOrbits[choiceElement] >= 0 {
						skipDeage = true
						continue jLoop
					}
				}

				//Success
				worse := op.splitBin(choicePosition, neighbours, currentBest, firstLeaf)
				path[len(path)-1] = j
				//Is the incremental leaf certificate worse?
				if worse {
					//Don't accept this step.
					continue jLoop
				}

				break stepLoop
			}

			//We weren't successful.
			//Take a step back.
			if !skipDeage {
				op.deage()
			} else {
				skipDeage = false
			}
			path = path[:len(path)-1]
			choices = choices[:len(choices)-1]
		}
		//End of stepping
		worse = equitableRefinementProcedure(neighbours, op, dws, nbs, space, timesSeen, maxCell, numberOfMax, currentBest, firstLeaf, options)
	}
}

//Below are various helper functions.

//zeroOut sets all the entries of a to be 0.
//Note that this will be optimised to a memclr call.
func zeroOut(a []int) {
	for i := range a {
		a[i] = 0
	}
}

//The following code is copied directly from sort.Sort and modified to work with []keyValue. This saves on both time and allocations compared to using sort.Sort.

// Copyright 2009 The Go Authors. All rights reserved.

//stable sorts the data by value using a stable sort.
func stable(data []keyValue, n int) {
	blockSize := 20 // must be > 0
	a, b := 0, blockSize
	for b <= n {
		insertionSortKeyValue(data, a, b)
		a = b
		b += blockSize
	}
	insertionSortKeyValue(data, a, n)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMerge(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		if m := a + blockSize; m < n {
			symMerge(data, a, m, n)
		}
		blockSize *= 2
	}
}

func insertionSortKeyValue(data []keyValue, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data[j].value < data[j-1].value; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

// SymMerge merges the two sorted subsequences data[a:m] and data[m:b] using
// the SymMerge algorithm from Pok-Son Kim and Arne Kutzner, "Stable Minimum
// Storage Merging by Symmetric Comparisons", in Susanne Albers and Tomasz
// Radzik, editors, Algorithms - ESA 2004, volume 3221 of Lecture Notes in
// Computer Science, pages 714-723. Springer, 2004.
//
// Let M = m-a and N = b-n. Wolog M < N.
// The recursion depth is bound by ceil(log(N+M)).
// The algorithm needs O(M*log(N/M + 1)) calls to data.Less.
// The algorithm needs O((M+N)*log(M)) calls to data.Swap.
//
// The paper gives O((M+N)*log(M)) as the number of assignments assuming a
// rotation algorithm which uses O(M+N+gcd(M+N)) assignments. The argumentation
// in the paper carries through for Swap operations, especially as the block
// swapping rotate uses only O(M+N) Swaps.
//
// symMerge assumes non-degenerate arguments: a < m && m < b.
// Having the caller check this condition eliminates many leaf recursion calls,
// which improves performance.
func symMerge(data []keyValue, a, m, b int) {
	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[a] into data[m:b]
	// if data[a:m] only contains one element.
	if m-a == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] >= data[a] for m <= i < b.
		// Exit the search loop with i == b in case no such index exists.
		i := m
		j := b
		for i < j {
			h := int(uint(i+j) >> 1)
			if data[h].value < data[a].value {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[a] reaches the position before i.
		for k := a; k < i-1; k++ {
			data[k], data[k+1] = data[k+1], data[k]
		}
		return
	}

	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[m] into data[a:m]
	// if data[m:b] only contains one element.
	if b-m == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] > data[m] for a <= i < m.
		// Exit the search loop with i == m in case no such index exists.
		i := a
		j := m
		for i < j {
			h := int(uint(i+j) >> 1)
			if data[m].value >= data[h].value {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[m] reaches the position i.
		for k := m; k > i; k-- {
			data[k], data[k-1] = data[k-1], data[k]
		}
		return
	}

	mid := int(uint(a+b) >> 1)
	n := mid + m
	var start, r int
	if m > mid {
		start = n - b
		r = mid
	} else {
		start = a
		r = m
	}
	p := n - 1

	for start < r {
		c := int(uint(start+r) >> 1)
		if data[p-c].value >= data[c].value {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < m && m < end {
		rotate(data, start, m, end)
	}
	if a < start && start < mid {
		symMerge(data, a, start, mid)
	}
	if mid < end && end < b {
		symMerge(data, mid, end, b)
	}
}

func swapRange(data []keyValue, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a+i], data[b+i] = data[b+i], data[a+i]
	}
}

// Rotate two consecutive blocks u = data[a:m] and v = data[m:b] in data:
// Data of the form 'x u v y' is changed to 'x v u y'.
// Rotate performs at most b-a many calls to data.Swap.
// Rotate assumes non-degenerate arguments: a < m && m < b.
func rotate(data []keyValue, a, m, b int) {
	i := m - a
	j := b - m

	for i != j {
		if i > j {
			swapRange(data, m-i, m, j)
			i -= j
		} else {
			swapRange(data, m-i, m+j-i, i)
			j -= i
		}
	}
	// i == j
	swapRange(data, m-i, m, i)
}
