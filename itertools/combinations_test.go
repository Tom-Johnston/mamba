package itertools

import (
	"testing"
)

func TestCombinations(t *testing.T) {
	//Truth data
	inputParameters := [][]int{{0, 0}, {4, 0}, {5, 1}, {5, 2}}
	correctOutputs := [][][]int{
		{{}},
		{{}},
		{{0}, {1}, {2}, {3}, {4}},
		{{0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}},
	}
	for i := range inputParameters {
		inputs := inputParameters[i]
		correctOutput := correctOutputs[i]
		iter := Combinations(inputs[0], inputs[1])
		index := 0
		for iter.Next() {
			if index >= len(correctOutput) {
				t.Log(index)
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
}

func TestCombinationsColex(t *testing.T) {
	//Truth data
	inputParameters := [][]int{{0, 0}, {4, 0}, {5, 1}, {5, 2}}
	correctOutputs := [][][]int{
		{{}},
		{{}},
		{{0}, {1}, {2}, {3}, {4}},
		{{0, 1}, {0, 2}, {1, 2}, {0, 3}, {1, 3}, {2, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
	}
	for i := range inputParameters {
		inputs := inputParameters[i]
		correctOutput := correctOutputs[i]
		iter := CombinationsColex(inputs[0], inputs[1])
		index := 0
		for iter.Next() {
			if index >= len(correctOutput) {
				t.Log(index)
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
}
