package dawg

import (
	"bufio"
	"bytes"
	"os"
	"sort"
	"testing"
)

var words [][]byte = [][]byte{[]byte("abject"), []byte("abjection"), []byte("abjections"), []byte("abjectly"), []byte("abjectness"), []byte("ablate"), []byte("ablated"), []byte("ablation"), []byte("ablations")}
var anagramWords [][]byte = [][]byte{[]byte("alerting"), []byte("altering"), []byte("integral"), []byte("post"), []byte("pots"), []byte("relating"), []byte("spot"), []byte("stop"), []byte("tops"), []byte("triangle"), []byte("ttps")}

func TestBuilder(t *testing.T) {
	db := new(Builder)
	for _, b := range words {
		err := db.Add(b)
		if err != nil {
			t.Fatal(err)
		}
	}
	dawg, err := db.Finish()
	if err != nil {
		t.Fatal(err)
	}
	if dawg.NumberOfWords() != len(words) {
		t.Fail()
	}
	if dawg.numberOfNodes() != 19 {
		t.Fail()
	}

}

func TestNew(t *testing.T) {
	dawg, err := New(words)
	if err != nil {
		t.Fatal(err)
	}
	if dawg.NumberOfWords() != len(words) {
		t.Fail()
	}
	if dawg.numberOfNodes() != 19 {
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	dawg, err := New(words)
	if err != nil {
		t.Fatal(err)
	}
	testPasses := words
	for _, test := range testPasses {
		if dawg.Contains(test) == false {
			t.Fail()
		}
	}

	testFails := [][]byte{[]byte("ab"), []byte(""), []byte("hello")}
	for _, test := range testFails {
		if dawg.Contains(test) == true {
			t.Fail()
		}
	}
}

func TestPattern(t *testing.T) {
	dawg, err := New(anagramWords)
	if err != nil {
		t.Fatal(err)
	}
	expectedResult := [][]byte{[]byte("post"), []byte("pots")}
	result := dawg.Pattern([]byte("po??"), 63)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}

	expectedResult = [][]byte{[]byte("tops"), []byte("ttps")}
	result = dawg.Pattern([]byte("t?ps"), 63)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}
}

func TestPatternSearcher(t *testing.T) {
	dawg, err := New(anagramWords)
	if err != nil {
		t.Fatal(err)
	}
	expectedResult := [][]byte{[]byte("post"), []byte("pots")}
	ps := NewPatternSearcher([]byte("po??"), 63)
	result := dawg.Search(ps)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}

	expectedResult = [][]byte{[]byte("tops"), []byte("ttps")}
	ps = NewPatternSearcher([]byte("t?ps"), 63)
	result = dawg.Search(ps)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}
}

func TestAnagrams(t *testing.T) {
	dawg, err := New(anagramWords)
	if err != nil {
		t.Fatal(err)
	}
	expectedResult := [][]byte{[]byte("post"), []byte("pots"), []byte("spot"), []byte("stop"), []byte("tops")}
	result := dawg.Anagrams([]byte("post"), 0)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}
	result = dawg.Anagrams([]byte{112, 111, 0, 0}, 0)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}

	result = dawg.Anagrams([]byte{112, 113, 0, 0}, 0)
	if len(result) != 0 {
		t.Log(result)
		t.Fail()
	}

	result = dawg.Anagrams([]byte("ttps"), 0)
	if len(result) != 1 {
		t.Fatal(result)
	}
}

func TestAnagramSearcher(t *testing.T) {
	dawg, err := New(anagramWords)
	if err != nil {
		t.Fatal(err)
	}
	expectedResult := [][]byte{[]byte("post"), []byte("pots"), []byte("spot"), []byte("stop"), []byte("tops")}
	as := NewAnagramSearcher([]byte("post"), 0)
	result := dawg.Search(as)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}
	as = NewAnagramSearcher([]byte{112, 111, 0, 0}, 0)
	result = dawg.Search(as)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}

	as = NewAnagramSearcher([]byte{112, 113, 0, 0}, 0)
	result = dawg.Search(as)
	if len(result) != 0 {
		t.Log(result)
		t.Fail()
	}

	as = NewAnagramSearcher([]byte("ttps"), 0)
	result = dawg.Search(as)
	if len(result) != 1 {
		t.Fatal(result)
	}
}

func TestAnagramPattern(t *testing.T) {
	dawg, err := New(anagramWords)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := [][]byte{[]byte("tops")}
	as := NewAnagramSearcher([]byte("o???"), 63)
	ps := NewPatternSearcher([]byte("t?ps"), 63)
	result := dawg.Search(as, ps)
	if len(result) != len(expectedResult) {
		t.Fail()
	} else {
		for i := range result {
			if bytes.Equal(result[i], expectedResult[i]) == false {
				t.Fail()
			}
		}
	}
}

func TestGob(t *testing.T) {
	dawg, err := New(words)
	if err != nil {
		t.Fatal(err)
	}
	b, err := dawg.GobEncode()
	if err != nil {
		t.Fatal(err)
	}
	ts := new(Dawg)
	err = ts.GobDecode(b)
	if err != nil {
		t.Fatal(err)
	}
	if ts.NumberOfWords() != len(words) {
		t.Fail()
	}
	if ts.numberOfNodes() != 19 {
		t.Fail()
	}

	dawg, err = crossWord()
	b, err = dawg.GobEncode()
	if err != nil {
		t.Fatal(err)
	}
	err = ts.GobDecode(b)
	if err != nil {
		t.Fatal(err)
	}
	if ts.NumberOfWords() != dawg.NumberOfWords() {
		t.Fail()
	}
	if ts.numberOfNodes() != dawg.numberOfNodes() {
		t.Fail()
	}

}

func BenchmarkNew(b *testing.B) {
	file, err := os.Open("testdata/CROSSWD.TXT")
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	words := [][]byte{}
	for scanner.Scan() {
		words = append(words, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		b.Log(err)
		b.FailNow()
	}
	less := func(i, j int) bool {
		return bytes.Compare(words[i], words[j]) < 0
	}
	sort.Slice(words, less)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(words)
	}
}

func BenchmarkContains(b *testing.B) {
	dawg, err := crossWord()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dawg.Contains([]byte("alerting"))
	}
}

func BenchmarkPattern(b *testing.B) {
	dawg, err := crossWord()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dawg.Pattern([]byte("al???ing"), 63)
	}
}

func BenchmarkAnagram(b *testing.B) {
	dawg, err := crossWord()
	if err != nil {
		b.Fatal(err)
	}

	testStrings := []string{"stop", "alerting", "??", "?top", "alert???"}

	for _, s := range testStrings {
		b.Run(s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dawg.Anagrams([]byte(s), 63)
			}
		})
	}
	for i := 0; i < b.N; i++ {
		dawg.Anagrams([]byte("?top"), 63)
	}
}

func BenchmarkNumberOfNodes(b *testing.B) {
	dawg, err := crossWord()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dawg.numberOfNodes()
	}
}

func BenchmarkGobEncode(b *testing.B) {
	dawg, err := crossWord()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dawg.GobEncode()
	}
}

func BenchmarkGobDecode(b *testing.B) {
	dawg, err := crossWord()
	if err != nil {
		b.Fatal(err)
	}
	byt, err := dawg.GobEncode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = dawg.GobDecode(byt)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func crossWord() (*Dawg, error) {
	file, err := os.Open("testdata/CROSSWD.TXT")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	db := new(Builder)
	for scanner.Scan() {
		err := db.Add([]byte(scanner.Text()))
		if err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return db.Finish()
}
