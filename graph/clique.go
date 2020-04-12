package graph

import (
	"math/rand"

	"github.com/Tom-Johnston/mamba/sortints"
)

//cliqueData is a wrapper used in the CliqueNumber.
type cliqueData struct {
	R []int
	P []int
	X []int
}

//AllMaximalCliques returns every maximal clique in g.
//This uses the Bron–Kerbosch algorithm with pivots chosen to reduce the number of branches at each point and currently doesn't use a vertex ordering for the first pass.
func AllMaximalCliques(g Graph, c chan []int) {
	n := g.N()

	R := make([]int, 0)
	P := make([]int, n)
	for i := range P {
		P[i] = i
	}
	X := make([]int, 0)
	var cd cliqueData
	toCheck := make([]cliqueData, 1)
	toCheck[0] = cliqueData{R, P, X}
	for len(toCheck) > 0 {
		cd, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]
		P = cd.P
		R = cd.R
		X = cd.X
		if len(P) == 0 && len(X) == 0 {
			c <- R
			continue
		}
		//Choose a pivot vertex
		pivotVertex := -1
		bestPivotSize := -1
		pivotSize := 0
		for _, v := range P {
			pivotSize = 0
			for _, u := range P {
				//TODO Would this be nicer to be the intersection of the neighbourhoods so it is quicker for low degree vertices.
				if u != v && g.IsEdge(u, v) {
					pivotSize++
				}
			}
			if pivotSize > bestPivotSize {
				bestPivotSize = pivotSize
				pivotVertex = v
			}
		}

		for _, v := range X {
			pivotSize = 0
			for _, u := range P {
				//TODO As above
				if u != v && g.IsEdge(u, v) {
					pivotSize++
				}
			}
			if pivotSize > bestPivotSize {
				bestPivotSize = pivotSize
				pivotVertex = v
			}
		}

		for i := len(P) - 1; i >= 0; i-- {
			v := P[i]

			if v != pivotVertex && g.IsEdge(v, pivotVertex) {
				continue
			}
			tmpR := make([]int, len(R)+1)
			tmpP := make([]int, 0, len(P))
			tmpX := make([]int, 0, len(X)+1)
			copy(tmpR, R)
			tmpR[len(tmpR)-1] = v
			for _, u := range P {
				if u != v && g.IsEdge(u, v) {
					tmpP = append(tmpP, u)
				}
			}
			for _, u := range X {
				if u != v && g.IsEdge(u, v) {
					tmpX = append(tmpX, u)
				}
			}

			toCheck = append(toCheck, cliqueData{tmpR, tmpP, tmpX})
			P[i] = P[len(P)-1]
			P = P[:len(P)-1]
			X = append(X, v)
		}
	}
	close(c)
}

//CliqueNumber returns the size of the largest clique in g.
//This effectively finds all maximal cliques by using the Bron–Kerbosch algorithm with pivots chosen to reduce the number of branches at each point. This doesn't currently use a vertex ordering for the first pass.
func CliqueNumber(g Graph) int {
	n := g.N()

	cliqueNumber := 0
	R := make([]int, 0)
	P := make([]int, n)
	for i := range P {
		P[i] = i
	}
	X := make([]int, 0)
	var cd cliqueData
	toCheck := make([]cliqueData, 1)
	toCheck[0] = cliqueData{R, P, X}
	for len(toCheck) > 0 {
		cd, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]
		P = cd.P
		R = cd.R
		X = cd.X
		if len(P) == 0 && len(X) == 0 {
			if cliqueNumber < len(R) {
				cliqueNumber = len(R)
			}
			continue
		}
		//Choose a pivot vertex
		pivotVertex := -1
		bestPivotSize := -1
		pivotSize := 0
		for _, v := range P {
			pivotSize = 0
			for _, u := range P {
				//TODO Would this be nicer to be the intersection of the neighbourhoods so it is quicker for low degree vertices.
				if u != v && g.IsEdge(u, v) {
					pivotSize++
				}
			}
			if pivotSize > bestPivotSize {
				bestPivotSize = pivotSize
				pivotVertex = v
			}
		}

		for _, v := range X {
			pivotSize = 0
			for _, u := range P {
				//TODO As above
				if u != v && g.IsEdge(u, v) {
					pivotSize++
				}
			}
			if pivotSize > bestPivotSize {
				bestPivotSize = pivotSize
				pivotVertex = v
			}
		}

		for i := len(P) - 1; i >= 0; i-- {
			v := P[i]

			if v != pivotVertex && g.IsEdge(v, pivotVertex) {
				continue
			}
			tmpR := make([]int, len(R)+1)
			tmpP := make([]int, 0, len(P))
			tmpX := make([]int, 0, len(X)+1)
			copy(tmpR, R)
			tmpR[len(tmpR)-1] = v
			for _, u := range P {
				if u != v && g.IsEdge(u, v) {
					tmpP = append(tmpP, u)
				}
			}
			for _, u := range X {
				if u != v && g.IsEdge(u, v) {
					tmpX = append(tmpX, u)
				}
			}
			toCheck = append(toCheck, cliqueData{tmpR, tmpP, tmpX})
			P[i] = P[len(P)-1]
			P = P[:len(P)-1]
			X = append(X, v)
		}
	}
	return cliqueNumber
}

//IndependenceNumber returns the size of the largest independent set in g.
//The independence number of g is the clique number of the complement of g and this is how it is calculated here.
func IndependenceNumber(g Graph) int {
	h := Complement(g)
	return CliqueNumber(h)
}

//RandomMaximalClique builds a random maximal clique by iteratively choosing an allowed vertex uniformly at random and adding it to the clique.
func RandomMaximalClique(g Graph, seed int64) []int {
	r := rand.New(rand.NewSource(seed))
	clique := make([]int, 0)
	options := sortints.Range(0, g.N(), 1)
	for len(options) > 0 {
		index := r.Intn(len(options))
		option := options[index]
		clique = append(clique, option)
		options = sortints.Intersection(options, g.Neighbours(option))
	}
	return clique
}
