package utils_test

import (
	"testing"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
)

func TestIsSubset(t *testing.T) {

	type subsetTest struct {
		A   []string
		B   []string
		res bool
	}

	spec := []subsetTest{
		subsetTest{
			A:   []string{"a"},
			B:   []string{"a", "b"},
			res: true,
		},
		subsetTest{
			A:   []string{"a", "b"},
			B:   []string{"a", "b", "c"},
			res: true,
		},
		subsetTest{
			A:   []string{"a", "b", "c"},
			B:   []string{"a", "b"},
			res: false,
		},
		subsetTest{
			A:   []string{"a", "b", "c"},
			B:   []string{"a", "b", "d"},
			res: false,
		},
	}

	for _, s := range spec {
		actualRes := utils.IsSubset(s.A, s.B)

		if actualRes != s.res {
			t.Fail()
		}
	}
}
