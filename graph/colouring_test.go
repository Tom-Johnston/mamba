package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/graph/search"
	"github.com/Tom-Johnston/mamba/ints"
)

//From http://keithbriggs.info/cgt.html
// n =  1   2   3   4    5    6     7      8        9       10
// k ----------------------------------------------------------
// 0    0   0   0   0    0    0     0      0        0
// 1    1   1   1   1    1    1     1      1        1        1
// 2    0   1   2   6   12   34    87    302     1118     5478   A076278
// 3    0   0   1   3   16   84   579   5721    87381  2104349   A076279
// 4    0   0   0   1    4   31   318   5366   155291  7855628   A076280
// 5    0   0   0   0    1    5    52    867    28722  1919895   A076281
// 6    0   0   0   0    0    1     6     81     2028   115391   A076282
// 7    0   0   0   0    0    0     1      7      118     4251
// 8    0   0   0   0    0    0     0      1        8      165
// 9    0   0   0   0    0    0     0      0        1        9
// 10   0   0   0   0    0    0     0      0        0        1
// 11   0   0   0   0    0    0     0      0        0        0
func TestChromaticNumber(t *testing.T) {
	//truthData[n][k] contains the number of graphs of size n + 1 with chromatic numer k
	truthData := make([][]int, 10)
	truthData[0] = []int{0, 1}
	truthData[1] = []int{0, 1, 1}
	truthData[2] = []int{0, 1, 2, 1}
	truthData[3] = []int{0, 1, 6, 3, 1}
	truthData[4] = []int{0, 1, 12, 16, 4, 1}
	truthData[5] = []int{0, 1, 34, 84, 31, 5, 1}
	truthData[6] = []int{0, 1, 87, 579, 318, 52, 6, 1}
	truthData[7] = []int{0, 1, 302, 5721, 5366, 867, 81, 7, 1}
	truthData[8] = []int{0, 1, 1118, 87381, 155291, 28722, 2028, 118, 8, 1}
	truthData[9] = []int{0, 1, 5478, 2104349, 7855628, 1919895, 115391, 4251, 165, 9, 1}
	for i := 1; i <= 8; i++ {
		foundData := make([]int, i+1)

		iter := search.All(i, 0, 1)
		for iter.Next() {
			g := iter.Value()
			chromaticNumber, colouring := graph.ChromaticNumber(g)
			if !graph.IsProperColouring(g, colouring) {
				t.Fatal(g, colouring)
			}
			foundData[chromaticNumber]++
		}
		if !ints.Equal(foundData, truthData[i-1]) {
			t.Log(foundData)
			t.Log(truthData[i-1])
			t.FailNow()
		}
	}

	g, err := graph.Graph6Decode("KlWW[EHD_BsC")
	if err != nil {
		t.Error(err)
	}
	cn, colouring := graph.ChromaticNumber(g)
	if !graph.IsProperColouring(g, colouring) || cn != 3 {
		t.Fail()
	}

	g, err = graph.Graph6Decode("KCOedae]SrLu")
	if err != nil {
		t.Error(err)
	}
	cn, colouring = graph.ChromaticNumber(g)
	if !graph.IsProperColouring(g, colouring) || cn != 4 {
		t.Fail()
	}
}

// k  | n=  1  2  3  4   5   6    7     8       9       10
// --------------------------------------------------------
// 0  |     1  1  1  1   1   1    1     1       1        1
// 1  |     0  1  1  2   2   3    3     4       4        5
// 2  |     0  0  1  3   5  10   15    26      37       58
// 3  |     0  0  1  5  14  46  123   375    1061     3331
// 4  |     0  0  0  0  10  58  347  2130   14039   103927
// 5  |     0  0  0  0   2  38  392  4895   68696  1140623
// 6  |     0  0  0  0   0   0  159  3855  113774  3953535
// 7  |     0  0  0  0   0   0    4  1060   64669  4607132
// 8  |     0  0  0  0   0   0    0     0   12378  1921822
// 9  |     0  0  0  0   0   0    0     0       9   274734
// 10  |     0  0  0  0   0   0    0     0       0        0
func TestChromaticIndex(t *testing.T) {
	truthData := make([][]int, 10)
	truthData[0] = []int{1}
	truthData[1] = []int{1, 1}
	truthData[2] = []int{1, 1, 1, 1}
	truthData[3] = []int{1, 2, 3, 5}
	truthData[4] = []int{1, 2, 5, 14, 10, 2}
	truthData[5] = []int{1, 3, 10, 46, 58, 38}
	truthData[6] = []int{1, 3, 15, 123, 347, 392, 159, 4}
	truthData[7] = []int{1, 4, 26, 375, 2130, 4895, 3855, 1060}
	truthData[8] = []int{1, 4, 37, 1061, 14039, 68696, 113774, 64669, 12378, 9}
	for i := 1; i <= 7; i++ {
		var foundData []int
		if i%2 == 0 || i == 1 {
			foundData = make([]int, i)
		} else {
			foundData = make([]int, i+1)
		}

		iter := search.All(i, 0, 1)
		for iter.Next() {
			g := iter.Value()
			chromaticIndex, _ := graph.ChromaticIndex(g)
			foundData[chromaticIndex]++
		}
		if !ints.Equal(foundData, truthData[i-1]) {
			t.Log(foundData)
			t.Log(truthData)
			t.Fail()
		}
	}
}

func TestChromaticPolynomial(t *testing.T) {
	truthData := make(map[string][]int)
	truthData["IheA@GUAo"] = []int{0, -704, 2606, -4305, 4275, -2861, 1353, -455, 105, -15, 1}
	truthData["Dhc"] = []int{0, 4, -10, 10, -5, 1}
	for g6 := range truthData {
		g, err := graph.Graph6Decode(g6)
		if err != nil {
			t.Log(err)
			t.Fail()
			continue
		}
		cp := graph.ChromaticPolynomial(g)
		if !ints.Equal(cp, truthData[g6]) {
			t.Fail()
		}
	}
}
