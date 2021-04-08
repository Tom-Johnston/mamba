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

func TestMultisetCombinations(t *testing.T) {
	correctOutput := [][]int{
		{0, 0, 0, 0, 1},
		{0, 0, 0, 1, 1},
		{0, 0, 1, 1, 1},
		{0, 0, 0, 0, 2},
		{0, 0, 0, 1, 2},
		{0, 0, 1, 1, 2},
		{0, 1, 1, 1, 2},
		{0, 0, 0, 2, 2},
		{0, 0, 1, 2, 2},
		{0, 1, 1, 2, 2},
		{1, 1, 1, 2, 2},
		{0, 0, 2, 2, 2},
		{0, 1, 2, 2, 2},
		{1, 1, 2, 2, 2},
		{0, 0, 0, 0, 3},
		{0, 0, 0, 1, 3},
		{0, 0, 1, 1, 3},
		{0, 1, 1, 1, 3},
		{0, 0, 0, 2, 3},
		{0, 0, 1, 2, 3},
		{0, 1, 1, 2, 3},
		{1, 1, 1, 2, 3},
		{0, 0, 2, 2, 3},
		{0, 1, 2, 2, 3},
		{1, 1, 2, 2, 3},
		{0, 2, 2, 2, 3},
		{1, 2, 2, 2, 3},
		{0, 0, 0, 3, 3},
		{0, 0, 1, 3, 3},
		{0, 1, 1, 3, 3},
		{1, 1, 1, 3, 3},
		{0, 0, 2, 3, 3},
		{0, 1, 2, 3, 3},
		{1, 1, 2, 3, 3},
		{0, 2, 2, 3, 3},
		{1, 2, 2, 3, 3},
		{2, 2, 2, 3, 3},
	}
	iter := MultisetCombinations([]int{4, 3, 3, 2}, 5)
	index := 0
	for iter.Next() {
		if index >= len(correctOutput) {
			t.Log("Found too many multisets.")
			t.FailNow()
		}
		found := iter.Value()
		truth := correctOutput[index]
		if len(found) != len(truth) {
			t.FailNow()
		}
		for i := range found {
			if found[i] != truth[i] {
				t.FailNow()
			}
		}
		index++
	}
	if index != len(correctOutput) {
		t.Error("Found too few multisets.")
	}
}
