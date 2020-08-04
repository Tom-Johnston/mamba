package search

import (
	"math/bits"

	"github.com/Tom-Johnston/mamba/comb"
	"github.com/Tom-Johnston/mamba/disjoint"
	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/ints"
	"github.com/Tom-Johnston/mamba/itertools"
)

//AllParallel sends all non-isomorphic graphs on n vertices to the output channel which it then closes. It automatically splits the work across numWorkers goroutines.
func AllParallel(n int, output chan *graph.DenseGraph, numWorkers int) {
	counter := numWorkers
	count := make(chan bool)
	for i := 0; i < numWorkers; i++ {
		c := make(chan *graph.DenseGraph)
		tmp := i
		go func() {
			go All(n, c, tmp, numWorkers)
			for v := range c {
				output <- v
			}
			count <- false
		}()
	}
	for {
		<-count
		counter--
		if counter == 0 {
			close(output)
		}
	}
}

//searchGraph holds a graph and the information from the canonical isomorph.
//Note that Perm, Generators and Orbits are not necessarily deep copies of the relevant part of storage.
type searchGraph struct {
	G          *graph.DenseGraph
	Neighbours [][]int //We don't keep this up to date yet.
	Perm       []int
	Generators [][]int
	Orbits     disjoint.Set
}

//All with a = 0 and m = 1 sends all non-isomorphic graphs on n vertices to the output channel which it then closes. In general, this will begin a search using canonical deletion to find all graphs on at most n vertices where the choice at level ceil(2n/3) is equal to a mod m. For small values of m this should produce a fairly even split and allow for some small parallelism.
func All(n int, output chan *graph.DenseGraph, a int, m int) {
	WithPruning(n, output, a, m, func(g *graph.DenseGraph) bool { return false }, func(g *graph.DenseGraph) bool { return false })
}

//WithPruningParallel uses numWorkers goroutines to output all non-isomorphic graphs on n vertices such that neither the graph itself nor any of its predecessors were pruned.
func WithPruningParallel(n int, output chan *graph.DenseGraph, numWorkers int, preprune, prune func(g *graph.DenseGraph) bool) {
	counter := numWorkers
	count := make(chan bool)
	for i := 0; i < numWorkers; i++ {
		c := make(chan *graph.DenseGraph)
		tmp := i
		go func() {
			go WithPruning(n, c, tmp, numWorkers, preprune, prune)
			for v := range c {
				output <- v
			}
			count <- false
		}()
	}
	for {
		<-count
		counter--
		if counter == 0 {
			close(output)
		}
	}
}

//WithPruning with a = 0 and m = 1 outputs all non-isomorphic graphs on n vertices such that neither the graph itself nor any of its predecessors were pruned.
//In general, only the choices which are equal to a mod m are output to allow a small amount of parallelism.
//The function prune is called once an augmentation has been determined to be canonical whereas preprune is called once the augmentation has been applied to the graph but before checking if the augmentation is canonical. Preprune is not applied to graphs on 0 or 1 vertices and prune is not checked for graphs on 0 vertices unless n is 0.
//Note that the vertices are added in the order 0, 1, ..., n-1 and, when adding the kth vertex, the graph induced on {0, 1 \dots, k-2} has already not been pruned.
func WithPruning(n int, output chan *graph.DenseGraph, a int, m int, preprune, prune func(g *graph.DenseGraph) bool) {
	//Handle a couple of small cases so that we can start with a graph with 1 vertex and expand from there.
	//This saves a check later on so we might as well do it here.
	if n == 0 {
		g := graph.NewDense(0, nil)
		if a == 0 && !prune(g) {
			output <- g
		}
		close(output)
		return
	}

	if n == 1 {
		g := graph.NewDense(1, nil)
		if a == 0 && !prune(g) {
			output <- g
		}
		close(output)
		return
	}

	//Intitialise a graph large enough to hold a graph on n vertices and then reset it to the graph on 1 vertex.
	g := graph.NewDense(n, nil)
	g.NumberOfVertices = 1
	g.DegreeSequence = g.DegreeSequence[:1]
	g.Edges = g.Edges[:0]

	//Check if the graph on 1 vertex should be pruned.
	if prune(g) {
		close(output)
		return
	}

	//Initialise storage for the neighbours.
	neighbours := make([][]int, n)
	for i := range neighbours {
		neighbours[i] = make([]int, 0, n)
	}

	sg := &searchGraph{G: g, Neighbours: neighbours, Generators: nil, Orbits: nil}

	//Initialise the storage etc. for CanonicalIsomorphCustom
	storage := graph.NewStorage(n, (n*(n-1))/2)
	op := graph.NewOrderedPartition(n, (n*(n-1))/2, nil)
	options := new(graph.CanonicalOptions)

	//Initialise storage which will be used when checking what sets of neighbours we should try to add.
	ds := make(disjoint.Set, comb.Coeff(n, n/2))

	stepForward := false
	choices := make([]uint, 0)
	currentPath := make([]int, 0, n)
	v := make([]int, 0, n)

	splitLevel := 2 * (n + 1) / 3
	splitLevel--
	for true {
		if sg.G.NumberOfVertices == n {
			output <- (sg.G.Copy()).(*graph.DenseGraph)
		} else {
			//Prepare to go deeper.
			numAugs := addAugmentations(sg, &choices, ds, op, storage, options)
			currentPath = append(currentPath, numAugs)
			stepForward = true
		}

		//Step loop
	stepLoop:
		for true {
			if len(choices) == 0 {
				close(output)
				return
			}
			for i := currentPath[len(currentPath)-1] - 1; i >= 0; i-- {
				//Splitting
				level := len(currentPath)
				x := choices[len(choices)-1]
				choices = choices[:len(choices)-1]

				if i%m != a && level == splitLevel {
					continue
				}

				v = v[:0]
				for x != 0 {
					//y isolates the least significant bit in x
					y := x & -x
					//The index is the number of trailing zeros.
					v = append(v, bits.TrailingZeros(x))
					//Unset the least significant bit
					x ^= y
				}

				//Are we moving deeper or to the next child?
				if !stepForward {
					sg.G.RemoveVertex(sg.G.NumberOfVertices - 1)
					clearAutomorphismGroup(sg)
				}

				stepForward = false

				//Take the step
				sg.G.AddVertex(v)
				clearAutomorphismGroup(sg)

				//Check if this should be prepruned.
				if preprune(sg.G) {
					continue
				}

				//Are we canonical and should we be pruned?
				if isCanonical(sg, v, op, storage, options) && !prune(sg.G) {
					//This step is valid so we are done.
					currentPath[len(currentPath)-1] = i
					break stepLoop
				}
				//fmt.Println("Not canonical")
			}
			//None of the options on this level worked so take a step back
			if !stepForward {
				sg.G.RemoveVertex(sg.G.NumberOfVertices - 1)
				clearAutomorphismGroup(sg)
			}
			stepForward = false
			currentPath = currentPath[:len(currentPath)-1]
			// choices = choices[:len(choices)-binSize]
		}
	}
}

func updateNeighbours(sg *searchGraph) {
	sg.Neighbours = sg.Neighbours[:sg.G.NumberOfVertices]
	for v := 0; v < sg.G.NumberOfVertices; v++ {
		r := sg.Neighbours[v][:0]
		tmp := (v * (v - 1)) / 2
		for i := 0; i < v; i++ {
			index := tmp + i
			if sg.G.Edges[index] > 0 {
				r = append(r, i)
			}
		}

		for i := v + 1; i < sg.G.NumberOfVertices; i++ {
			index := (i*(i-1))/2 + v
			if sg.G.Edges[index] > 0 {
				r = append(r, i)
			}
		}
		sg.Neighbours[v] = r
	}
}

func getAutomorphismGroup(sg *searchGraph, op *graph.CanonicalOrderedPartition, storage *graph.CanonicalStorage, options *graph.CanonicalOptions) {
	op.Reset(sg.G.NumberOfVertices, sg.G.NumberOfEdges, nil)
	updateNeighbours(sg)
	perm, orbits, generators := graph.CanonicalIsomorphAllocated(sg.G.NumberOfVertices, sg.G.NumberOfEdges, sg.Neighbours, op, storage, options)
	sg.Perm = perm
	sg.Generators = generators
	sg.Orbits = orbits
}

func clearAutomorphismGroup(sg *searchGraph) {
	sg.Perm = nil
	sg.Generators = nil
	sg.Orbits = nil
}

func addAugmentations(sg *searchGraph, choices *[]uint, ds disjoint.Set, op *graph.CanonicalOrderedPartition, storage *graph.CanonicalStorage, options *graph.CanonicalOptions) int {
	n := sg.G.NumberOfVertices
	minDegree := ints.Min(sg.G.DegreeSequence)
	maxSize := minDegree + 1

	numFound := 0

	//TODO Reuse some preallocated space in augs?

	//TODO Do we want to be blocking obviously non-canonical sets?
	//No point checking a set of of minDegree + 1 if we aren't adjacent to a vertex of min degree for example.

	//Handle k == 0 separately
	numFound++
	*choices = append(*choices, 0)

	//For all the other large sets we will want the generators
	if sg.Perm == nil {
		options.CheckViability = false
		getAutomorphismGroup(sg, op, storage, options)
	}

	//Handle k == 1 separately
	for i, v := range sg.Orbits {
		if v < 0 {
			*choices = append(*choices, (1 << uint(i)))
			numFound++
		}
	}

	//We are just going to be dumb here.
	//This really needs improving.
	//TODO IMPORTANT This is basically the slowest bit
	for k := 2; k <= maxSize; k++ {
		ds = ds[:comb.Coeff(n, k)]
		for i := range ds {
			ds[i] = -1
		}
		//TODO Should these be passed in.
		buf := make([]int, n)
		iter := itertools.CombinationsColex(n, k)
		c2 := make([]int, k)
		for i := 0; i < len(ds); i++ {
			iter.Next()
			c := iter.Value()
			for _, g := range sg.Generators {
				for j := 0; j < k; j++ {
					c2[j] = g[c[j]]
				}
				ints.Sort(c2)
				ds.UnionBuffered(i, comb.Rank(c2), buf)
			}
		}
		iter = itertools.CombinationsColex(n, k)
		for i := 0; i < len(ds); i++ {
			iter.Next()
			if ds[i] < 0 {
				x := 0
				for _, v := range iter.Value() {
					x |= (1 << uint(v))
				}
				*choices = append(*choices, uint(x))
				numFound++
			}
		}
	}
	return numFound
}

func isCanonical(sg *searchGraph, aug []int, op *graph.CanonicalOrderedPartition, storage *graph.CanonicalStorage, options *graph.CanonicalOptions) bool {
	n := sg.G.NumberOfVertices

	viableBits := uint(0)
	degrees := sg.G.DegreeSequence

	//Check the degree
	degree := degrees[n-1]
	for i := 0; i < n-1; i++ {
		if degrees[i] < degree {
			return false
		} else if degrees[i] == degree {
			viableBits |= (1 << uint(i))
		}
	}

	if viableBits == 0 {
		return true
	}

	//Sums and squares of degrees
	sum := 0
	square := 0
	for _, v := range aug {
		sum += degrees[v]
		square += degrees[v] * degrees[v]
	}
	sumV := 0
	squareV := 0
	x := uint(viableBits)
	for x != 0 {
		//y isolates the least significant bit in x
		y := x & -x
		//The index is the number of trailing zeros.
		v := bits.TrailingZeros(x)
		//Unset the least significant bit
		x ^= y
		sumV = 0
		squareV = 0
		for j := 0; j < n; j++ {
			if v > j {
				if sg.G.Edges[(v*(v-1))/2+j] == 1 {
					sumV += degrees[j]
					squareV += degrees[j] * degrees[j]
				}
			} else if v < j {
				if sg.G.Edges[(j*(j-1))/2+v] == 1 {
					sumV += degrees[j]
					squareV += degrees[j] * degrees[j]
				}
			}
		}
		if sumV > sum {
			return false
		} else if sumV < sum {
			viableBits ^= (1 << uint(v))
		} else if squareV > square {
			return false
		} else if squareV < square {
			viableBits ^= (1 << uint(v))
		}
	}
	if viableBits == 0 {
		return true
	}

	//TODO Lexicographic degrees?

	//We must now check the canonical isomorph
	if sg.Perm == nil {
		//Set the options for checking viability since it will sometimes let us return early.
		options.CheckViability = true
		options.ViableBits = viableBits
		getAutomorphismGroup(sg, op, storage, options)
	}

	if sg.Perm == nil {
		//Still nil means that we returned early.
		return false
	}

	correctAnswer := sg.Orbits.Find(n - 1)

	for _, u := range sg.Perm {
		if u == n-1 {
			return true
		}

		if viableBits>>uint(u)&1 == 1 {
			return correctAnswer == sg.Orbits.Find(u)
		}
	}
	return true
}
