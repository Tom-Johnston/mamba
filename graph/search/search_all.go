package search

import (
	"runtime"
	"sort"

	"github.com/Tom-Johnston/gigraph/comb"
	"github.com/Tom-Johnston/gigraph/disjoint"
	. "github.com/Tom-Johnston/gigraph/graph"
)

//AllParallel sends all non-isomorphic graphs on n vertices to the output channel which it then closes. It automatically splits the work across GOMAXPROCS goroutines.
func AllParallel(n int, output chan *DenseGraph) {
	m := runtime.GOMAXPROCS(0)
	counter := m
	count := make(chan bool)
	for i := 0; i < m; i++ {
		c := make(chan *DenseGraph)
		tmp := i
		go func() {
			go All(n, c, tmp, m)
			for v := range c {
				output <- v
			}
			count <- false
		}()
	}
	for true {
		<-count
		counter--
		if counter == 0 {
			close(output)
		}
	}
}

//All with a= 0 and m = 1 sends all non-isomorphic graphs on n vertices to the output channel which it then closes. In general, this will begin a search using canonical deletion to find all graphs on at most n vertices where the choice at level ceil(2n/3) is equal to a mod m. For small values of m this should produce a fairly even split and allow for some small parallelism.
func All(n int, output chan *DenseGraph, a int, m int) {
	var object *DenseGraph
	augs := make([][]int, 1)
	objs := make([]*DenseGraph, 1)
	objectsToCheck := []*DenseGraph{NewDense(0, nil)}
	splitLevel := 2 * (n + 1) / 3
	if n == 0 {
		if a == 0 {
			output <- objectsToCheck[0]
		}
		close(output)
		return
	}
	for len(objectsToCheck) > 0 {
		object, objectsToCheck = objectsToCheck[len(objectsToCheck)-1], objectsToCheck[:len(objectsToCheck)-1]

		if object.NumberOfVertices == n {
			output <- object
			continue
		}
		augs = getAugmentations(object, augs)
		objs = objs[:0]
		for _, v := range augs {
			if w := applyAugmentation(object, v); isCanonical(object, v, w) {
				objs = append(objs, w)
			}
		}

		for i := range objs {
			if object.NumberOfVertices+1 != splitLevel || i%m == a {
				objectsToCheck = append(objectsToCheck, objs[i])
			}
		}
	}
	close(output)
}

func pruneAndOutput(gI interface{}, output chan interface{}) bool {
	g := gI.(DenseGraph)
	output <- g
	return false
}

func getAugmentations(g *DenseGraph, augs [][]int) [][]int {
	n := g.NumberOfVertices
	minDegree := MinDegree(g)
	maxSize := minDegree + 1
	augs = augs[:0]

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	f := func(i, j int) int {
		return 1
	}
	_, _, generators := CanonicalIsomorphCustom(g, f, 1, OrderedPartition{Order: order, BinSizes: []int{n}, Path: []int{}, SplitPoint: 0})
	for k := 0; k <= maxSize; k++ {
		ds := disjoint.New(comb.Coeff(n, k))
		for i := 0; i < len(ds); i++ {
			c := comb.Unrank(i, k)
			c2 := make([]int, k)
			for _, g := range generators {
				for j := 0; j < k; j++ {
					c2[j] = g[c[j]]
				}
				sort.Ints(c2)
				ds.Union(i, comb.Rank(c2))
			}
		}
		for i := 0; i < len(ds); i++ {
			if ds[i] < 0 {
				augs = append(augs, comb.Unrank(i, k))
			}
		}
	}
	return augs
}

func applyAugmentation(g *DenseGraph, aug []int) *DenseGraph {
	newGraph := g.Copy()
	newGraph.AddVertex(aug)
	n := g.N() + 1
	edges := make([]byte, (n*(n-1))/2)
	copy(edges, g.Edges)
	p := ((n - 1) * (n - 2)) / 2
	for _, v := range aug {
		edges[p+v] = 1
	}
	return newGraph.(*DenseGraph)
}

func isCanonical(g *DenseGraph, aug []int, h *DenseGraph) bool {
	n := g.NumberOfVertices

	viable := make([]int, 0, n+1)
	degrees := h.Degrees()

	//Check the degree
	degree := degrees[n]
	for i := 0; i < n+1; i++ {
		if degrees[i] < degree {
			return false
		} else if degrees[i] == degree {
			viable = append(viable, i)
		}
	}

	if len(viable) == 1 {
		return true
	}

	//Sums of degrees
	sum := 0
	for i := 0; i < n; i++ {
		if h.Edges[(n*(n-1))/2+i] == 1 {
			sum += degrees[i]
		}
	}
	sumV := 0
	for i := len(viable) - 1; i >= 0; i-- {
		sumV = 0
		for j := 0; j < n+1; j++ {
			if v := viable[i]; v > j {
				if h.Edges[(v*(v-1))/2+j] == 1 {
					sumV += degrees[j]
				}
			} else if v < j {
				if h.Edges[(j*(j-1))/2+v] == 1 {
					sumV += degrees[j]
				}
			}
		}

		if sumV > sum {
			return false
		} else if sumV < sum {
			viable[i] = viable[len(viable)-1]
			viable = viable[:len(viable)-1]
		}
	}
	if len(viable) == 1 {
		return true
	}

	//TODO Degree sequence

	order := make([]int, h.NumberOfVertices)
	for i := 0; i < h.NumberOfVertices; i++ {
		order[i] = i
	}
	f := func(i, j int) int {
		return 1
	}
	perm, orbits, _ := CanonicalIsomorphCustom(h, f, 1, OrderedPartition{Order: order, BinSizes: []int{h.NumberOfVertices}, Path: []int{}, SplitPoint: 0})
	for _, u := range perm {
		for _, v := range viable {
			if u == v {
				if orbits.Find(h.NumberOfVertices-1) == orbits.Find(u) {
					return true
				}
				return false
			}
		}
	}
	return true
}
