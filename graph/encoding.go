package graph

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strings"
)

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)

//Graph6Decode returns the graph with Graph6 encoding s or an error.
//For the definition of the format see: https://users.cecs.anu.edu.au/~bdm/data/formats.txt
//The empty string decodes as the empty graph.
func Graph6Decode(s string) (*DenseGraph, error) {
	//Check the bytes are in 63-126 (inclusive). This could be done during decoding but it is neater to check here.
	for i := 0; i < len(s); i++ {
		if s[i] < 63 || s[i] > 126 {
			return &DenseGraph{}, fmt.Errorf("Byte out of range. Index: %v Value: %v", i, s[i])
		}
	}

	var n uint64
	i := 0

	if len(s) == 0 {
		return NewDense(0, nil), nil
	}

	if s[0] != 126 {
		n = uint64(s[0] - 63)
		i = 1
	} else if s[1] != 126 {
		if len(s) < 4 {
			return &DenseGraph{}, errors.New("String too short - unable to decode n")
		}
		n = (uint64(s[1]-63) << 12) + (uint64(s[2]-63) << 6) + uint64(s[3]-63)
		i = 4
	} else {
		if len(s) < 8 {
			return &DenseGraph{}, errors.New("String too short - unable to decode n")
		}
		n = (uint64(s[2]-63) << 30) + (uint64(s[3]-63) << 24) + (uint64(s[4]-63) << 18) + (uint64(s[5]-63) << 12) + (uint64(s[6]-63) << 6) + uint64(s[7]-63)
		i = 8
		MaxN := 0.5 + math.Sqrt(2*float64(maxInt)+0.25)
		if float64(n) > MaxN {
			return &DenseGraph{}, errors.New("Graph too large")
		}
	}

	if i+int(((n*(n-1))/2)+5)/6 > len(s) {
		return &DenseGraph{}, errors.New("String too short - unable to decode edges")
	}

	edges := make([]uint8, (n*(n-1))/2)
	for j := 0; j < len(edges); j++ {
		edges[j] = (uint8(s[i+j/6]-63) & (1 << uint(5-(j%6)))) >> uint(5-(j%6))
	}

	return NewDense(int(n), edges), nil
}

//Graph6Encode returns the Graph6 encoding of g.
//For the definition of the format see: https://users.cecs.anu.edu.au/~bdm/data/formats.txt
func Graph6Encode(g Graph) string {
	var s []byte
	n := g.N()
	if n <= 1 {
		return string(n + 63)
	} else if n <= 62 {
		s = make([]byte, 1, 1+((n*(n-1))/2+5)/6)
		s[0] = byte(n + 63)
	} else if n <= 258047 {
		s = make([]byte, 4, 4+((n*(n-1))/2+5)/6)
		s[0] = 126
		s[1] = byte((n>>12)&63) + 63
		s[2] = byte((n>>6)&63) + 63
		s[3] = byte(n&63) + 63
	} else if n <= 68719476735 {
		s = make([]byte, 8, 8+((n*(n-1))/2+5)/6)
		s[0] = 126
		s[1] = 126
		s[2] = byte((n>>30)&63) + 63
		s[3] = byte((n>>24)&63) + 63
		s[4] = byte((n>>18)&63) + 63
		s[5] = byte((n>>12)&63) + 63
		s[6] = byte((n>>6)&63) + 63
		s[7] = byte(n&63) + 63
	} else {
		panic("Graph too large")
	}

	var b byte
	bIndex := 0
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if g.IsEdge(i, j) {
				b += 1 << uint(5-bIndex)
			}
			bIndex++
			if bIndex == 6 {
				s = append(s, b+63)
				bIndex = 0
				b = 0
			}
		}
	}

	if bIndex != 0 {
		s = append(s, b+63)
	}

	return string(s)
}

//Sparse6Decode decode returns the graph with Sparse6 encoding s or an error if no such graph exists.
//For the definition of the format see: https://users.cecs.anu.edu.au/~bdm/data/formats.txt
func Sparse6Decode(s string) (*SparseGraph, error) {
	//Check the initial byte and remove it.
	if s[0] != 58 {
		return &SparseGraph{}, fmt.Errorf("Incorrect first character. Expected: : Found: %v", s[0])
	}
	s = s[1:]

	//Check every byte is in the correct range.
	for i := 0; i < len(s); i++ {
		if s[i] < 63 || s[i] > 126 {
			return &SparseGraph{}, fmt.Errorf("Byte out of range (63-126). Index: %v Value: %v", i, s[i])
		}
	}

	var n uint64
	i := 0

	if s[0] != 126 {
		n = uint64(s[0] - 63)
		i = 1
	} else if s[1] != 126 {
		if len(s) < 4 {
			return &SparseGraph{}, errors.New("String too short - unable to decode n")
		}
		n = (uint64(s[1]-63) << 12) + (uint64(s[2]-63) << 6) + uint64(s[3]-63)
		i = 4
	} else {
		if len(s) < 8 {
			return &SparseGraph{}, errors.New("String too short - unable to decode n")
		}
		n = (uint64(s[2]-63) << 30) + (uint64(s[3]-63) << 24) + (uint64(s[4]-63) << 18) + (uint64(s[5]-63) << 12) + (uint64(s[6]-63) << 6) + uint64(s[7]-63)
		i = 8
	}

	g := NewSparse(int(n), nil)
	v := 0
	k := 64 - bits.LeadingZeros64(uint64(n-1))
	var bitIndex uint
	for true {
		b := ((s[i] - 63) >> (5 - bitIndex)) & 1
		bitIndex++
		if bitIndex == 6 {
			bitIndex = 0
			i++
			if i >= len(s) {
				return g, nil
			}
		}
		if b == 1 {
			v++
		}
		x := 0
		for j := 0; j < k; j++ {
			if ((s[i]-63)>>(5-bitIndex))&1 == 1 {
				x |= 1 << uint(k-j-1)
			}
			bitIndex++
			if bitIndex == 6 {
				bitIndex = 0
				i++
				if i >= len(s) {
					return g, nil
				}
			}
		}
		if x > v {
			v = x
		} else {
			g.AddEdge(v, x)
		}
	}
	return g, nil
}

//Sparse6Encode returns an encoding of g. Note that the encoding is not unique but this should align with the format used by showg, geng, nauty etc.
//For the definition of the format see: https://users.cecs.anu.edu.au/~bdm/data/formats.txt
func Sparse6Encode(g Graph) string {
	n := g.N()
	m := g.M()
	//Number of bits needed to express n -1.
	k := 64 - bits.LeadingZeros64(uint64(n-1))

	var s []byte

	//To encode an edge we may possibly need to a 1 + k bits to move to the vertex and then another 1 + k bits to encode the edge.

	if n <= 1 {
		return string([]byte{58, byte(n + 63)})
	} else if n <= 62 {
		s = make([]byte, 2, 2+((k+1)*2*m+5)/6)
		s[0] = 58
		s[1] = byte(n + 63)
	} else if n <= 258047 {
		s = make([]byte, 5, 5+((k+1)*2*m+5)/6)
		s[0] = 58
		s[1] = 126
		s[2] = byte((n>>12)&63) + 63
		s[3] = byte((n>>6)&63) + 63
		s[4] = byte(n&63) + 63
	} else if n <= 68719476735 {
		s = make([]byte, 9, 9+((k+1)*2*m+5)/6)
		s[0] = 58
		s[1] = 126
		s[2] = 126
		s[3] = byte((n>>30)&63) + 63
		s[4] = byte((n>>24)&63) + 63
		s[5] = byte((n>>18)&63) + 63
		s[6] = byte((n>>12)&63) + 63
		s[7] = byte((n>>6)&63) + 63
		s[8] = byte(n&63) + 63
	} else {
		panic("Graph too large")
	}

	var b byte

	v := 0
	currentBitPosition := 0
	for i := 1; i < n; i++ {
		neighbours := g.Neighbours(i)
		for _, u := range neighbours {
			if u > i {
				break
			}

			if i == v {
				//b[-] == 0
				currentBitPosition++
				if currentBitPosition == 6 {
					s = append(s, b+63)
					b = 0
					currentBitPosition = 0
				}
				//x[-] = u
				for j := 0; j < k; j++ {
					if (u>>uint(k-j-1))&1 == 1 {
						b |= 1 << uint(5-currentBitPosition)
					}
					currentBitPosition++
					if currentBitPosition == 6 {
						s = append(s, b+63)
						b = 0
						currentBitPosition = 0
					}
				}
			} else if i == v+1 {
				v++
				//b[-] == 1
				b |= 1 << uint(5-currentBitPosition)
				currentBitPosition++
				if currentBitPosition == 6 {
					s = append(s, b+63)
					b = 0
					currentBitPosition = 0
				}
				//x[-] = u
				for j := 0; j < k; j++ {
					if (u>>uint(k-j-1))&1 == 1 {
						b |= 1 << uint(5-currentBitPosition)
					}
					currentBitPosition++
					if currentBitPosition == 6 {
						s = append(s, b+63)
						b = 0
						currentBitPosition = 0
					}
				}
			} else {
				v = i
				//First we move.

				//I believe this is arbitrarily 0 or 1 but I think 1 is used by default.
				//b[-] == 1
				b |= 1 << uint(5-currentBitPosition)
				currentBitPosition++
				if currentBitPosition == 6 {
					s = append(s, b+63)
					b = 0
					currentBitPosition = 0
				}
				//x[-] = i
				for j := 0; j < k; j++ {
					if (i>>uint(k-j-1))&1 == 1 {
						b |= 1 << uint(5-currentBitPosition)
					}
					currentBitPosition++
					if currentBitPosition == 6 {
						s = append(s, b+63)
						b = 0
						currentBitPosition = 0
					}
				}

				//Now add the edge.

				//b[-] == 0
				currentBitPosition++
				if currentBitPosition == 6 {
					s = append(s, b+63)
					b = 0
					currentBitPosition = 0
				}
				//x[-] = u
				for j := 0; j < k; j++ {
					if (u>>uint(k-j-1))&1 == 1 {
						b |= 1 << uint(5-currentBitPosition)
					}
					currentBitPosition++
					if currentBitPosition == 6 {
						s = append(s, b+63)
						b = 0
						currentBitPosition = 0
					}
				}

			}

		}
	}

	if currentBitPosition == 0 {
		return string(s)
	}

	//Padding

	if (n == 2 || n == 4 || n == 8 || n == 16) && 6-currentBitPosition > k+1 {
		degrees := g.Degrees()
		if degrees[n-2] > 0 && degrees[n-1] == 0 {
			currentBitPosition++
		}
	}

	//Pad with 1s.
	for j := currentBitPosition; j < 6; j++ {
		b += 1 << uint(5-j)
	}

	s = append(s, b+63)

	return string(s)
}

//MulticodeEncode returns the Multicode encoding of g.
func MulticodeEncode(g Graph) []byte {
	n := g.N()
	if n > 255 {
		panic("Graph too large for Multicode")
	} else if n == 0 {
		return []byte{0}
	}
	s := make([]byte, g.M()+n)
	s[0] = byte(n)
	index := 1
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if g.IsEdge(i, j) {
				s[index] = byte(j + 1)
				index++
			}
		}
		s[index] = 0
		index++
	}
	return s
}

//MulticodeDecode returns the graph with Multicode encoding s.
func MulticodeDecode(s []byte) *DenseGraph {
	n := int(s[0])
	m := 0
	degrees := make([]int, n)
	edges := make([]byte, (n*(n-1))/2)
	currentVertex := 0
	for i := 1; i < len(s); i++ {
		if s[i] == 0 {
			currentVertex++
		} else {
			edges[(int(s[i]-1)*int(s[i]-2))/2+currentVertex] = 1
			degrees[s[i]]++
			degrees[currentVertex]++
			m++
		}
	}
	if n > 0 && currentVertex != n-1 {
		panic("Incomplete encoding")
	}
	return &DenseGraph{NumberOfVertices: n, NumberOfEdges: m, DegreeSequence: degrees, Edges: edges}
}

//MulticodeDecodeMultiple returns an array of graphs encoded in s.
func MulticodeDecodeMultiple(s []byte) []*DenseGraph {
	graphs := make([]*DenseGraph, 0)
	startOfGraph := 0
	var numberOfVerticesLeft byte
	for i := 0; i < len(s); i++ {
		if numberOfVerticesLeft == 0 {
			numberOfVerticesLeft = s[i] - 1
			startOfGraph = i
		}
		if s[i] == 0 {
			numberOfVerticesLeft--
			if numberOfVerticesLeft == 0 {
				graphs = append(graphs, MulticodeDecode(s[startOfGraph:i+1]))
			}
		}
	}
	return graphs
}

//PruferEncode returns the Prüfer code of the labelled tree g.
//https://en.wikipedia.org/wiki/Pr%C3%BCfer_sequence
func PruferEncode(g Graph) []int {
	n := g.N()
	verticesLeftToRemove := make([]int, n)
	for i := range verticesLeftToRemove {
		verticesLeftToRemove[i] = i
	}
	prufer := make([]int, 0, n-2)
	degrees := g.Degrees()
	for i := 0; i < n-2; i++ {
		for j, v := range verticesLeftToRemove {
			if degrees[v] == 1 {
				//Find the unique neighbour of v
				for _, u := range verticesLeftToRemove {
					if g.IsEdge(u, v) {
						prufer = append(prufer, u)
						degrees[u]--
						break
					}
				}
				copy(verticesLeftToRemove[j:], verticesLeftToRemove[j+1:])
				break
			}
		}
	}
	return prufer
}

//PruferDecode returns the labelled tree given by the Prüfer code p.
func PruferDecode(p []int) *DenseGraph {
	n := len(p) + 2
	degrees := make([]int, n)
	for i := 0; i < n; i++ {
		degrees[i] = 1
	}
	edges := make([]byte, (n*(n-1))/2)
	for _, v := range p {
		degrees[v]++
	}
	for _, v := range p {
		for j := 0; j < n; j++ {
			if degrees[j] == 1 {
				if j > v {
					edges[(j*(j-1))/2+v] = 1
				} else {
					edges[(v*(v-1))/2+j] = 1
				}
				degrees[j]--
				degrees[v]--
				break
			}

		}
	}
	for i := 0; i < n; i++ {
		if degrees[i] != 1 {
			continue
		}
		for j := i + 1; j < n; j++ {
			if degrees[j] == 1 {
				edges[(j*(j-1))/2+i] = 1
				break
			}
		}
		break
	}
	return &DenseGraph{NumberOfVertices: n, Edges: edges}
}

//AdjacencyMatrixEncode returns an encoding of the adjacency matrix of g suitable for copying into MATLAB.
//The matrix starts with a [ and ends with a ]. Elements on the same row are separated by , and rows are separated by ;
func AdjacencyMatrixEncode(g Graph) string {
	var sb strings.Builder
	n := g.N()
	if n == 0 {
		sb.WriteString("[]")
		return sb.String()
	}
	sb.WriteString("[")
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1; j++ {
			if g.IsEdge(i, j) {
				sb.WriteString("1,")
			} else {
				sb.WriteString("0,")
			}
		}
		if g.IsEdge(i, n-1) {
			sb.WriteString("1")
		} else {
			sb.WriteString("0")
		}
		sb.WriteString(";")
	}
	for j := 0; j < n-1; j++ {
		if g.IsEdge(n-1, j) {
			sb.WriteString("1,")
		} else {
			sb.WriteString("0,")
		}
	}
	if g.IsEdge(n-1, n-1) {
		sb.WriteString("1")
	} else {
		sb.WriteString("0")
	}
	sb.WriteString("]")
	return sb.String()
}
