package graph

import (
	"math/rand"

	"github.com/Tom-Johnston/gigraph/comb"
	"github.com/Tom-Johnston/gigraph/sortints"
)

//Random Graphs

//RandomGraph returns an Erdős–Rényi graph with n vertices and edge probability p. The pseudorandomness is determined by the seed.
func RandomGraph(n int, p float64, seed int64) *DenseGraph {
	r := rand.New(rand.NewSource(seed))
	g := NewDense(n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if r.Float64() < p {
				g.AddEdge(i, j)
			}
		}
	}
	return g
}

//RandomTree returns a tree chosen uniformly at random from all trees on n vertices.
//This constructs a random Prufer code and converts into a tree.
func RandomTree(n int, seed int64) *DenseGraph {
	code := make([]int, n-2)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < n-2; i++ {
		code[i] = r.Intn(n)
	}
	return PruferDecode(code)
}

//Specific Graphs

//CompleteGraph returns a copy of the complete graph on n vertices.
func CompleteGraph(n int) *DenseGraph {
	edges := make([]byte, (n*(n-1))/2)
	m := (n * (n - 1)) / 2
	degrees := make([]int, n)
	for i := range degrees {
		degrees[i] = n - 1
	}
	for i := 0; i < len(edges); i++ {
		edges[i] = 1
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: m, DegreeSequence: degrees, Edges: edges}
}

//CompletePartiteGraph returns the complete partite graph where the ith part has nums[i] elements.
func CompletePartiteGraph(nums ...int) *DenseGraph {
	n := 0
	for _, v := range nums {
		n += v
	}
	edges := make([]byte, (n*(n-1))/2)
	degrees := make([]int, n)
	m := 0
	start := 0
	end := 0
	for _, v := range nums {
		end += v
		degree := n - v
		for i := start; i < end; i++ {
			degrees[i] = degree
		}
		for k := end; k < n; k++ {
			for j := start; j < end; j++ {
				edges[(k*(k-1))/2+j] = 1
				m++
			}
		}
		start = end
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: m, DegreeSequence: degrees, Edges: edges}
}

//Path returns a copy of the path on n vertices.
func Path(n int) *DenseGraph {
	edges := make([]byte, (n*(n-1))/2)
	for i := 0; i < n-1; i++ {
		edges[((i+1)*i)/2+i] = 1
	}

	degrees := make([]int, n)
	if n > 0 {
		degrees[0] = 1
		degrees[n-1] = 1
		for i := 1; i < n-1; i++ {
			degrees[i] = 2
		}
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: n - 1, DegreeSequence: degrees, Edges: edges}
}

//Cycle returns a copy of the cycle on n vertices.
func Cycle(n int) *DenseGraph {
	edges := make([]byte, (n*(n-1))/2)
	for i := 0; i < n-1; i++ {
		edges[((i+1)*i)/2+i] = 1
	}
	edges[((n-1)*(n-2))/2] = 1

	degrees := make([]int, n)
	for i := range degrees {
		degrees[i] = 2
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: n, DegreeSequence: degrees, Edges: edges}
}

//Star returns a copy of the star on n vertices.
func Star(n int) *DenseGraph {
	edges := make([]byte, (n*(n-1))/2)
	for i := 1; i < n; i++ {
		edges[(i*(i-1))/2] = 1
	}

	degrees := make([]int, n)
	if n > 0 {
		degrees[0] = n - 1
		for i := 1; i < n; i++ {
			degrees[i] = 1
		}
	}

	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: n - 1, DegreeSequence: degrees, Edges: edges}
}

//RookGraph returns the n x m Rook graph i.e. the graph representing the moves of a rook on an n x m chessboard.
//If m = n, the graph is also called a Latin square graph.
func RookGraph(n, m int) *DenseGraph {
	return LineGraphDense(CompletePartiteGraph(n, m))
}

//FlowerSnark returns the flower snark on 4n vertices for n odd.
//See https://en.wikipedia.org/wiki/Flower_snark
func FlowerSnark(n int) *DenseGraph {
	if n&1 == 0 {
		panic("n must be odd")
	}
	N := 4 * n
	edges := make([]byte, (N*(N-1))/2)
	for i := 0; i < n; i++ {
		a := (4 * i)
		b := a + 1
		c := a + 2
		d := a + 3
		//Make a,b,c,d into a star with centre a.
		edges[(b*(b-1))/2+a] = 1
		edges[(c*(c-1))/2+a] = 1
		edges[(d*(d-1))/2+a] = 1
		if i < n-1 {
			edges[((b+4)*(b+3))/2+b] = 1
			edges[((c+4)*(c+3))/2+c] = 1
			edges[((d+4)*(d+3))/2+d] = 1
		} else {
			edges[(b*(b-1))/2+1] = 1
			edges[(c*(c-1))/2+3] = 1
			edges[(d*(d-1))/2+2] = 1
		}
	}
	return NewDense(N, edges)
}

//HypercubeGraph returns the dim-dimensional hypercube graph on 2^(dim) vertices.
func HypercubeGraph(dim int) *DenseGraph {
	if dim < 0 {
		panic("dim must be positive")
	}
	g := NewDense(1<<uint(dim), nil)
	for i := 0; i < 1<<uint(dim); i++ {
		var j uint
		for j = 0; j < uint(dim); j++ {
			g.AddEdge(i, i^(1<<uint(j)))
		}
	}
	return g
}

//FoldedHypercubeGraph returns a hypergraph on 2^{(dim - 1)} vertices except each point is joined with its antipodal point.
func FoldedHypercubeGraph(dim int) *DenseGraph {
	if dim < 1 {
		panic("dim must be at least 1")
	}
	g := HypercubeGraph(dim - 1)
	mask := 1<<uint(dim-1) - 1
	for i := 0; i < 1<<uint(dim-2); i++ {
		g.AddEdge(i, mask&^i)
	}
	return g
}

//KneserGraph returns the Kneser graph with parameters n and k. It has a vertex for each unordered subset of {0,...,n-1} of size k and an edge between two subsets if they are disjoint.
//The subsets are ordered in co-lexicographic.
func KneserGraph(n, k int) *DenseGraph {
	N := comb.Coeff(n, k)
	g := NewDense(N, nil)
	for i := 0; i < N; i++ {
		combi := comb.Unrank(i, k)
		for j := i; j < N; j++ {
			combj := comb.Unrank(j, k)
			if sortints.IntersectionSize(combi, combj) == 0 {
				g.AddEdge(i, j)
			}
		}
	}
	return g
}

//BipartiteKneserGraph returns the a bipartite graph vertices [n]^{(k)} on one side and [n]^{(n-k)} on the other and an edge between two sets if they are on different sides and one is a subset of the other.
func BipartiteKneserGraph(n, k int) *DenseGraph {
	N := comb.Coeff(n, k)
	size, overflow := addHasOverflowed(N, N)
	if overflow {
		panic("calculating the size of the graph overflows int")
	}

	g := NewDense(size, nil)
	for i := 0; i < N; i++ {
		combi := comb.Unrank(i, k)
		for j := 0; j < N; j++ {
			combj := comb.Unrank(j, n-k)
			if sortints.IntersectionSize(combi, combj) == k {
				g.AddEdge(i, N+j)
			}
		}
	}
	return g
}

//CirculantGraph generates a graph on n vertices where there is an edge between i and j if j-i is in diffs.
func CirculantGraph(n int, diffs ...int) *DenseGraph {
	g := NewDense(n, nil)
	for i := 0; i < n; i++ {
		for _, v := range diffs {
			targetVertex := (i + v) % n
			if targetVertex < 0 {
				targetVertex += n
			}
			g.AddEdge(i, targetVertex)
		}
	}
	return g
}

//CirculantBipartiteGraph creates a bipartite graph with vertices a_1, \dots, a_n in one class, b_1, \dots, b_m in the other and a edge between a_i and b_j if j - i is in diffs (mod m).
func CirculantBipartiteGraph(n int, m int, diffs ...int) *DenseGraph {
	g := NewDense(n+m, nil)
	for i := 0; i < n; i++ {
		for _, v := range diffs {
			targetVertex := (i + v) % m
			if targetVertex < 0 {
				targetVertex += m
			}
			g.AddEdge(i, n+targetVertex)
		}
	}
	return g
}

//GeneralisedPetersenGraph returns the graph with vertices {u_0,...u_{n-1}, v_0,...v_{n-1}} with edges {u_i u_{i+1}, u_i v_i, v_i v_{i+k}: i = 0,...,n − 1} where subscripts are to be read modulo n.
// n must be at least 3 and k must be between 1 and floor((n-1)/2).
//The Petersen graph is GeneralisedPetersenGraph(5,2) and several other named graphs are generalised Petersen graphs.
func GeneralisedPetersenGraph(n, k int) *DenseGraph {
	if n < 3 {
		panic("n must be at least 3.")
	}
	if k < 0 || k > (n-1)/2 {
		panic("k must be at least 1 and at most floor((n-1)/2).")
	}
	g := NewDense(2*n, nil)
	for i := 0; i < n; i++ {
		g.AddEdge(i, (i+1)%n)
		g.AddEdge(i, n+i)
		g.AddEdge(n+i, n+((i+k)%n))
	}
	return g
}

//FriendshipGraph returns the graph woth 2n + 1 vertices and 3n edges which is n copies of C_3 which share a common vertex.
func FriendshipGraph(n int) *DenseGraph {
	g := NewDense(2*n+1, nil)
	for i := 0; i < n; i++ {
		g.AddEdge(2*i+1, 2*i+2)
		g.AddEdge(0, 2*i+1)
		g.AddEdge(0, 2*i+2)
	}
	return g
}
