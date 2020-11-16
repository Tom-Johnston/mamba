package search_test

import (
	"fmt"
	"math/bits"
	"os"
	"time"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/graph/search"
)

func Example_permutationGraphs() {
	//This program generates all circle graphs on 8 vertices using a canonical deletion method and writes the Graph6 encoding of each graph to os.Stdout.

	//Time how long it takes to generate the graphs.
	start := time.Now()

	n := 8
	fmt.Fprintf(os.Stderr, "Enumerating permutation graphs on %v vertices.\n", n)

	//Reuse the storage space for the isPermutationGraph function.
	U := make([]uint, n)
	D := make([]uint, n)
	implicants := []intPair{}
	neighbours := make([]uint, n)

	//Initialise an iterator.
	iter := search.WithPruning(n, 0, 1, func(g *graph.DenseGraph) bool { return false }, func(g *graph.DenseGraph) bool { return !isPermutationGraph(g, U, D, neighbours, implicants) })

	//Counter to keep track of how many graphs we find.
	counter := 0
	//Keep iterating until there are no more graphs.
	for iter.Next() {
		//Get the value of the iterator. Note that we must not edit the value.
		g := iter.Value()
		//Encode the graph and write to Stdout.
		s := graph.Graph6Encode(g)
		fmt.Println(s)
		counter++
	}

	fmt.Fprintln(os.Stderr, "Graphs: ", counter)
	fmt.Fprintf(os.Stderr, "Took %v\n", time.Since(start))
}

func isPermutationGraph(g *graph.DenseGraph, U, D, neighbours []uint, implicants []intPair) bool {

	n := g.NumberOfVertices

	neighbours = neighbours[:n]
	for i := range neighbours {
		neighbours[i] = 0
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if g.IsEdge(i, j) {
				neighbours[i] |= (1 << uint(j))
				neighbours[j] |= (1 << uint(i))
			}
		}
	}

	if !isTransitive(neighbours, U, D, implicants) {
		return false
	}

	var mask uint = (1 << uint(n)) - 1
	for i := range neighbours {
		neighbours[i] = ^neighbours[i]
		neighbours[i] ^= (1 << uint(i))
		neighbours[i] &= mask
	}

	return isTransitive(neighbours, U, D, implicants)
}

type intPair struct {
	i int
	j int
}

func isTransitive(neighbours []uint, U, D []uint, implicants []intPair) bool {
	n := len(neighbours)
	U = U[:n]
	copy(U, neighbours)
	D = D[:n]

	var edge intPair
	implicants = implicants[:0]
algLoop:
	for {
		for i := range D {
			D[i] = 0
		}

		//Step A
		for i, v := range U {
			if v != 0 {
				j := bits.TrailingZeros(v)
				U[i] ^= (1 << uint(j))
				U[j] ^= (1 << uint(i))
				D[i] |= (1 << uint(j))
				implicants = append(implicants, intPair{i: i, j: j})
				break
			}
		}
		//Step B
		for len(implicants) > 0 {
			edge, implicants = implicants[len(implicants)-1], implicants[:len(implicants)-1]
			i, j := edge.i, edge.j
			//Get the i' which are neighbours of i but not neighbours of j. Note we
			c := U[i] & (^neighbours[j])
			//Iterate over all i' and direct the edge from i to i'
			for c != 0 {
				y := c & -c
				v := bits.TrailingZeros(c)
				c ^= y

				U[i] ^= (1 << uint(v))
				U[v] ^= (1 << uint(i))
				D[i] ^= (1 << uint(v))
				implicants = append(implicants, intPair{i: i, j: v})
			}

			//Get the j' which are neighbours of j but not neighbours of i.
			c = U[j] & (^neighbours[i])
			//Iterate over all j' and direct the edge from j' to j
			for c != 0 {
				y := c & -c
				v := bits.TrailingZeros(c)
				c ^= y

				U[j] ^= (1 << uint(v))
				U[v] ^= (1 << uint(j))
				D[v] ^= (1 << uint(j))
				implicants = append(implicants, intPair{i: v, j: j})
			}
		}

		//Step C
		//The TRD test

		for i, c := range D {
			var w uint
			for c != 0 {
				y := c & -c
				j := bits.TrailingZeros(c)
				c ^= y

				w |= D[j]
			}
			//Refresh c
			c = D[i]

			//Check if there are any bits in w which are not set in c
			if w&(^c) != 0 {
				return false
			}
		}

		//Step D
		//Check if there are any undirected edges left.
		for _, v := range U {
			if v != 0 {
				continue algLoop
			}
		}

		return true
	}
}
