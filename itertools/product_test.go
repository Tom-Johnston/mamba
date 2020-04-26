package itertools

import "testing"

func TestProduct(t *testing.T) {
	correctOutput := [][]int{}

	iter := Product(1, 0, 1)
	index := 0
	for iter.Next() {
		if index >= len(correctOutput) {
			t.Log("Found too many product.")
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
		t.Error("Found too few product.")
	}

	correctOutput = [][]int{
		{0, 0, 0},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{0, 2, 0},
		{0, 2, 1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 0},
		{1, 1, 1},
		{1, 2, 0},
		{1, 2, 1},
		{2, 0, 0},
		{2, 0, 1},
		{2, 1, 0},
		{2, 1, 1},
		{2, 2, 0},
		{2, 2, 1},
	}

	iter = Product(3, 3, 2)
	index = 0
	for iter.Next() {
		if index >= len(correctOutput) {
			t.Log("Found too many options.")
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
		t.Error("Found too few options.")
	}
}
