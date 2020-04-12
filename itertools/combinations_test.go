package itertools

import (
	"testing"
)

func TestCombinations(t *testing.T) {
	//Truth data
	correctOutput := [][]int{
		[]int{},
	}

	iter := Combinations(4, 0)
	index := 0
	for iter.Next() {
		if index >= len(correctOutput) {
			t.Log("Found too many combinations.")
			t.FailNow()
		}
		found := iter.Value()
		truth := correctOutput[index]
		if len(found) != len(truth) {
			t.Fail()
		}
		for i := range found {
			if found[i] != truth[i] {
				t.Fail()
			}
		}
		index++
	}
	if index != len(correctOutput) {
		t.Error("Found too few combinations.")
	}
}
