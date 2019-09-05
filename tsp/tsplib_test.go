package tsp

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestLIB(t *testing.T) {
	n := 11
	shift := 100
	weights := func(i, j int) int {
		if i == j {
			return 0
		}
		if i > j {
			return shift*j + i
		}
		return shift*i + j
	}
	buf := new(bytes.Buffer)
	truthData, err := ioutil.ReadFile("testdata/lib.golden")
	if err != nil {
		t.Error(err)
	}
	LIB(buf, n, weights)
	if !bytes.Equal(buf.Bytes(), truthData) {
		t.Log(buf.Bytes())
		t.Log(truthData)
		t.Fail()
	}
}
