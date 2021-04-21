package search_test

import (
	"fmt"
	"math/bits"
	"os"
	"strings"
	"time"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/graph/search"
)

func Example_circleGraphs() {

	//This program generates all circle graphs on 8 vertices using a canonical deletion method and writes the Graph6 encoding of each graph to os.Stdout.

	//Time how long it takes to generate the graphs.
	start := time.Now()

	n := 8
	fmt.Fprintf(os.Stderr, "Enumerating circle graphs on %v vertices.\n", n)

	//Reuse the storage space for the isCircleGraph function.
	mat := Zeros(0, 0)
	b := make([]byte, 0)
	neighbours := make([]uint, n)

	//Make an iterator.
	iter := search.WithPruning(n, 0, 1, func(g *graph.DenseGraph) bool { return false }, func(g *graph.DenseGraph) bool { return !isCircleGraph(g, mat, b, neighbours) })

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

//isCircleGraph checks if the graph g is a circle graph. It will crash if n is larger than the number of bits in uint.
//It will use the space allocated in mat, b and neighbours. neighbours must be of length at least n.
//The method used is in Naji, Walid. "Reconnaissance des graphes de cordes." Discrete Mathematics 54.3 (1985): 329-337.
func isCircleGraph(g *graph.DenseGraph, mat *Matrix, b []byte, neighbours []uint) bool {

	components := graph.ConnectedComponents(g)

	for _, comp := range components {
		n := len(comp)

		numVars := n * n //We will have variables for (i,i) even if we don't need them.
		mat.N = numVars
		mat.M = 0
		for i := range mat.Entries {
			mat.Entries[i] = 0
		}
		mat.Entries = mat.Entries[:0]
		mat.RowPerm = mat.RowPerm[:0]
		numEqns := 0
		b = b[:0]
		for i := range neighbours {
			neighbours[i] = 0
		}
		neighbours = neighbours[:n]
		for i, v := range comp {
			for j := i + 1; j < n; j++ {
				u := comp[j]
				if g.IsEdge(v, u) {
					neighbours[i] |= (1 << uint(j))
					neighbours[j] |= (1 << uint(i))
				}
			}
		}

		for i := 0; i < n; i++ {
			v := comp[i]
			for j := i + 1; j < n; j++ {
				u := comp[j]
				if g.IsEdge(v, u) {
					//Find the vertices which are non-neighbours of both i and j.
					c := (^neighbours[i]) & (^neighbours[j])
					//Restrict to the first n bits.
					c &= (1 << uint(n)) - 1

					mat.AddRows(bits.OnesCount(c) + 1)
					//First type of equation
					mat.SetEntry(numEqns, i*n+j, 1)
					mat.SetEntry(numEqns, j*n+i, 1)
					b = append(b, 1)
					numEqns++

					//Other type of equation
					for c != 0 {
						y := c & -c
						v := bits.TrailingZeros(c)
						c ^= y

						mat.SetEntry(numEqns, v*n+i, 1)
						mat.SetEntry(numEqns, v*n+j, 1)
						b = append(b, 0)
						numEqns++
					}
					continue
				}
				c := neighbours[i] & neighbours[j]
				mat.AddRows(bits.OnesCount(c))
				for c != 0 {
					y := c & -c
					v := bits.TrailingZeros(c)
					c ^= y
					mat.SetEntry(numEqns, v*n+i, 1)
					mat.SetEntry(numEqns, v*n+j, 1)
					mat.SetEntry(numEqns, i*n+j, 1)
					mat.SetEntry(numEqns, j*n+i, 1)
					b = append(b, 1)
					numEqns++
				}
			}
		}
		if mat.HasSolution(b) == false {
			return false
		}
	}
	return true
}

//Matrix is a M x N dense matrix over F2.
//Each entry in the matrix is a bit and they are ordered left to right, top to bottom. The bits in a row are split into blocks of length 64 and each block is reversed before being encoded in a uint64.
//The entries are indexed from 0.
type Matrix struct {
	M       int
	N       int
	Entries []uint64
	RowPerm []int //The ith row in the current state of the matrix is the RowPerm[i]th row according to entries.
}

//Zeros initialises an M x N matrix which is all zeros.
func Zeros(M, N int) *Matrix {
	entries := make([]uint64, M*((N+63)/64))
	rowPerm := make([]int, M)
	for i := range rowPerm {
		rowPerm[i] = i
	}
	return &Matrix{M: M, N: N, Entries: entries, RowPerm: rowPerm}
}

//AddRows adds numRows extra zero rows to the bottom of the matrix.
func (m *Matrix) AddRows(numRows int) {
	width := (m.N + 63) / 64
	numNeeded := (m.M + numRows) * width
	if cap(m.Entries) >= numNeeded {
		m.Entries = m.Entries[:numNeeded]
		for i := m.M * width; i < numNeeded; i++ {
			m.Entries[i] = 0
		}
	} else {
		tmp := make([]uint64, numNeeded)
		copy(tmp, m.Entries)
		m.Entries = tmp
	}
	numNeeded = (m.M + numRows)
	if cap(m.RowPerm) >= numNeeded {
		m.RowPerm = m.RowPerm[:numNeeded]
		for i := m.M; i < numNeeded; i++ {
			m.RowPerm[i] = i
		}
	} else {
		tmp := make([]int, numNeeded)
		copy(tmp, m.RowPerm)
		m.RowPerm = tmp
		for i := m.M; i < numNeeded; i++ {
			m.RowPerm[i] = i
		}
	}
	m.M += numRows
}

//SetEntry sets the (i,j) entry to the value b.
//The value b must be either 0 or 1.
func (m *Matrix) SetEntry(i, j int, b byte) {
	r := m.RowPerm[i]
	switch b {
	case 0:
		m.Entries[r*((m.N+63)/64)+j/64] &^= (1 << uint(j%64))
		return
	case 1:
		m.Entries[r*((m.N+63)/64)+j/64] |= (1 << uint(j%64))
		return
	}
	panic("b is not 0 or 1")
}

//GetEntry returns the (i,j) entry of m.
func (m *Matrix) GetEntry(i, j int) uint64 {
	r := m.RowPerm[i]
	return (m.Entries[r*((m.N+63)/64)+j/64] >> uint(j%64)) & 1
}

//AddRowTo replaces the dst row with dst row + src row.
func (m *Matrix) AddRowTo(src, dst int) {
	rSrc := m.RowPerm[src]
	rDst := m.RowPerm[dst]
	width := (m.N + 63) / 64
	for k := 0; k < width; k++ {
		m.Entries[rDst*width+k] ^= m.Entries[rSrc*width+k]
	}
}

//SwapRows swaps the rows i and j.
func (m *Matrix) SwapRows(i, j int) {
	m.RowPerm[i], m.RowPerm[j] = m.RowPerm[j], m.RowPerm[i]
}

//Copy returns a deep copy of the matrix.
func (m *Matrix) Copy() *Matrix {
	entries := make([]uint64, len(m.Entries))
	copy(entries, m.Entries)
	rowPerm := make([]int, m.M)
	copy(rowPerm, m.RowPerm)
	return &Matrix{M: m.M, N: m.N, Entries: entries, RowPerm: rowPerm}
}

//HasSolution checks if the equation m*x = b has a solution over F_2^N using dense Gaussian Elimination.
//This modifies the matrix m.
func (m *Matrix) HasSolution(b []byte) bool {
	r := 0
	c := 0
mainLoop:
	for r < m.M && c < m.N {
		//Find a pivot row
		for i := r; i < m.M; i++ {
			if m.GetEntry(i, c) == 1 {
				//Swap the rows
				m.SwapRows(i, r)
				b[i], b[r] = b[r], b[i]
				for k := r + 1; k < m.M; k++ {
					if m.GetEntry(k, c) == 1 {
						m.AddRowTo(r, k)
						b[k] ^= b[r]
					}
				}
				//Next Iteration
				r++
				c++
				continue mainLoop
			}
		}
		//No pivot found. Skip this column and repeat.
		c++
	}
	for i := r; i < m.M; i++ {
		if b[i] != 0 {
			return false
		}
	}
	return true
}

func (m Matrix) String() string {
	var s strings.Builder
	for i := 0; i < m.M; i++ {
		for j := 0; j < m.N; j++ {
			err := s.WriteByte(48 + byte(m.GetEntry(i, j)))
			if err != nil {
				panic(err)
			}
		}
		err := s.WriteByte(10)
		if err != nil {
			panic(err)
		}
	}
	return s.String()
}
